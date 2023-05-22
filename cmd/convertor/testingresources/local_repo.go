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

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content/oci"
)

// REPOSITORY
type RepoStore struct {
	path         string
	fileStore    *oci.Store
	inmemoryRepo *inmemoryRepo
	opts         *RegistryOptions
}

type inmemoryRepo struct {
}

func (r *RepoStore) LoadStore(ctx context.Context) error {
	// File Store is already initialized
	if r.fileStore != nil {
		return nil
	}

	if r.opts.InmemoryOnly {
		return errors.New("LoadStore should not be invoked if registry is memory only")
	}

	if r.path == "" {
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

func (r *RepoStore) Resolve(ctx context.Context, tag string) (v1.Descriptor, error) {
	if !r.opts.InmemoryOnly {
		if err := r.LoadStore(ctx); err != nil {
			return v1.Descriptor{}, err
		}

		return r.fileStore.Resolve(ctx, tag)
	}
	return v1.Descriptor{}, errors.New("Content not found")
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
