package categories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.uber.org/fx"

	"github.com/orion-tec/oriondns/internal/ai"
	"github.com/orion-tec/oriondns/internal/domains"
)

type syncer struct {
	ai         ai.AI
	categoryDB DB
	domainDB   domains.DB

	secondsToSleep int
}

type CategoryAIAnswer struct {
	Category []string `json:"category"`
}

type Syncer interface {
	Sync() error
}

func NewSyncer(lc fx.Lifecycle, ai ai.AI, categoryDB DB, domainsDB domains.DB) Syncer {
	timeToSleepStr := os.Getenv("SYNCER_SECONDS_TO_SLEEP")
	secondsToSleep, err := strconv.Atoi(timeToSleepStr)
	if err != nil {
		secondsToSleep = 60
	}

	s := &syncer{ai, categoryDB, domainsDB, secondsToSleep}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				for {
					fmt.Println("Syncing categories")
					err := s.Sync()
					if err != nil {
						log.Printf("Error on syncer: %s\n", err)
					}

					time.Sleep(1 * time.Minute)
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
			Considering domain %s, which content category from these listed below you think is a fit?
			Answer me only with a json with a 'category' key and an array with sanitized categories all in lower case, with spaces replaced with underscore and ordered by relevance without markdown. Consides maximum 3 categories.

			- adult
			- advertisements
			- alcohol
			- animals_and_pets
			- arts
			- astrology
			- auctions
			- business_and_industry
			- cannabis
			- chat_and_instant_messaging
			- cheating_and_plagiarism
			- child_abuse_content
			- cloud_and_data_centers
			- computer_security
			- computers_and_internet
			- conventions_conferences_and_trade_shows
			- cryptocurrency
			- dating
			- digital_postcards
			- dining_and_drinking
			- diy_projects
			- encrypted_dns
			- dynamic_and_residential
			- education
			- entertainment
			- extreme
			- fashion
			- file_transfer_services
			- filter_avoidance
			- finance
			- freeware_and_shareware
			- gambling
			- games
			- government_and_law
			- hacking
			- hate_speech
			- health_and_medicine
			- humor
			- hunting
			- illegal_activities
			- illegal_downloads
			- illegal_drugs
			- infrastructure_and_content_delivery_networks
			- internet_of_things
			- internet_telephony
			- job_search
			- lingerie_and_swimsuits
			- lotteries
			- military
			- mobile_phones
			- museums
			- nature_and_conservation
			- news
			- non_governmental_organizations
			- non_sexual_nudity
			- not_actionable
			- online_communities
			- online_document_sharing_and_collaboration
			- online_meetings
			- online_storage_and_backup
			- online_trading
			- organizational_email
			- paranormal
			- parked_domains
			- peer_file_transfer
			- personal_sites
			- personal_vpn
			- photo_search_and_images
			- politics
			- pornography
			- private_ip_addresses_as_host
			- professional_networking
			- real_estate
			- recipes_and_food
			- reference
			- regional_restricted_law_germany
			- regional_restricted_law_great_britain
			- regional_restricted_law_italy
			- regional_restricted_law_poland
			- religion
			- saas_and_b2b
			- safe_for_kids
			- science_and_technology
			- search_engines_and_portals
			- sex_education
			- shopping
			- social_networking
			- social_science
			- society_and_culture
			- software_updates
			- sports_and_recreation
			- streaming_audio
			- streaming_video
			- terrorism
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
		time.Sleep(time.Duration(s.secondsToSleep) * time.Second)
	}

	return nil
}
