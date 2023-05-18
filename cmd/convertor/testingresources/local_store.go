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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/reference"
	"github.com/containerd/containerd/remotes"
	"github.com/containers/image/v5/directory"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/pkg/blobinfocache/none"
	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type internalImage struct {
	descriptor v1.Descriptor
	ref        types.ImageReference
}

type internalRegistry map[string]map[string]*internalImage

// RESOLVER
type MockLocalResolver struct {
	localReg *localRegistry
}

func NewMockLocalResolver(ctx context.Context, localRegistryPath string) (*MockLocalResolver, error) {
	reg, err := newLocalRegistry(ctx, localRegistryPath)
	if err != nil {
		return nil, err
	}

	return &MockLocalResolver{
		localReg: reg,
	}, nil
}

func (r *MockLocalResolver) Resolve(ctx context.Context, ref string) (string, v1.Descriptor, error) {
	desc, err := r.localReg.findDescriptorFromRef(ctx, ref)
	if err != nil {
		return "", v1.Descriptor{}, err
	}
	return "", desc, nil
}

func (r *MockLocalResolver) Fetcher(ctx context.Context, ref string) (remotes.Fetcher, error) {
	_, repository, _, err := parseRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	return &MockLocalFetcher{
		localReg:   r.localReg,
		repository: repository,
	}, nil
}

func (r *MockLocalResolver) Pusher(ctx context.Context, ref string) (remotes.Pusher, error) {
	_, repository, _, err := parseRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	return &MockLocalPusher{
		localReg:   r.localReg,
		repository: repository,
	}, nil
}

// FETCHER
type MockLocalFetcher struct {
	localReg   *localRegistry
	repository string
}

func (f *MockLocalFetcher) Fetch(ctx context.Context, desc v1.Descriptor) (io.ReadCloser, error) {
	switch desc.MediaType {
	// Get Manifest
	case images.MediaTypeDockerSchema2Manifest,
		images.MediaTypeDockerSchema2ManifestList,
		images.MediaTypeDockerSchema1Manifest,
		v1.MediaTypeImageManifest,
		v1.MediaTypeImageIndex:

		content, err := f.localReg.GetManifest(ctx, f.repository, desc)

		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(content)
		return io.NopCloser(reader), nil
	}
	return f.localReg.GetBlob(ctx, f.repository, desc)
}

// PUSHER
type MockLocalPusher struct {
	localReg   *localRegistry
	repository string
}

// Not used by overlaybd conversion
func (p MockLocalPusher) Writer(ctx context.Context, opts ...content.WriterOpt) (content.Writer, error) {
	return nil, errors.New("Not implemented")
}

func (p MockLocalPusher) Push(ctx context.Context, desc v1.Descriptor) (content.Writer, error) {
	return nil, errors.New("Not implemented")
}

type localRegistry struct {
	registryInternal internalRegistry
}

func newLocalRegistry(ctx context.Context, localRegistryPath string) (*localRegistry, error) {
	registryInternal, err := findImagesFromSource(ctx, localRegistryPath)
	if err != nil {
		return nil, err
	}

	return &localRegistry{
		registryInternal: registryInternal,
	}, nil
}

func (l *localRegistry) GetBlob(ctx context.Context, repository string, desc v1.Descriptor) (io.ReadCloser, error) {
	findBlobGivenRef := func(ctx context.Context, ref types.ImageReference, desc v1.Descriptor) (io.ReadCloser, error) {
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

	// Foreach to go over all images in a repository
	for _, image := range l.registryInternal[repository] {
		stream, err := findBlobGivenRef(ctx, image.ref, desc)
		if err != nil {
			continue // Try all refs for blob in the repository
		}
		return stream, nil
	}
	return nil, errors.New("Blob not found")
}

func (l *localRegistry) GetManifest(ctx context.Context, repository string, desc v1.Descriptor) ([]byte, error) {
	// If we are looking for the main manifest (In case of manifest list)

	if _, ok := l.registryInternal[repository]; !ok {
		return nil, errors.New("Repository not found")
	}

	// Exact image
	internalImg, ok := l.registryInternal[repository][desc.Digest.String()]
	if ok {
		img, err := internalImg.ref.NewImage(ctx, nil)
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
	// Sub manifest in a manifest list
	for _, subImg := range l.registryInternal[repository] {
		mnfstSrc, err := subImg.ref.NewImageSource(ctx, nil)
		if err != nil {
			continue
		}
		mnfst, _, err := mnfstSrc.GetManifest(ctx, &desc.Digest)
		if err != nil {
			return nil, err
		}
		return mnfst, nil
	}
	return nil, errors.New("Manifest not found")
}

func (l *localRegistry) findDescriptorFromRef(ctx context.Context, refStr string) (v1.Descriptor, error) {
	_, repository, object, err := parseRef(ctx, refStr)
	if err != nil {
		return v1.Descriptor{}, err
	}

	if _, ok := l.registryInternal[repository]; !ok {
		return v1.Descriptor{}, errors.New("Repository not found")
	}
	img, ok := l.registryInternal[repository][object]
	if !ok {
		return v1.Descriptor{}, errors.New("Reference not found")
	}
	return img.descriptor, nil
}

func getDesc(ctx context.Context, ref types.ImageReference) (v1.Descriptor, error) {

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

func findImagesFromSource(ctx context.Context, localRegistryPath string) (internalRegistry, error) {

	var sourceReferences []types.ImageReference

	// Repository -> Image -> (Ref + Descriptor)
	registryInternal := make(internalRegistry)
	re := regexp.MustCompile(`^.+testingresources/mocks/registry/(.+):(.+)`) // Regex to get repository and reference name

	err := filepath.WalkDir(localRegistryPath, func(path string, d fs.DirEntry, err error) error {
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
			desc, err := getDesc(ctx, ref)
			if err != nil {
				return err
			}

			matches := re.FindStringSubmatch(dirname)
			repo := matches[1]   // Get repository name
			imgRef := matches[2] // Get reference name

			if _, ok := registryInternal[repo]; !ok {
				registryInternal[repo] = make(map[string]*internalImage) // Create repository
			}

			img := internalImage{
				descriptor: desc,
				ref:        ref,
			}

			registryInternal[repo][desc.Digest.String()] = &img
			registryInternal[repo][imgRef] = &img

			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return registryInternal, nil
}

func parseRef(ctx context.Context, ref string) (string, string, string, error) {
	refspec, err := reference.Parse(ref)
	if err != nil {
		return "", "", "", err
	}
	host := refspec.Hostname()
	repository := strings.TrimPrefix(refspec.Locator, host+"/")
	object := refspec.Object
	return host, repository, object, nil
}
