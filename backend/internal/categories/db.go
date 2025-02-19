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
	Insert(ctx context.Context, domain string, categories []string) error
}

func New(db *db.DB) DB {
	return &categoriesDB{db}
}

func (b *categoriesDB) Insert(ctx context.Context, domain string, categories []string) error {
	tx, err := b.db.Begin(ctx)
	if err != nil {
		return err
	}

	for _, category := range categories {
		_, err := tx.Exec(ctx, `
			INSERT INTO domain_categories (domain, category)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, domain, category)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (b *categoriesDB) GetAll(ctx context.Context) ([]Category, error) {
	rows, err := b.db.Query(ctx, `
		SELECT DISTINCT category
		FROM domain_categories
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
