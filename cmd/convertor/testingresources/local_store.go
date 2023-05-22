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
	"io/fs"
	"path/filepath"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content/oci"
)

type RegistryOptions struct {
	InmemoryOnly      bool // Specifies if the registry should not load any resources from storage
	localRegistryPath string
}

type InmemoryRepo struct {
}

type RepoStore struct {
	path         string
	fileStore    *oci.Store
	inmemoryRepo *InmemoryRepo
	opts         *RegistryOptions
}

func (r *RepoStore) LoadStore(ctx context.Context) error {
	// File Store is already initialized
	if r.fileStore != nil {
		return nil
	}

	if r.opts.InmemoryOnly {
		return errors.New("LoadStore should not be invoked if registry is memory only")
	}

	if r.path != "" {
		return errors.New("LoadStore was not provided a path")
	}

	// Load an OCI layout store
	store, err := oci.New(r.path)
	if err != nil {
		return err
	}

	r.fileStore = store
	return nil
}

func (r *RepoStore) Fetch(ctx context.Context, descriptor v1.Descriptor) (io.ReadCloser, error) {
	if !r.opts.InmemoryOnly {
		if err := r.LoadStore(ctx); err != nil {
			return nil, err
		}

		return r.fileStore.Fetch(ctx, descriptor)
	}
	return nil, errors.New("Content not found")
}

func (r *RepoStore) Resolve(ctx context.Context, tag string) (v1.Descriptor, error) {
	if !r.opts.InmemoryOnly {
		if err := r.LoadStore(ctx); err != nil {
			return v1.Descriptor{}, err
		}

		return r.fileStore.Resolve(ctx, tag)
	}
	return v1.Descriptor{}, errors.New("Content not found")
}

type internalRegistry map[string]*RepoStore

type TestRegistry struct {
	internalRegistry internalRegistry
	opts             RegistryOptions
}

func NewTestRegistry(ctx context.Context, opts RegistryOptions) (*TestRegistry, error) {
	TestRegistry := TestRegistry{
		internalRegistry: make(internalRegistry),
	}
	if !opts.InmemoryOnly {
		err := filepath.WalkDir(opts.localRegistryPath, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				// Actual load from storage is done in a deferred manner as an optimization
				TestRegistry.internalRegistry[d.Name()] = &RepoStore{
					path: path,
				}
			}
			return nil
		})

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
	repo := r.internalRegistry[repository]
	return repo.Resolve(ctx, tag)
}

func (r *TestRegistry) Fetch(ctx context.Context, repository string, descriptor v1.Descriptor) (io.ReadCloser, error) {
	// Add in memory store for overrides/pushes
	if repo, ok := r.internalRegistry[repository]; ok {
		return repo.Fetch(ctx, descriptor)
	}
	return nil, errors.New("Repository not found")
}
