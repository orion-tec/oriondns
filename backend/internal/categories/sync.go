package categories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/orion-tec/oriondns/internal/ai"
	"github.com/orion-tec/oriondns/internal/domains"
	"go.uber.org/fx"
)

type syncer struct {
	ai         ai.AI
	categoryDB DB
	domainDB   domains.DB
}

type CategoryAIAnswer struct {
	Category []string `json:"category"`
}

type Syncer interface {
	Sync() error
}

func NewSyncer(lc fx.Lifecycle, ai ai.AI, categoryDB DB, domainsDB domains.DB) Syncer {
	s := &syncer{ai, categoryDB, domainsDB}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Println("Starting syncer")
				for {
					err := s.Sync()
					if err != nil {
						log.Printf("Error on syncer: %s\n", err)
					}
					time.Sleep(5 * time.Second)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return s
}

func (s *syncer) Sync() error {
	domains, err := s.domainDB.GetDomainsWithoutCategory(context.Background())
	if err != nil {
		return err
	}

	log.Printf("Found %d domains without category\n", len(domains))
	for _, domain := range domains {
		log.Printf("Processing domain %s\n", domain.Domain)
		query := fmt.Sprintf(`
			Considering domain %s, which content category you thing would be a good fit for it?
			Answer me only with a json with a 'category' key and an array with sanitized categories all in lower case.
		`, domain.Domain)
		answer, err := s.ai.Query(query)
		if errors.Is(err, ai.ErrRateLimit) {
			log.Printf("Rate limit exceeded, waiting 10 minutes\n")
			time.Sleep(10 * time.Minute)
			continue
		}
		if err != nil {
			log.Printf("Error on query AI for domain %s: %s\n", domain.Domain, err)
			continue
		}

		c := CategoryAIAnswer{}
		err = json.Unmarshal([]byte(answer), &c)
		if err != nil {
			log.Printf("Error on decode data from AI for domain %s: %s\n", domain.Domain, err)
			continue
		}

		err = s.categoryDB.Insert(context.Background(), domain.Domain, c.Category)
		if err != nil {
			log.Printf("Error on insert data on db for domain %s: %s\n", domain.Domain, err)
			continue
		}

		log.Printf("Domain %s categorized as %v\n", domain.Domain, c.Category)
		time.Sleep(10 * time.Second)
	}

	return nil
}
