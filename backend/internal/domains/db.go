package domains

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/orion-tec/oriondns/db"
)

type domainsDB struct {
	db *db.DB
}

type DB interface {
	Insert(ctx context.Context, domain string) error
	GetAll(ctx context.Context) ([]Domain, error)
	GetByDomain(ctx context.Context, domain string) (*Domain, error)
	GetDomainsWithoutCategory(ctx context.Context) ([]Domain, error)
}

func New(db *db.DB) DB {
	return &domainsDB{db}
}

func (b *domainsDB) GetDomainsWithoutCategory(ctx context.Context) ([]Domain, error) {
	row, err := b.db.Query(ctx, `
		SELECT d.*
		FROM domains d
						 LEFT JOIN public.domain_categories dc ON d.domain = dc.domain
		WHERE category IS NULL
		ORDER BY used_count DESC
		LIMIT 10
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

func (b *domainsDB) GetByDomain(ctx context.Context, domain string) (*Domain, error) {
	row, err := b.db.Query(ctx, `
		SELECT domain, used_count, created_at, updated_at
		FROM domains
		WHERE domain = $1
	`, domain)
	if err != nil {
		return nil, err
	}

	domainStr, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[Domain])
	if err != nil {
		return nil, err
	}

	return &domainStr, nil
}

func (b *domainsDB) Insert(ctx context.Context, domain string) error {
	// Try to insert the domain with used_count = 1
	_, err := b.db.Exec(ctx, `
		INSERT INTO domains (domain, used_count)
		VALUES ($1, 1)
	`, domain)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation
			// If domain already exists, increment used_count
			_, updateErr := b.db.Exec(ctx, `
				UPDATE domains
				SET used_count = used_count + 1, updated_at = now()
				WHERE domain = $1
			`, domain)
			return updateErr
		}
		return err // Return other errors
	}

	return nil
}

func (b *domainsDB) GetAll(ctx context.Context) ([]Domain, error) {
	rows, err := b.db.Query(ctx, `
		SELECT domain, created_at, updated_at, used_count
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
