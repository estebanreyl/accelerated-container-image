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

package builder

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	testingresources "github.com/containerd/accelerated-container-image/cmd/convertor/testingresources"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func getResolver(t *testing.T, ctx context.Context) remotes.Resolver {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	localRegistryPath := fmt.Sprintf("%s/../testingresources/mocks/registry", cwd)
	resolver, err := testingresources.NewMockLocalResolver(ctx, localRegistryPath)
	if err != nil {
		t.Error(err)
	}
	return resolver
}

func getTestFetcherFromResolver(t *testing.T, ctx context.Context, resolver remotes.Resolver, ref string) remotes.Fetcher {
	fetcher, err := resolver.Fetcher(ctx, ref)
	if err != nil {
		t.Error(err)
	}
	return fetcher
}

// TESTS
func Test_fetchManifest(t *testing.T) {
	ctx := context.Background()
	resolver := getResolver(t, ctx)
	type args struct {
		ctx     context.Context
		fetcher remotes.Fetcher
		desc    v1.Descriptor
	}
	tests := []struct {
		name        string
		args        args
		want        *v1.Manifest
		wantErr     bool
		wantSubDesc v1.Descriptor
	}{
		{
			name: "Fetch existing manifest",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Manifest,
					Digest:    "sha256:92c7f9c92844bbbb5d0a101b22f7c2a7949e40f8ea90c8b3bc396879d95e899a",
					Size:      524,
				},
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Fetch manifest List",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/redis:alpine3.18"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2ManifestList,
					Digest:    "sha256:7bea1298286d063ed53cb52f3eaf4e70574269b7c0fe2b8f85ea699497e9cba6",
					Size:      2584,
				},
				ctx: ctx,
			},
			// The manifest list is expected to select the first manifest that can be converted
			// in the list, for this image that is the very first one.
			wantSubDesc: v1.Descriptor{
				MediaType: images.MediaTypeDockerSchema2Manifest,
				Digest:    "sha256:e20345b7ec692815860c07f0209eb0465687b0c28cd85df412811ae1ac7b653e",
				Size:      1571,
				Platform: &v1.Platform{
					Architecture: "amd64",
					OS:           "linux",
				},
			},
			wantErr: false,
		},
		{
			name: "Fetch unknown manifest errors",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Manifest,
					Digest:    "sha256:82c7f9c92844bbab5d0a101b12f7c2a7949e40f8ee90c8b3bc396879d95e899a",
					Size:      524,
				},
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Fetch invalid digest",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Manifest,
					Digest:    "sha256:829d95e899a",
					Size:      524,
				},
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Fetch manifest with unsupported mediaType (docker v1)",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema1Manifest,
					Digest:    "sha256:82c7f9c92844bbab5d0a101b12f7c2a7949e40f8ee90c8b3bc396879d95e899a",
					Size:      524,
				},
				ctx: ctx,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest, err := fetchManifest(tt.args.ctx, tt.args.fetcher, tt.args.desc)
			if (err == nil) && tt.wantErr {
				t.Error("fetchManifest() error was expected but no error was returned")
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fetchManifest() unexpectedly returned error %v", err)
				}
				return
			}
			content, err := testingresources.ConsistentManifestMarshal(manifest)
			if err != nil {
				t.Errorf("Could not parse obtained manifest, got: %v", err)
			}

			contentDigest := digest.FromBytes(content)

			if tt.args.desc.MediaType != images.MediaTypeDockerSchema2ManifestList &&
				tt.args.desc.MediaType != v1.MediaTypeImageIndex {

				if tt.args.desc.Digest != contentDigest {
					t.Errorf("fetchManifest() = %v, want %v", manifest, tt.want)
				}
			} else {
				if tt.wantSubDesc.Digest != contentDigest {
					t.Errorf("fetchManifest() = %v, want %v", manifest, tt.want)
				}
			}
		})
	}
}

func Test_fetchConfig(t *testing.T) {
	ctx := context.Background()
	resolver := getResolver(t, ctx)

	type args struct {
		ctx     context.Context
		fetcher remotes.Fetcher
		desc    v1.Descriptor
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.Image
		wantErr bool
	}{
		{
			name: "Fetch Config with supported mediaType (docker v2)",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Config,
					Digest:    "sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
					Size:      1510,
				},
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Fetch Config with supported mediaType (oci)",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest"),
				desc: v1.Descriptor{
					MediaType: v1.MediaTypeImageConfig,
					Digest:    "sha256:82c7f9c92844bbab5d0a101b12f7c2a7949e40f8ee90c8b3bc396879d95e899a",
					Size:      524,
				},
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Fetch unknown config",
			args: args{
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema1Manifest,
					Digest:    "sha256:82c7f9c92844bbab5d0a101b12f7c2a7949e40f8ee90c8b3bc396879d95e899a",
					Size:      524,
				},
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchConfig(tt.args.ctx, tt.args.fetcher, tt.args.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
