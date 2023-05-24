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
	"encoding/json"
	"fmt"
	"path"
	"testing"

	testingresources "github.com/containerd/accelerated-container-image/cmd/convertor/testingresources"
	"github.com/containerd/containerd/images"
	_ "github.com/containerd/containerd/pkg/testutil" // Handle custom root flag
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
					Size:      testingresources.DockerV2_Manifest_Simple_Size,
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

// Simple test to ensure that the uploadBytes function is working as expected for a few different scenarios.
func Test_uploadBytes(t *testing.T) {
	ctx := context.Background()
	sourceManifest := testingresources.DockerV2_Manifest_Simple_Ref
	targetManifest := "sample.localstore.io/hello-world:another"
	resolver := testingresources.GetTestResolver(t, ctx)

	_, desc, err := resolver.Resolve(ctx, sourceManifest)
	if err != nil {
		t.Error(err)
	}
	fetcher := testingresources.GetTestFetcherFromResolver(t, ctx, resolver, sourceManifest)
	pusher := testingresources.GetTestPusherFromResolver(t, ctx, resolver, targetManifest)

	// Load manifest
	content, err := fetcher.Fetch(ctx, desc)
	if err != nil {
		t.Error(err)
	}

	test_uploadBytes := func(manifest v1.Manifest, pusher remotes.Pusher) error {
		manifestBytes, err := testingresources.ConsistentManifestMarshal(&manifest)
		if err != nil {
			return err
		}
		newDesc := v1.Descriptor{
			MediaType: images.MediaTypeDockerSchema2Manifest,
			Digest:    digest.FromBytes(manifestBytes),
			Size:      int64(len(manifestBytes)),
		}
		err = uploadBytes(ctx, pusher, newDesc, manifestBytes)
		if err != nil {
			return err
		}
		return nil
	}

	// Docker v2 manifest
	manifest := v1.Manifest{}
	json.NewDecoder(content).Decode(&manifest)

	// Re-Push Manifest  error should be handled
	testingresources.Assert(t, test_uploadBytes(manifest, testingresources.GetTestPusherFromResolver(t, ctx, resolver, sourceManifest)) == nil, "Could not upload Re upload Docker v2 Manifest with layers present") // Docker v2 manifest

	// Modify manifest to change digest
	manifest.Annotations = map[string]string{
		"test": "test",
	}
	testingresources.Assert(t, test_uploadBytes(manifest, pusher) == nil, "Could not upload Docker v2 Manifest with layers present") // Docker v2 manifest

	// OCI manifest
	manifest.MediaType = v1.MediaTypeImageManifest
	for i := range manifest.Layers {
		manifest.Layers[i].MediaType = v1.MediaTypeImageLayerGzip
	}
	testingresources.Assert(t, test_uploadBytes(manifest, pusher) == nil, "Could not upload OCI Manifest with layers present") // Docker v2 manifest

	// Missing layer
	manifest.Layers[0].Digest = digest.FromString("not there")
	testingresources.Assert(t, test_uploadBytes(manifest, pusher) != nil, "Expected layer not found error") // Docker v2 manifest
}

func Test_uploadBlob(t *testing.T) {
	ctx := context.Background()
	// Create a new inmemory registry to push to
	reg := testingresources.GetTestRegistry(t, ctx, testingresources.RegistryOptions{
		InmemoryOnly:              true,
		ManifestPushIgnoresLayers: false,
	})

	resolver := testingresources.GetCustomTestResolver(t, ctx, reg)
	pusher := testingresources.GetTestPusherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest")
	blobPath := path.Join(testingresources.GetLocalRegistryPath(), "hello-world", "blobs", "sha256", digest.Digest(testingresources.DockerV2_Manifest_Simple_Layer_0_Digest).Encoded())

	desc := v1.Descriptor{
		MediaType: images.MediaTypeDockerSchema2Manifest,
		Digest:    testingresources.DockerV2_Manifest_Simple_Layer_0_Digest,
		Size:      testingresources.DockerV2_Manifest_Simple_Layer_0_Size,
	}

	testingresources.Assert(t, uploadBlob(ctx, pusher, blobPath, desc) == nil, "uploadBlob() expected no error but got one")

	// Uploads already present shuld give no issues
	testingresources.Assert(t, uploadBlob(ctx, pusher, blobPath, desc) == nil, "uploadBlob() retry expected no error but got one")
	// Validate manifest of pushed blob
	fetcher := testingresources.GetTestFetcherFromResolver(t, ctx, resolver, "sample.localstore.io/hello-world:latest")
	blob, err := fetcher.Fetch(ctx, desc)
	if err != nil {
		t.Error(err)
	}
	blobDigest, err := digest.FromReader(blob)
	if err != nil {
		t.Error(err)
	}
	testingresources.Assert(t, blobDigest == desc.Digest, "uploadBlob() blob digest does not match stored value")
}

func Test_getFileDesc(t *testing.T) {
	test_getFileDesc := func(blobPath string, compressed bool, expectedDigest string, expectedSize int64) {
		desc, err := getFileDesc(blobPath, compressed)
		if err != nil {
			t.Error(err)
		}
		testingresources.Assert(t, desc.Digest.String() == expectedDigest, "getFileDesc() wrong digest returned")
		testingresources.Assert(t, desc.Size == expectedSize, "getFileDesc() wrong size returned")
	}
	blobPath := path.Join(testingresources.GetLocalRegistryPath(), "hello-world", "blobs", "sha256")

	// Compressed blob
	test_getFileDesc(
		path.Join(blobPath, digest.Digest(testingresources.DockerV2_Manifest_Simple_Layer_0_Digest).Encoded()),
		false,
		testingresources.DockerV2_Manifest_Simple_Layer_0_Digest,
		testingresources.DockerV2_Manifest_Simple_Layer_0_Size)

	// Uncompressed blob
	test_getFileDesc(
		path.Join(blobPath, digest.Digest(testingresources.DockerV2_Manifest_Simple_Config_Digest).Encoded()),
		false,
		testingresources.DockerV2_Manifest_Simple_Config_Digest,
		testingresources.DockerV2_Manifest_Simple_Config_Size)
}

// TODO: Helper functions writing to the file system are not currently unit tested. This would involve mocking the
// existing filesystem access which involves additional refactoring out of the scope of initial coverage. The functions
// should be covered in more complete end to end tests but this will be a nice to have.
