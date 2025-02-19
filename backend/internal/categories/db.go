package categories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/orion-tec/oriondns/db"
)

type categoriesDB struct {
	db *db.DB
}

type DB interface {
	GetAll(ctx context.Context) ([]Category, error)
}

func New(db *db.DB) DB {
	return &categoriesDB{db}
}

func (b *categoriesDB) GetAll(ctx context.Context) ([]Category, error) {
	rows, err := b.db.Query(ctx, `
		SELECT id, name
		FROM categories
	`)
	if err != nil {
		return nil, err
	}

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[Category])
	if err != nil {
		return nil, err
	}

	return categories, nil
}
