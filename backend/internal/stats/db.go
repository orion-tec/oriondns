package stats

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/orion-tec/oriondns/db"
)

type statsDB struct {
	db *db.DB
}

type DB interface {
	Insert(ctx context.Context, t time.Time, domain, domainType string) error
	GetMostUsedDomains(ctx context.Context, from, to time.Time, limit int) ([]MostUsedDomainResponse, error)
	GetUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, domains []string) ([]MostUsedDomainResponse, error)
	GetMostUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time) ([]MostUsedDomainResponse, error)
}

func New(db *db.DB) DB {
	return &statsDB{db}
}

func (s *statsDB) GetMostUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time) ([]MostUsedDomainResponse, error) {
	domains, err := s.GetMostUsedDomains(ctx, from, to, 10)
	if err != nil {
		return nil, err
	}

	domainsNames := make([]string, len(domains))
	for i, domain := range domains {
		domainsNames[i] = domain.Domain
	}

	return s.GetUsedDomainsByTimeAggregation(ctx, from, to, domainsNames)
}

func (s *statsDB) GetUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, domains []string) ([]MostUsedDomainResponse, error) {
	query := `
			SELECT domain, SUM(count) as count
      FROM stats_aggregated
      WHERE time >= $1 AND time <= $2 AND domain = ANY($3)
      GROUP BY domain
      ORDER BY count DESC
	`
	rows, err := s.db.Query(ctx, query, from, to, domains)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []MostUsedDomainResponse
	for rows.Next() {
		var r MostUsedDomainResponse
		if err := rows.Scan(&r.Domain, &r.Count); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func (s *statsDB) GetMostUsedDomains(ctx context.Context, from, to time.Time, limit int) ([]MostUsedDomainResponse, error) {
	query := `
			SELECT domain, SUM(count) as count
      FROM stats_aggregated
      WHERE time >= $1 AND time <= $2 AND q_type != 'PTR'
      GROUP BY domain
      ORDER BY count DESC LIMIT $3
	`
	rows, err := s.db.Query(ctx, query, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[MostUsedDomainResponse])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *statsDB) Insert(ctx context.Context, t time.Time, domain, domainType string) error {
	// Truncate the time to the minute
	// This is to ensure that we can aggregate stats by 10 minutes
	minute := t.Minute()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute-t.Minute()%10, 0, 0, t.Location())

	_, err := s.db.Exec(ctx, `
		INSERT INTO stats_aggregated (time, domain, count, q_type)
			VALUES ($1, $2, 1, $3) ON CONFLICT(time, domain) DO
		UPDATE SET count = stats_aggregated.count + 1, updated_at = NOW()
	`, newTime, domain, domainType)
	if err != nil {
		return err
	}

	return nil
}
