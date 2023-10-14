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

package database

import (
	"context"

	"github.com/opencontainers/go-digest"
)

type ConversionDatabase interface {
	GetEntryForManifest(ctx context.Context, host string, repository string, manifest digest.Digest) *Entry
	GetEntryForRepo(ctx context.Context, host string, repository string, chainID string) *Entry
	GetCrossRepoEntries(ctx context.Context, host string, chainID string) []*Entry
	CreateEntry(ctx context.Context, host string, repository string, convertedDigest digest.Digest, chainID string, manifestDigest digest.Digest, size int64, entryType string) error
	DeleteEntry(ctx context.Context, host string, repository string, chainID string, manifest digest.Digest, entryType string) error
}

type Entry struct {
	ManifestDigest  digest.Digest // Present when entryType is manifest
	ConvertedDigest digest.Digest
	DataSize        int64
	Repository      string
	ChainID         string // Present when entryType is layer
	Host            string
	EntryType       string // One of "manifest" or "layer"
}
