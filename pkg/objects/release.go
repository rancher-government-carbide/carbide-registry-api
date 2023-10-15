package objects

import "time"

type Release struct {
	Id          int32
	ProductId   *int32
	Name        *string
	TarballLink *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Images []Image
}
