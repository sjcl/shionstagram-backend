package model

import (
	"time"
)

type (
	Message struct {
		ID          uint64    `json:"id,string" db:"id"`
		UUID        string    `json:"-" db:"uuid"`
		TwitterName string    `json:"twitter" db:"twitter_name"`
		Name        string    `json:"name" db:"name"`
		Location    string    `json:"location" db:"location"`
		Message     string    `json:"message" db:"message"`
		ImageSrc    string    `json:"img_src,omitempty" db:"img_src"`
		Avatar      uint      `json:"pfp" db:"avatar"`
		Pending     bool      `json:"pending,omitempty" db:"is_pending"`
		CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	}
)

func (m *model) AddMessage(msg *Message) (err error) {
	_, err = m.db.Exec(`INSERT INTO
						messages (twitter_name, name, location, message, pfp, image)
						VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return
	}
	
	return
}