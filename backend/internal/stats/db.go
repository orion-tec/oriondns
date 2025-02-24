package stats

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/huandu/go-sqlbuilder"

	"github.com/orion-tec/oriondns/db"
)

type statsDB struct {
	db *db.DB
}

type DB interface {
	Insert(ctx context.Context, t time.Time, domain, domainType string) error
	GetMostUsedDomains(ctx context.Context, from, to time.Time, categories []string, limit int) ([]MostUsedDomainResponse, error)
	GetUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, domains []string) ([]MostUsedDomainResponse, error)
	GetMostUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, categories []string) ([]MostUsedDomainResponse, error)
}

func New(db *db.DB) DB {
	return &statsDB{db}
}

func (s *statsDB) GetMostUsedDomainsByTimeAggregation(ctx context.Context, from, to time.Time, categories []string) ([]MostUsedDomainResponse, error) {
	domains, err := s.GetMostUsedDomains(ctx, from, to, categories, 10)
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

func (s *statsDB) GetMostUsedDomains(ctx context.Context, from, to time.Time, categories []string, limit int) ([]MostUsedDomainResponse, error) {
	var inData string
	if len(categories) > 0 {
		inData = strings.Join(categories, ",")
	} else {
		inData = "SELECT DISTINCT category FROM domain_categories"
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	cond := sqlbuilder.NewCond()
	where := sqlbuilder.NewWhereClause().
		AddWhereExpr(cond.Args,
			cond.And(
				cond.GreaterEqualThan("time", from),
				cond.LessEqualThan("time", to),
				cond.NotEqual("q_type", "PTR"),
				cond.In("dc.category", inData),
			),
		)

	sb.Select("sa.domain", "SUM(count) as count").
		From("stats_aggregated sa").
		Join("domain_categories dc on dc.domain = sa.domain").
		GroupBy("sa.domain").
		OrderBy("count DESC").
		Limit(limit)
	sb.WhereClause = where

	query, args := sb.Build()
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

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
