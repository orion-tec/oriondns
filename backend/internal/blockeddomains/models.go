package blockeddomains

import "time"

type BlockedDomain struct {
	ID        int64
	Domain    string
	Recursive bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
