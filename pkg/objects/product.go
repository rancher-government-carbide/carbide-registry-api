package objects

import (
	"time"
)

type Product struct {
	Id        int32
	Name      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
