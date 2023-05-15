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
	"os"
	"testing"

	"github.com/containerd/accelerated-container-image/cmd/convertor/builder/mocks"
	"github.com/containerd/containerd/remotes"
	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func Test_builderEngineBase_isGzipLayer(t *testing.T) {
	ctx := context.Background()
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	localRegistryPath := fmt.Sprintf("%s/mocks/registry", cwd)

	resolver, err := mocks.NewMockLocalResolver(ctx, localRegistryPath)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		fetcher  remotes.Fetcher
		manifest specs.Manifest
	}
	getFields := func(ctx context.Context, ref string) fields {
		_, desc, err := resolver.Resolve(ctx, ref)
		if err != nil {
			t.Error(err)
		}

		fetcher, err := resolver.Fetcher(ctx, ref)
		if err != nil {
			t.Error(err)
		}

		manifestStream, err := fetcher.Fetch(ctx, desc)
		if err != nil {
			t.Error(err)
		}

		if err != nil {
			t.Error(err)
		}

		parsedManifest := v1.Manifest{}
		decoder := json.NewDecoder(manifestStream)
		if err = decoder.Decode(&parsedManifest); err != nil {
			t.Error(err)
		}

		return fields{
			fetcher:  fetcher,
			manifest: parsedManifest,
		}
	}

	type args struct {
		ctx context.Context
		idx int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO Add more layers types for validation
		// Unknown Layer Type
		// Uncompressed Layer Type
		{
			name:   "Valid Gzip Layer",
			fields: getFields(ctx, "sample.localstore.io/hello-world:latest"),
			args: args{
				ctx: ctx,
				idx: 0,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Layer Not Found",
			fields: func() fields {
				fields := getFields(ctx, "sample.localstore.io/hello-world:latest")
				fields.manifest.Layers[0].Digest = digest.FromString("sample")
				return fields
			}(),
			args: args{
				ctx: ctx,
				idx: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &builderEngineBase{
				fetcher:  tt.fields.fetcher,
				manifest: tt.fields.manifest,
			}
			got, err := e.isGzipLayer(tt.args.ctx, tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("builderEngineBase.isGzipLayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("builderEngineBase.isGzipLayer() = %v, want %v", got, tt.want)
			}
		})
	}
}
