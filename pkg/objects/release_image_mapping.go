package objects

import (
	"time"
)

type Release_Image_Mapping struct {
	Id        int32
	ReleaseId *int32
	ImageId   *int32
	CreatedAt time.Time
	UpdatedAt time.Time
}
