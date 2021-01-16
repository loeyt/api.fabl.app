package sql

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"

	"api.fabl.app/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type itemRepo struct {
	db *sqlx.DB
}

func (r *itemRepo) Get(ctx context.Context, id ulid.ULID) (*repository.Item, error) {
	v := struct {
		Data      []byte    `db:"item_data"`
		AccountID uuid.UUID `db:"account_id"`
	}{}
	err := r.db.GetContext(ctx, &v, `
		SELECT
			item_data, account_id
		FROM
			item
			INNER JOIN item_data ON item.sum256 = item_data.sum256
		WHERE
			id = $1;
	`, id)
	return &repository.Item{
		ULID:      id,
		AccountID: v.AccountID,
		TimeMs:    id.Time(),
		Data:      v.Data,
	}, err
}

func (r *itemRepo) Create(ctx context.Context, item *repository.Item) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	var v []byte
	err = tx.GetContext(ctx, &v, `
		INSERT
		INTO
			item_data (item_data)
		VALUES
			($1)
		ON CONFLICT
		DO
			NOTHING
		RETURNING
			sum256;`,
		item.Data,
	)
	if err == sql.ErrNoRows {
		sum256 := sha256.Sum256(item.Data)
		item.Sum256 = &sum256
	} else if err != nil {
		return err
	} else {
		item.Sum256 = new([32]byte)
		copy((*item.Sum256)[:], v)
	}
	if item.TimeMs == 0 {
		item.TimeMs = ulid.Now()
	}
	item.ULID, err = ulid.New(item.TimeMs, bytes.NewReader(item.Sum256[:]))
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT
		INTO
			item (id, sum256, account_id)
		VALUES
			($1, $2, $3);`,
		item.ULID[:], item.Sum256[:], item.AccountID,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *itemRepo) List(ctx context.Context, accountID uuid.UUID) ([]*repository.Item, error) {
	var v []*struct {
		ULID      ulid.ULID `db:"id"`
		AccountID uuid.UUID `db:"account_id"`
		Sum256    []byte    `db:"sum256"`
	}
	err := r.db.SelectContext(ctx, &v,
		`SELECT id, account_id, sum256 FROM item WHERE account_id = $1 ORDER BY id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	items := make([]*repository.Item, len(v))
	for i, item := range v {
		items[i] = &repository.Item{
			ULID:      item.ULID,
			AccountID: item.AccountID,
			Sum256:    new([32]byte),
		}
		copy(items[i].Sum256[:], item.Sum256)
	}
	return items, err
}

func (r *itemRepo) GetData(ctx context.Context, sum256 [32]byte) ([]byte, error) {
	return nil, status.Errorf(codes.Unimplemented, "sql.itemRepo method GetData not implemented")
}
