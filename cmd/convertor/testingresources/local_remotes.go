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

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/remotes"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// RESOLVER
type MockLocalResolver struct {
	testReg *TestRegistry
}

func NewMockLocalResolver(ctx context.Context, localRegistryPath string) (*MockLocalResolver, error) {
	reg, err := NewTestRegistry(ctx, RegistryOptions{
		localRegistryPath: localRegistryPath,
		InmemoryOnly:      false,
	})
	if err != nil {
		return nil, err
	}

	return &MockLocalResolver{
		testReg: reg,
	}, nil
}

func NewCustomMockLocalResolver(ctx context.Context, localRegistryPath string, testReg *TestRegistry) (*MockLocalResolver, error) {
	return &MockLocalResolver{
		testReg: testReg,
	}, nil
}

func (r *MockLocalResolver) Resolve(ctx context.Context, ref string) (string, v1.Descriptor, error) {
	desc, err := r.testReg.Resolve(ctx, ref)
	if err != nil {
		return "", v1.Descriptor{}, err
	}
	return "", desc, nil
}

func (r *MockLocalResolver) Fetcher(ctx context.Context, ref string) (remotes.Fetcher, error) {
	_, repository, _, err := ParseRef(ctx, ref)
	if err != nil {
		return nil, err
	}

	return &MockLocalFetcher{
		testReg:    r.testReg,
		repository: repository,
	}, nil
}

func (r *MockLocalResolver) Pusher(ctx context.Context, ref string) (remotes.Pusher, error) {
	_, repository, _, err := ParseRef(ctx, ref)
	if err != nil {
		return nil, err
	}

	return &MockLocalPusher{
		testReg:    r.testReg,
		repository: repository,
	}, nil
}

// FETCHER
type MockLocalFetcher struct {
	testReg    *TestRegistry
	repository string
}

func (f *MockLocalFetcher) Fetch(ctx context.Context, desc v1.Descriptor) (io.ReadCloser, error) {
	return f.testReg.Fetch(ctx, f.repository, desc)
}

// PUSHER
type MockLocalPusher struct {
	testReg    *TestRegistry
	repository string
}

// Not used by overlaybd conversion
func (p MockLocalPusher) Writer(ctx context.Context, opts ...content.WriterOpt) (content.Writer, error) {
	return nil, errors.New("Not implemented")
}

func (p MockLocalPusher) Push(ctx context.Context, desc v1.Descriptor) (content.Writer, error) {
	return nil, errors.New("Not implemented")
}
