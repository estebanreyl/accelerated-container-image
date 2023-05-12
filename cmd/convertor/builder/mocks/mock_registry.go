package mocks

import (
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type mockImage struct {
	blobs        map[string]string
	manifest     *v1.Manifest
	manifestList *v1.Index
	config       *v1.Image
}

type mockRegistry struct {
	images map[string]mockImage
}

func NewImageRegistry() *mockRegistry {
	return &mockRegistry{
		images: make(map[string]mockImage),
	}
}

func (r *mockRegistry) LoadImages(name, digest string) {

}
