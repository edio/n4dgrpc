// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/linkerd/linkerd/mesh/core/src/main/protobuf (interfaces: InterpreterClient)

// Package mock_protobuf is a generated GoMock package.
package mock_protobuf

import (
	gomock "github.com/golang/mock/gomock"
	protobuf "github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockInterpreterClient is a mock of InterpreterClient interface
type MockInterpreterClient struct {
	ctrl     *gomock.Controller
	recorder *MockInterpreterClientMockRecorder
}

// MockInterpreterClientMockRecorder is the mock recorder for MockInterpreterClient
type MockInterpreterClientMockRecorder struct {
	mock *MockInterpreterClient
}

// NewMockInterpreterClient creates a new mock instance
func NewMockInterpreterClient(ctrl *gomock.Controller) *MockInterpreterClient {
	mock := &MockInterpreterClient{ctrl: ctrl}
	mock.recorder = &MockInterpreterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInterpreterClient) EXPECT() *MockInterpreterClientMockRecorder {
	return m.recorder
}

// GetBoundTree mocks base method
func (m *MockInterpreterClient) GetBoundTree(arg0 context.Context, arg1 *protobuf.BindReq, arg2 ...grpc.CallOption) (*protobuf.BoundTreeRsp, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBoundTree", varargs...)
	ret0, _ := ret[0].(*protobuf.BoundTreeRsp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoundTree indicates an expected call of GetBoundTree
func (mr *MockInterpreterClientMockRecorder) GetBoundTree(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoundTree", reflect.TypeOf((*MockInterpreterClient)(nil).GetBoundTree), varargs...)
}

// StreamBoundTree mocks base method
func (m *MockInterpreterClient) StreamBoundTree(arg0 context.Context, arg1 *protobuf.BindReq, arg2 ...grpc.CallOption) (protobuf.Interpreter_StreamBoundTreeClient, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StreamBoundTree", varargs...)
	ret0, _ := ret[0].(protobuf.Interpreter_StreamBoundTreeClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamBoundTree indicates an expected call of StreamBoundTree
func (mr *MockInterpreterClientMockRecorder) StreamBoundTree(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamBoundTree", reflect.TypeOf((*MockInterpreterClient)(nil).StreamBoundTree), varargs...)
}
