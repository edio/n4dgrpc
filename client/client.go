// Copyright © 2017 Dmytro Kostiuchenko edio@archlinux.us
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Stateful client
// Call Connect() before calling anything else
// You may want also to set custom to exported variables
package client

import (
	"context"
	"fmt"
	"github.com/edio/n4dgrpc/convertions"
	mesh "github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

// exported
var (
	// Timeout for initial dial
	DialTimeout time.Duration = 1 * time.Second
	// Timeout for completing operation (shared if multiple calls to n4d are involved)
	OpTimeout time.Duration = 1 * time.Second
)

type ResolutionStrategy func(ctx context.Context, boundPath *mesh.Path) ([]*mesh.Endpoint, error)

var (
	ResolutionStrategySmart      ResolutionStrategy = resolveGetThenStream
	ResolutionStrategyNoStream   ResolutionStrategy = resolveGet
	ResolutionStrategyStreamOnly ResolutionStrategy = resolveStream
)

// private
var (
	connection  *grpc.ClientConn
	resolver    mesh.ResolverClient
	interpreter mesh.InterpreterClient
	ctx         = context.Background()
)

// Connect to n4d
func Connect(namerdAddress string) (err error) {
	lctx, cancel := context.WithTimeout(ctx, DialTimeout)
	defer cancel()
	connection, err = grpc.DialContext(lctx, namerdAddress, grpc.WithInsecure())
	resolver = mesh.NewResolverClient(connection)
	interpreter = mesh.NewInterpreterClient(connection)
	return
}

func Close() {
	if connection != nil {
		connection.Close()
		connection = nil
		resolver = nil
		resolver = nil
	}
}

// Bind a name in a namespace specified by root
func Bind(root *mesh.Path, name *mesh.Path) ([]*mesh.Path, error) {
	lctx, cancel := context.WithTimeout(ctx, OpTimeout)
	defer cancel()
	return bind(lctx, root, name)
}

func bind(ctx context.Context, root *mesh.Path, name *mesh.Path) ([]*mesh.Path, error) {
	bindReq := &mesh.BindReq{
		Root: root,
		Name: name,
	}

	resp, err := interpreter.GetBoundTree(ctx, bindReq)
	if err != nil {
		return nil, err
	}

	// TODO support not only leafs
	switch resp.Tree.Node.(type) {
	case *mesh.BoundNameTree_Leaf_:
		return []*mesh.Path{resp.Tree.GetLeaf().Id}, nil
	case *mesh.BoundNameTree_Neg_:
		return []*mesh.Path{}, &ErrNegBinding{Name: convertions.PathToStr(name)}
	default:
		return nil, fmt.Errorf("Not supported yet: %v", resp)
	}
	return nil, fmt.Errorf("Something unexpected has happened")
}

// Resolve a name in a namespace specified by root
func Resolve(root *mesh.Path, name *mesh.Path, resolveFn ResolutionStrategy) ([]*mesh.Endpoint, error) {
	lctx, cancel := context.WithTimeout(ctx, OpTimeout)
	defer cancel()

	boundPaths, err := bind(lctx, root, name)
	if err != nil || len(boundPaths) == 0 {
		return nil, err
	}

	var endpoints []*mesh.Endpoint
	// TODO concurrent resolution
	// TODO or stop resolution as soon as something is resolved. Provide --all option to resolve all
	// TODO clarify default linkerd behavior. Ask on linkerd slack?
	for _, path := range boundPaths {
		endpnts, errLast := resolveFn(lctx, path)
		if errLast != nil {
			log.Printf("Error resolving [%v]: %v", path, errLast)
		}
		endpoints = append(endpoints, endpnts...)
		err = errLast
	}

	// right now last resolution error is returned. But multiple bound trees are rare at least for me, ok so far
	// TODO collect all errors
	return endpoints, err
}

func resolveStream(ctx context.Context, boundPath *mesh.Path) ([]*mesh.Endpoint, error) {
	replicasReq := &mesh.ReplicasReq{
		Id: boundPath,
	}

	stream, err := resolver.StreamReplicas(ctx, replicasReq)
	if err != nil {
		return nil, err
	}
	defer stream.CloseSend()

	return recvEndpoints(stream)
}

func recvEndpoints(stream mesh.Resolver_StreamReplicasClient) (endpoints []*mesh.Endpoint, err error) {
	for endpoints == nil && err == nil {
		replicas, e := stream.Recv()
		if replicas != nil && replicas.GetBound() != nil {
			endpoints = replicas.GetBound().Endpoints
		}
		if e != nil {
			err = e
		}
	}
	if endpoints != nil && err == io.EOF {
		// do not treat EOF as error if endpoints received
		err = nil
	}
	return endpoints, err
}

func resolveGet(ctx context.Context, boundPath *mesh.Path) ([]*mesh.Endpoint, error) {
	replicasReq := &mesh.ReplicasReq{
		Id: boundPath,
	}

	resp, err := resolver.GetReplicas(ctx, replicasReq)
	if err != nil {
		return nil, err
	}

	endpoints := []*mesh.Endpoint{}
	if resp.GetBound() != nil {
		endpoints = resp.GetBound().Endpoints
		err = nil
	} else if resp.GetPending() != nil {
		err = &ErrResolutionPending{Path: convertions.PathToStr(boundPath)}
	} else {
		err = fmt.Errorf("Not supported: %v", resp.GetResult())
	}
	return endpoints, err
}

func resolveGetThenStream(ctx context.Context, boundPath *mesh.Path) ([]*mesh.Endpoint, error) {
	endpoints, err := resolveGet(ctx, boundPath)
	switch err.(type) {
	case *ErrResolutionPending:
		endpoints, err = resolveStream(ctx, boundPath)
		break
	default:
		break
	}

	return endpoints, err
}
