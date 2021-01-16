package sql

import (
	"context"

	"api.fabl.app/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accountRepo struct {
	db *sqlx.DB
}

func (r *accountRepo) Get(ctx context.Context, id uuid.UUID) (*repository.Account, error) {
	var acc struct {
		ID             uuid.UUID `db:"id"`
		HashedPassword []byte    `db:"hashed_password"`
		Nickname       string    `db:"nickname"`
	}
	err := r.db.GetContext(ctx, &acc, `
		SELECT
			id, hashed_password, nickname
		FROM
			account
		WHERE
			id = $1;`,
		id,
	)
	return &repository.Account{
		ID:             acc.ID,
		HashedPassword: acc.HashedPassword,
		Nickname:       acc.Nickname,
	}, err
}

func (r *accountRepo) FromToken(ctx context.Context, token string) (*repository.Account, error) {
	return nil, status.Errorf(codes.Unimplemented, "sql.accountRepo method FromToken not implemented")
}
