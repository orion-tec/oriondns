package domains

import "time"

type Domain struct {
	Domain    string
	UsedCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}
