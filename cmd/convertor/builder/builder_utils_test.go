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
	"testing"

	testingresources "github.com/containerd/accelerated-container-image/cmd/convertor/testingresources"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// TESTS
func Test_fetchManifest(t *testing.T) {
	ctx := context.Background()
	resolver := testingresources.GetTestResolver(t, ctx)
	_, desc, _ := resolver.Resolve(ctx, testingresources.Docker_Manifest_List_Ref)
	fmt.Println(desc)
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
				fetcher: testingresources.GetTestFetcherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Manifest,
					Digest:    testingresources.DockerV2_Manifest_Simple_Digest,
					Size:      525,
				},
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Fetch manifest List",
			args: args{
				fetcher: testingresources.GetTestFetcherFromResolver(t, ctx, resolver, testingresources.Docker_Manifest_List_Ref),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2ManifestList,
					Digest:    testingresources.Docker_Manifest_List_Digest,
					Size:      2069,
				},
				ctx: ctx,
			},
			// The manifest list is expected to select the first manifest that can be converted
			// in the list, for this image that is the very first one.
			wantSubDesc: v1.Descriptor{
				MediaType: images.MediaTypeDockerSchema2Manifest,
				Digest:    testingresources.DockerV2_Manifest_Simple_Digest,
				Size:      525,
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
				fetcher: testingresources.GetTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:idontexist"),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Manifest,
					Digest:    "sha256:82c7f9c92844bbbb5d0a101b12f7c2a7949e40f8ee90c8b3bc396879d95e899a",
					Size:      524,
				},
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Fetch invalid digest",
			args: args{
				fetcher: testingresources.GetTestFetcherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Manifest,
					Digest:    "sha256:829d95e899a",
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
	resolver := testingresources.GetTestResolver(t, ctx)

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
		// TODO: "Fetch Config with supported mediaType (oci)",
		{
			name: "Fetch Config with supported mediaType (docker v2)",
			args: args{
				fetcher: testingresources.GetTestFetcherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Config,
					Digest:    testingresources.DockerV2_Manifest_Simple_Config_Digest,
					Size:      1470,
					Platform: &v1.Platform{
						Architecture: "amd64",
						OS:           "linux",
					},
				},
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Fetch unknown config",
			args: args{
				fetcher: testingresources.GetTestFetcherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref),
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
			if (err == nil) && tt.wantErr {
				t.Error("fetchConfig() error was expected but no error was returned")
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fetchConfig() unexpectedly returned error %v", err)
				}
				return
			}
			if got.Architecture != tt.args.desc.Platform.Architecture ||
				got.OS != tt.args.desc.Platform.OS {
				t.Errorf("fetchConfig() config is not as expected")
			}

			if len(got.RootFS.DiffIDs) == 0 {
				t.Errorf("fetchConfig() Expected some DiffIds")
			}
			if len(got.History) == 0 {
				t.Errorf("fetchConfig() Expected layer history")
			}
		})
	}
}

func Test_uploadBytes(t *testing.T) {
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
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Config,
					Digest:    testingresources.DockerV2_Manifest_Simple_Config_Digest,
					Size:      1470,
					Platform: &v1.Platform{
						Architecture: "amd64",
						OS:           "linux",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchConfig(tt.args.ctx, tt.args.fetcher, tt.args.desc)
			if (err == nil) && tt.wantErr {
				t.Error("fetchConfig() error was expected but no error was returned")
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fetchConfig() unexpectedly returned error %v", err)
				}
				return
			}

		})
	}
}

func Test_uploadBlob(t *testing.T) {
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
				fetcher: getTestFetcherFromResolver(t, ctx, resolver, testingresources.DockerV2_Manifest_Simple_Ref),
				desc: v1.Descriptor{
					MediaType: images.MediaTypeDockerSchema2Config,
					Digest:    testingresources.DockerV2_Manifest_Simple_Config_Digest,
					Size:      1470,
					Platform: &v1.Platform{
						Architecture: "amd64",
						OS:           "linux",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchConfig(tt.args.ctx, tt.args.fetcher, tt.args.desc)
			if (err == nil) && tt.wantErr {
				t.Error("fetchConfig() error was expected but no error was returned")
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("fetchConfig() unexpectedly returned error %v", err)
				}
				return
			}

		})
	}
}

func Test_getFileDesc(t *testing.T) {

}

// TODO: Helper functions writing to the file system are not currently unit tested. This would involve mocking the
// existing filesystem access which involves additional refactoring out of the scope of initial coverage. The functions
// should be covered in more complete end to end tests but this will be a nice to have.
