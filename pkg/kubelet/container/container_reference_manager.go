/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package container

import (
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
)

// RefManager manages the references for the containers.
// The references are used for reporting events such as creation,
// failure, etc. This manager is thread-safe, no locks are necessary
// for the caller.
type RefManager struct {
	sync.RWMutex
	// TODO(yifan): To use strong type.
	containerIDToRef map[string]*api.ObjectReference
}

// NewRefManager creates and returns a container reference manager
// with empty contents.
func NewRefManager() *RefManager {
	return &RefManager{containerIDToRef: make(map[string]*api.ObjectReference)}
}

// SetRef stores a reference to a pod's container, associating it with the given container ID.
func (c *RefManager) SetRef(id string, ref *api.ObjectReference) {
	c.Lock()
	defer c.Unlock()
	c.containerIDToRef[id] = ref
}

// ClearRef forgets the given container id and its associated container reference.
// TODO(yifan): This is currently never called. Consider to remove this function,
// or figure out when to clear the references.
func (c *RefManager) ClearRef(id string) {
	c.Lock()
	defer c.Unlock()
	delete(c.containerIDToRef, id)
}

// GetRef returns the container reference of the given ID, or (nil, false) if none is stored.
func (c *RefManager) GetRef(id string) (ref *api.ObjectReference, ok bool) {
	c.RLock()
	defer c.RUnlock()
	ref, ok = c.containerIDToRef[id]
	return ref, ok
}
