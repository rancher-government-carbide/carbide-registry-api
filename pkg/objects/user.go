package objects

import "time"

type User struct {
	Id        int64
	Username  *string
	Password  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
