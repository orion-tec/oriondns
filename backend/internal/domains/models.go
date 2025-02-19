package domains

import "time"

type Domain struct {
	Name      string
	UsedCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}
