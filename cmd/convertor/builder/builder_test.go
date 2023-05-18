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
	"math/rand"
	"testing"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type mockBuilder struct {
	layers int
	config v1.Image
	engine builderEngine
}

func Test_overlaybdBuilder_Build(t *testing.T) {
	type fields struct {
		layers int
		config v1.Image
		engine builderEngine
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &overlaybdBuilder{
				layers: tt.fields.layers,
				config: tt.fields.config,
				engine: tt.fields.engine,
			}
			if err := b.Build(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("overlaybdBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockBuilderEngine struct{}

func NewmockBuilderEngine() builderEngine {
	return &mockBuilderEngine{}
}

func (e *mockBuilderEngine) DownloadLayer(ctx context.Context, idx int) error {
	if rand.Float64() < 0.3 {
		return fmt.Errorf("random error on download")
	}
	return nil
}

func (e *mockBuilderEngine) BuildLayer(ctx context.Context, idx int) error {
	if rand.Float64() < 0.3 {
		return fmt.Errorf("random error on download")
	}
	return nil
}

func (e *mockBuilderEngine) UploadLayer(ctx context.Context, idx int) error {
	if rand.Float64() < 0.3 {
		return fmt.Errorf("random error on download")
	}
	return nil
}

func (e *mockBuilderEngine) UploadImage(ctx context.Context) error {
	if rand.Float64() < 0.3 {
		return fmt.Errorf("random error on download")
	}
	return nil
}

func (e *mockBuilderEngine) CheckForConvertedLayer(ctx context.Context, idx int) (specs.Descriptor, error) {
	if rand.Float64() < 0.3 {
		return specs.Descriptor{}, fmt.Errorf("random error on download")
	}
	return specs.Descriptor{}, nil
}

func (e *mockBuilderEngine) StoreConvertedLayerDetails(ctx context.Context, idx int) error {
	if rand.Float64() < 0.3 {
		return fmt.Errorf("random error on download")
	}
	return nil
}

func (e *mockBuilderEngine) DownloadConvertedLayer(ctx context.Context, idx int, desc specs.Descriptor) error {
	if rand.Float64() < 0.3 {
		return fmt.Errorf("random error on download")
	}
	return nil
}

func (e *mockBuilderEngine) Cleanup() {
}
