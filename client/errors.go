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

package client

import "fmt"

// Indicates that a Name can't be bound to Path
type ErrNegBinding struct {
	Name string
}

func (e *ErrNegBinding) Error() string {
	return fmt.Sprint("Neg:", e.Name)
}

// Indicates that a Path resolution is still pending
type ErrResolutionPending struct {
	Path string
}

func (e *ErrResolutionPending) Error() string {
	return fmt.Sprint("Pending:", e.Path)
}
