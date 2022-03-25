package model

import (
	"time"
)

type (
	Post struct {
		ID        uint64    `json:"id,string" db:"id"`
		UUID      string    `json:"-" db:"uuid"`
		Name      string    `json:"name" db:"name"`
		Location  string    `json:"location" db:"location"`
		Message   string    `json:"message" db:"message"`
		ImageSrc  string    `json:"img_src,omitempty" db:"img_src"`
		Avatar    uint      `json:"avatar" db:"avatar"`
		Pending   bool      `json:"pending" db:"is_pending"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}
)