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
	"database/sql"
	"fmt"

	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"

	"github.com/opencontainers/go-digest"
)

type sqldb struct {
	db *sql.DB
}

func NewSqlDB(db *sql.DB) ConversionDatabase {
	return &sqldb{
		db: db,
	}
}

func (m *sqldb) GetEntryForRepo(ctx context.Context, host string, repository string, chainID string) *Entry {
	var entry Entry

	row := m.db.QueryRowContext(ctx, "select host, repo, chain_id, data_digest, data_size from overlaybd_layers where host=? and repo=? and chain_id=? where type=?", host, repository, chainID, "layer")
	if err := row.Scan(&entry.Host, &entry.Repository, &entry.ChainID, &entry.ConvertedDigest, &entry.DataSize); err != nil {
		return nil
	}

	return &entry
}

func (m *sqldb) GetEntryForManifest(ctx context.Context, host string, repository string, manifest digest.Digest) *Entry {
	var entry Entry
	row := m.db.QueryRowContext(ctx, "select host, repo, chain_id, data_digest, data_size, manifest_digest from overlaybd_layers where host=? and repo=? and manifest_digest=? where type=?", host, repository, manifest, "manifest")
	if err := row.Scan(&entry.Host, &entry.Repository, &entry.ChainID, &entry.ConvertedDigest, &entry.DataSize); err != nil {
		return nil
	}

	return &entry
}

func (m *sqldb) GetCrossRepoEntries(ctx context.Context, host string, chainID string) []*Entry {
	rows, err := m.db.QueryContext(ctx, "select host, repo, chain_id, data_digest, data_size from overlaybd_layers where host=? and chain_id=? where type=?", host, chainID, "layer")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.G(ctx).Infof("query error %v", err)
		return nil
	}
	var entries []*Entry
	for rows.Next() {
		var entry Entry
		err = rows.Scan(&entry.Host, &entry.Repository, &entry.ChainID, &entry.ConvertedDigest, &entry.DataSize)
		if err != nil {
			continue
		}
		entries = append(entries, &entry)
	}

	return entries
}

func (m *sqldb) CreateEntry(ctx context.Context, host string, repository string, convertedDigest digest.Digest, chainID string, manifestDigest digest.Digest, size int64, entryType string) error {
	if entryType != "manifest" && entryType != "layer" {
		return errors.Errorf("invalid entry type %s", entryType)
	}
	_, err := m.db.ExecContext(ctx, "insert into overlaybd_layers(host, repo, chain_id, data_digest, data_size, manifest_digest, type) values(?, ?, ?, ?, ?, ?)", host, repository, chainID, convertedDigest, size, manifestDigest, entryType)
	return err
}

func (m *sqldb) DeleteEntry(ctx context.Context, host string, repository string, chainID string, manifest digest.Digest, entryType string) error {
	var err error
	if entryType == "manifest" {
		_, err = m.db.Exec("delete from overlaybd_layers where host=? and repo=? and chain_id=? and type=?", host, repository, chainID, entryType)
	} else if entryType == "layer" {
		_, err = m.db.Exec("delete from overlaybd_layers where host=? and repo=? and manifest_digest=? and type=?", host, repository, chainID, manifest, entryType)
	} else {
		return errors.Errorf("invalid entry type %s", entryType)
	}
	if err != nil {
		return fmt.Errorf("failed to remove invalid record in db: %w", err)
	}
	return nil
}
