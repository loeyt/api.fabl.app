package repository

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type Item struct {
	ULID      ulid.ULID
	AccountID uuid.UUID
	TimeMs    uint64
	Data      []byte
	// Sum256 should be present when queried without Data.
	Sum256 *[32]byte
}

type ItemRepository interface {
	Get(ctx context.Context, id ulid.ULID) (*Item, error)
	GetData(ctx context.Context, sum256 [32]byte) ([]byte, error)
	Create(ctx context.Context, item *Item) error
	List(ctx context.Context, accountID uuid.UUID) ([]*Item, error)
}

func (i *Item) Import(s string) error {
	if len(s) == 0 {
		return errors.New("empty import string")
	}
	firstRune, s := s[0], s[1:]
	if firstRune != '0' {
		return fmt.Errorf("unrecognized first character: 0x%x", firstRune)
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	r, err := zlib.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, r)
	if err != nil {
		return err
	}
	i.Data = buf.Bytes()
	return nil
}

func (i *Item) Export() (string, error) {
	s := new(strings.Builder)
	s.WriteString("0")
	enc := base64.NewEncoder(base64.StdEncoding, s)
	w, err := zlib.NewWriterLevel(enc, zlib.BestCompression)
	if err != nil {
		return "", err
	}
	r := bytes.NewReader(i.Data)
	_, err = io.Copy(w, r)
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	err = enc.Close()
	if err != nil {
		return "", err
	}
	return s.String(), nil

}
