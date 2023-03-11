package database

import "github.com/opencontainers/go-digest"

type ConversionDatabase interface {
	GetEntryForRepo(host string, repository string, chainID string) *Entry
	GetCrossRepoEntry(host string, chainID string) *Entry
	CreateEntry(host string, repository string, convertedDigest digest.Digest, chainID string) error
	DeleteEntry(host string, repository string, chainID string) error
}

type Entry struct {
	ConvertedDigest digest.Digest
	DataSize        int64
	Repository      string
	ChainID         string
	Host            string
}
