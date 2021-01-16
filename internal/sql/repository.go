package sql

import (
	"api.fabl.app/internal/repository"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Account repository.AccountRepository
	Item    repository.ItemRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Item:    &itemRepo{db: db},
		Account: &accountRepo{db: db},
	}
}
