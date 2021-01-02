package service

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
)

var errScanValue = errors.New("source must be a byte slice")

type sha256Sum [32]byte

// Scan implements the sql.Scanner interface. It only supports scanning byte
// slices.
func (s *sha256Sum) Scan(src interface{}) error {
	switch x := src.(type) {
	case nil:
		return nil
	case []byte:
		if len(x) != 32 {
			return fmt.Errorf("expected 32 bytes, got %d", len(x))
		}
		copy(s[:], x)
		return nil
	}

	return errScanValue
}

type Item struct {
	ULID ulid.ULID `db:"id"`
	Sum  sha256Sum `db:"sum"`
}

func extractItemData(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, errors.New("empty import string")
	}
	firstRune, s := s[0], s[1:]
	if firstRune != '0' {
		return nil, fmt.Errorf("unrecognized first character: 0x%x", firstRune)
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	r, err := zlib.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ItemStore is the storage type used by ItemService.
type ItemStore interface {
	GetData(ctx context.Context, id ulid.ULID) ([]byte, error)
	Create(ctx context.Context, data []byte, timeMs uint64) (*Item, error)
	List(ctx context.Context) ([]*Item, error)
}

type sqlxItemStore struct {
	db *sqlx.DB
}

// NewItemStore creates an itemStore based on a sqlx DB.
func NewItemStore(db *sqlx.DB) ItemStore {
	return &sqlxItemStore{
		db: db,
	}
}

func (s *sqlxItemStore) Create(ctx context.Context, data []byte, timeMs uint64) (*Item, error) {
	item := &Item{
		Sum: sha256.Sum256(data),
	}
	if timeMs == 0 {
		timeMs = ulid.Now()
	}
	var err error
	item.ULID, err = ulid.New(timeMs, bytes.NewReader(item.Sum[:]))
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO item_data (data) VALUES ($1) ON CONFLICT DO NOTHING`,
		data,
	)
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO item (id, sum) VALUES ($1, $2)`,
		item.ULID[:], item.Sum[:],
	)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return item, nil
}

func (s *sqlxItemStore) GetData(ctx context.Context, id ulid.ULID) ([]byte, error) {
	v := struct {
		Data []byte `db:"data"`
	}{}
	err := s.db.GetContext(ctx, &v,
		`SELECT data FROM item_data WHERE sum = (SELECT sum FROM item WHERE id = $1)`,
		id[:],
	)
	return v.Data, err
}

func (s *sqlxItemStore) List(ctx context.Context) ([]*Item, error) {
	var items []*Item
	err := s.db.SelectContext(ctx, &items, `SELECT id, sum FROM item ORDER BY id`)
	return items, err
}
