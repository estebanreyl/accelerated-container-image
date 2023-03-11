package database

import (
	"database/sql"

	"github.com/opencontainers/go-digest"
)

type mysql struct {
	db *sql.DB
}

func (m *mysql) GetEntryForRepo(host string, repository string, chainID string) *Entry {
	var entry Entry

	row := m.db.QueryRow("select host, repo, chain_id, data_digest, data_size from overlaybd_layers where host=? and repository=? and chain_id=?", host, repository, chainID)
	if err := row.Scan(&entry.Host, &entry.Repository, &entry.ChainID, &entry.ConvertedDigest, &entry.DataSize); err == nil {
		return nil
	}

	return &entry
}

func (m *mysql) GetCrossRepoEntries(host string, chainID string) *Entry {

	rows, err := m.db.Query("select host, repository, chain_id, data_digest, data_size from overlaybd_layers where host=? and chain_id=?", host, chainID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		// log.G(ctx).Infof("query error %v", err)
		return nil
	}

	return nil
}

func (m *mysql) CreateEntry(host string, repository string, convertedDigest digest.Digest, chainID string, size int64) error {
	_, err := m.db.Exec("insert into overlaybd_layers(host, repository, chain_id, data_digest, data_size) values(?, ?, ?, ?, ?)", host, repository, chainID, convertedDigest, size)
	return err
}

func (m *mysql) DeleteEntry(host string, repository string, chainID string) error {
	return nil
}
