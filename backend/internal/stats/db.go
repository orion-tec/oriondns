package stats

import (
	"context"
	"time"

	"github.com/orion-tec/oriondns/db"
)

type statsDB struct {
	db *db.DB
}

type DB interface {
	Insert(ctx context.Context, t time.Time, domain string) error
}

func New(db *db.DB) DB {
	return &statsDB{db}
}

func (s *statsDB) Insert(ctx context.Context, t time.Time, domain string) error {
	// Truncate the time to the minute
	// This is to ensure that we can aggregate stats by 10 minutes
	m := t.Minute()
	newTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m-t.Minute()%10, 0, 0, t.Location())

	_, err := s.db.Exec(ctx, `
		INSERT INTO stats_aggregated (time, domain, count) 
			VALUES ($1, $2, 1) ON CONFLICT(time, domain) DO 
		UPDATE SET count = stats_aggregated.count + 1, updated_at = NOW()
	`, newTime, domain)
	if err != nil {
		return err
	}

	return nil
}
