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

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"io"
// 	"reflect"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/containerd/containerd/content"
// 	"github.com/containerd/containerd/remotes"
// 	"github.com/opencontainers/go-digest"
// 	"github.com/opencontainers/image-spec/specs-go"
// 	specs "github.com/opencontainers/image-spec/specs-go/v1"
// 	v1 "github.com/opencontainers/image-spec/specs-go/v1"
// 	"github.com/pkg/errors"
// )

// // TESTS
// func Test_fetchManifest(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		fetcher remotes.Fetcher
// 		desc    v1.Descriptor
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *v1.Manifest
// 		wantErr bool
// 	}{
// 		{
// 			name: "Fetch existing manifest",
// 			args: args{fetcher: &mockHelloWorldFetcher{}, desc: v1.Descriptor{
// 				MediaType: helloworldManifest.MediaType,
// 				Digest:    digest.FromString("sha256:1234567890"),
// 				Size:      1234,
// 			}},
// 			want:    &helloworldManifest,
// 			wantErr: false,
// 		},
// 		{
// 			name: "Fetch unknown manifest",
// 			args: args{fetcher: &mockHelloWorldFetcher{}, desc: v1.Descriptor{
// 				MediaType: helloworldManifest.MediaType,
// 				Digest:    digest.FromString("sha256:1234567890"),
// 				Size:      1234,
// 			}},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fetch invalid descriptor",
// 			args: args{fetcher: &mockHelloWorldFetcher{}, desc: v1.Descriptor{
// 				MediaType: helloworldManifest.MediaType,
// 				Digest:    digest.FromString("sha256:1234567890"),
// 				Size:      1234,
// 			}},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fetch manifest with unsupported mediaType",
// 			args: args{fetcher: &mockHelloWorldFetcher{}, desc: v1.Descriptor{
// 				MediaType: helloworldManifest.MediaType,
// 				Digest:    digest.FromString("sha256:1234567890"),
// 				Size:      1234,
// 			}},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fetch manifest list",
// 			args: args{fetcher: &mockHelloWorldFetcher{}, desc: v1.Descriptor{
// 				MediaType: helloworldManifest.MediaType,
// 				Digest:    digest.FromString("sha256:1234567890"),
// 				Size:      1234,
// 			}},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := fetchManifest(tt.args.ctx, tt.args.fetcher, tt.args.desc)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("fetchManifest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("fetchManifest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_fetchConfig(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		fetcher remotes.Fetcher
// 		desc    v1.Descriptor
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *v1.Image
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := fetchConfig(tt.args.ctx, tt.args.fetcher, tt.args.desc)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("fetchConfig() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("fetchConfig() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// // We need a mock fetcher to fetch the manifest, config and layers of the above image
// // Also need to mock the content store to store the layers and the utils to commit, etc
// // Lets start by mocking the fetcher
// var helloworldManifest v1.Manifest = v1.Manifest{
// 	Versioned: specs.Versioned{
// 		SchemaVersion: 2,
// 	},
// 	MediaType: "application/vnd.docker.distribution.manifest.v2+json",
// 	Config: v1.Descriptor{
// 		MediaType: "application/vnd.docker.container.image.v1+json",
// 		Size:      1510,
// 		Digest:    "sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
// 	},
// 	Layers: []v1.Descriptor{
// 		{
// 			MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
// 			Size:      977,
// 			Digest:    "sha256:1b930d010525941c1d56ec53b97bd057a67ae1865eebf042686d2a2d18271ced",
// 		},
// 	},
// }

// // RESOLVER
// type mockHelloWorldResolver struct{}

// func (r *mockHelloWorldResolver) Resolve(ctx context.Context, ref string) (string, v1.Descriptor, error) {
// 	content, _ := json.MarshalIndent(helloworldManifest, "", "  ")

// 	return "", v1.Descriptor{
// 		MediaType: helloworldManifest.MediaType,
// 		Digest:    digest.FromBytes(content),
// 		Size:      int64(len(content)),
// 	}, nil
// }

// func (r *mockHelloWorldResolver) Fetcher(ctx context.Context, ref string) (remotes.Fetcher, error) {
// 	return &mockHelloWorldFetcher{}, nil
// }

// func (r *mockHelloWorldResolver) Pusher(ctx context.Context, ref string) (remotes.Pusher, error) {
// 	return &mockHelloWorldPusher{}, nil
// }

// // FETCHER
// type mockHelloWorldFetcher struct{}

// func (f *mockHelloWorldFetcher) Fetch(ctx context.Context, desc v1.Descriptor) (io.ReadCloser, error) {
// 	return io.NopCloser(strings.NewReader("")), nil
// }

// // PUSHER
// type mockHelloWorldPusher struct{}

// // Not used by overlaybd conversion
// func (p mockHelloWorldPusher) Writer(ctx context.Context, opts ...content.WriterOpt) (content.Writer, error) {
// 	return nil, errors.New("Not implemented")
// }

// func (p mockHelloWorldPusher) Push(ctx context.Context, desc v1.Descriptor) (content.Writer, error) {
// 	return nil, errors.New("Not implemented")
// }
