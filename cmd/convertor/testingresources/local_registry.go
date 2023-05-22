/*
Copyright The Accelerated Container Image Authors

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

package testingresources

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// REGISTRY
type TestRegistry struct {
	internalRegistry internalRegistry
	opts             RegistryOptions
}

type RegistryOptions struct {
	InmemoryOnly      bool // Specifies if the registry should not load any resources from storage
	localRegistryPath string
}

type internalRegistry map[string]*RepoStore

func NewTestRegistry(ctx context.Context, opts RegistryOptions) (*TestRegistry, error) {
	TestRegistry := TestRegistry{
		internalRegistry: make(internalRegistry),
		opts:             opts,
	}
	if !opts.InmemoryOnly {
		files, err := os.ReadDir(opts.localRegistryPath)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				// Actual load from storage is done in a deferred manner as an optimization
				TestRegistry.internalRegistry[file.Name()] = &RepoStore{
					path: filepath.Join(opts.localRegistryPath, file.Name()),
					opts: &opts,
				}
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return &TestRegistry, nil
}

func (r *TestRegistry) Resolve(ctx context.Context, ref string) (v1.Descriptor, error) {
	_, repository, tag, err := ParseRef(ctx, ref)
	if err != nil {
		return v1.Descriptor{}, err
	}
	if repo, ok := r.internalRegistry[repository]; ok {
		return repo.Resolve(ctx, tag)
	}
	return v1.Descriptor{}, errors.New("Repository not found")
}

func (r *TestRegistry) Fetch(ctx context.Context, repository string, descriptor v1.Descriptor) (io.ReadCloser, error) {
	// Add in memory store for overrides/pushes
	if repo, ok := r.internalRegistry[repository]; ok {
		return repo.Fetch(ctx, descriptor)
	}
	return nil, errors.New("Repository not found")
}
