package database

import (
	"context"

	"github.com/opencontainers/go-digest"
)

type ConversionDatabase interface {
	GetEntryForRepo(ctx context.Context, host string, repository string, chainID string) *Entry
	GetCrossRepoEntries(ctx context.Context, host string, chainID string) []*Entry
	CreateEntry(ctx context.Context, host string, repository string, convertedDigest digest.Digest, chainID string, size int64) error
	DeleteEntry(ctx context.Context, host string, repository string, chainID string) error
}

type Entry struct {
	ConvertedDigest digest.Digest
	DataSize        int64
	Repository      string
	ChainID         string
	Host            string
}
