package repository

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID             uuid.UUID
	HashedPassword []byte
	Nickname       string
}

type AccountRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*Account, error)
	FromToken(ctx context.Context, token string) (*Account, error)
}

func (u *Account) CheckPassword(password []byte) error {
	if u.HashedPassword == nil || password == nil {
		return bcrypt.ErrMismatchedHashAndPassword
	}
	return bcrypt.CompareHashAndPassword(u.HashedPassword, password)
}
