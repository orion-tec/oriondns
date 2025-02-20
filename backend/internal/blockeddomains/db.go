package blockeddomains

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/orion-tec/oriondns/db"
)

type blockedDomainsDB struct {
	db *db.DB
}

type DB interface {
	Insert(ctx context.Context, domain string, recursive bool) error
	GetAll(ctx context.Context) ([]BlockedDomain, error)
}

func New(db *db.DB) DB {
	return &blockedDomainsDB{db}
}

func (b *blockedDomainsDB) Insert(ctx context.Context, domain string, recursive bool) error {
	_, err := b.db.Exec(ctx, `
		INSERT INTO blocked_domains (domain, recursive) 
			VALUES ($1, $2)
	`, domain, recursive)
	if err != nil {
		return err
	}

	return nil
}

func (b *blockedDomainsDB) GetAll(ctx context.Context) ([]BlockedDomain, error) {
	rows, err := b.db.Query(ctx, `
		SELECT id, domain, recursive, created_at, updated_at, deleted_at
		FROM blocked_domains
	`)
	if err != nil {
		return nil, err
	}

	blockedDomains, err := pgx.CollectRows(rows, pgx.RowToStructByName[BlockedDomain])
	if err != nil {
		return nil, err
	}

	return blockedDomains, nil
}
