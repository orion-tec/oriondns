package domains

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/orion-tec/oriondns/db"
)

type domainsDB struct {
	db *db.DB
}

type DB interface {
	Insert(ctx context.Context, domain string) error
	GetAll(ctx context.Context) ([]Domain, error)
	GetByName(ctx context.Context, name string) (*Domain, error)
	GetDomainsWithoutCategory(ctx context.Context) ([]Domain, error)
}

func New(db *db.DB) DB {
	return &domainsDB{db}
}

func (b *domainsDB) GetDomainsWithoutCategory(ctx context.Context) ([]Domain, error) {
	row, err := b.db.Query(ctx, `
		SELECT d.name as name
		FROM domains d
						 left join domain_categories dc
											 on d.name = dc.domain_name
						 left join categories c
											 on dc.category_id = c.id
		where c.id is null;
	`)
	if err != nil {
		return nil, err
	}

	domains, err := pgx.CollectRows(row, pgx.RowToStructByName[Domain])
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (b *domainsDB) GetByName(ctx context.Context, name string) (*Domain, error) {
	row, err := b.db.Query(ctx, `
		SELECT name
		FROM domains
		WHERE name = $1
	`, name)
	if err != nil {
		return nil, err
	}

	domain, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[Domain])
	if err != nil {
		return nil, err
	}

	return &domain, nil
}

func (b *domainsDB) Insert(ctx context.Context, domain string) error {
	_, err := b.db.Exec(ctx, `
		INSERT INTO domains (domain)
			VALUES ($1)
		ON CONFLICT DO UPDATE SET domain = $1, updated_at = now(), used_count = used_count + 1
	`, domain)
	if err != nil {
		return err
	}

	return nil
}

func (b *domainsDB) GetAll(ctx context.Context) ([]Domain, error) {
	rows, err := b.db.Query(ctx, `
		SELECT name
		FROM domains
	`)
	if err != nil {
		return nil, err
	}

	domains, err := pgx.CollectRows(rows, pgx.RowToStructByName[Domain])
	if err != nil {
		return nil, err
	}

	return domains, nil
}
