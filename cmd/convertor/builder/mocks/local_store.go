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
package mocks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/containers/image/v5/directory"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/pkg/blobinfocache/none"
	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// RESOLVER
type MockLocalResolver struct {
	localReg *localRegistry
}

func (r *MockLocalResolver) Resolve(ctx context.Context, ref string) (string, v1.Descriptor, error) {
	desc, err := r.localReg.findRef(ctx, ref)
	if err != nil {
		return "", v1.Descriptor{}, err
	}
	return "", desc, nil
}

func (r *MockLocalResolver) Fetcher(ctx context.Context, ref string) (remotes.Fetcher, error) {
	return &MockLocalFetcher{
		localReg: r.localReg,
	}, nil
}

func (r *MockLocalResolver) Pusher(ctx context.Context, ref string) (remotes.Pusher, error) {
	return &MockLocalPusher{
		localReg: r.localReg,
	}, nil
}

// FETCHER
type MockLocalFetcher struct {
	localReg *localRegistry
}

func (f *MockLocalFetcher) Fetch(ctx context.Context, desc v1.Descriptor) (io.ReadCloser, error) {
	switch desc.MediaType {
	// Get Manifest
	case images.MediaTypeDockerSchema2Manifest,
		images.MediaTypeDockerSchema2ManifestList,
		images.MediaTypeDockerSchema1Manifest,
		v1.MediaTypeImageManifest,
		v1.MediaTypeImageIndex:

		content, err := f.localReg.GetManifest(ctx, desc)

		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(content)
		return ioutil.NopCloser(reader), nil
	}
	return f.localReg.GetBlob(ctx, desc)
}

// PUSHER
type MockLocalPusher struct {
	localReg *localRegistry
}

// Not used by overlaybd conversion
func (p MockLocalPusher) Writer(ctx context.Context, opts ...content.WriterOpt) (content.Writer, error) {
	return nil, errors.New("Not implemented")
}

func (p MockLocalPusher) Push(ctx context.Context, desc v1.Descriptor) (content.Writer, error) {
	return nil, errors.New("Not implemented")
}

type localRegistry struct {
	images []types.ImageReference
}

func newLocalRegistry(ctx context.Context) (*localRegistry, error) {
	refs, err := findImagesFromSource(ctx)
	if err != nil {
		return nil, err
	}

	return &localRegistry{
		images: refs,
	}, nil
}

func (l *localRegistry) GetBlob(ctx context.Context, desc v1.Descriptor) (io.ReadCloser, error) {
	ref := l.images[0]

	img, err := ref.NewImageSource(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	stream, _, err := img.GetBlob(ctx, manifest.BlobInfoFromOCI1Descriptor(desc), none.NoCache)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (l *localRegistry) GetManifest(ctx context.Context, desc v1.Descriptor) ([]byte, error) {
	ref := l.images[0]

	img, err := ref.NewImage(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	mnfst, _, err := img.Manifest(ctx)
	if err != nil {
		return nil, err
	}

	return mnfst, nil
}

func (l *localRegistry) findRef(ctx context.Context, refStr string) (v1.Descriptor, error) {
	ref := l.images[0]

	if ref.DockerReference().Name() != refStr {
		return v1.Descriptor{}, fmt.Errorf("refStr is improper found %s, wanted %s", ref.DockerReference().Name(), refStr)
	}

	img, err := ref.NewImage(ctx, nil)
	if err != nil {
		return v1.Descriptor{}, err
	}
	defer img.Close()

	content, _, err := img.Manifest(ctx)
	if err != nil {
		return v1.Descriptor{}, err
	}

	parsedManifest := v1.Manifest{}
	err = json.Unmarshal(content, &parsedManifest)
	if err != nil {
		return v1.Descriptor{}, err
	}

	return v1.Descriptor{
		MediaType: parsedManifest.MediaType,
		Digest:    digest.FromBytes(content),
		Size:      int64(len(content)),
	}, nil
}

func findImagesFromSource(ctx context.Context) ([]types.ImageReference, error) {

	var sourceReferences []types.ImageReference
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(fmt.Sprintf("%s/registry", cwd), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "manifest.json" {
			dirname := filepath.Dir(path)
			ref, err := directory.Transport.ParseReference(dirname)
			if err != nil {
				return err
			}
			sourceReferences = append(sourceReferences, ref)
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return sourceReferences, nil
}
