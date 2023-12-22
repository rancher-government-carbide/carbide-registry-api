package objects

import (
	"time"
)

type ReleaseImageMapping struct {
	Id        int32
	ReleaseId *int32
	ImageId   *int32
	CreatedAt time.Time
	UpdatedAt time.Time
}
