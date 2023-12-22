package objects

import "time"

type Image struct {
	Id            int32
	ImageName     *string
	ImageSigned   *bool
	TrivySigned   *bool
	TrivyValid    *bool
	SbomSigned    *bool
	SbomValid     *bool
	LastScannedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Releases      []Release
}
