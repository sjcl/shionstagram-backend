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
		Image       string    `json:"image,omitempty" db:"image"`
		Avatar      int       `json:"pfp" db:"avatar"`
		Pending     bool      `json:"-" db:"is_pending"`
		CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	}
)

func (m *model) AddMessage(msg *Message) (int64, error) {
	res, err := m.db.NamedExec(`INSERT INTO
						messages (uuid, twitter_name, name, location, message, avatar, image)
						VALUES (:uuid, :twitter_name, :name, :location, :message, :avatar, :image)`, msg)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	
	return id, nil
}

func (m *model) GetMessage(id uint64) (*Message, error) {
	var msg Message
	if err := m.db.Get(&msg, `SELECT * FROM messages WHERE id = ?`, id); err != nil {
		return nil, err
	}
	
	return &msg, nil
}

func (m *model) AcceptMessage(id uint64) (err error) {
	_, err = m.db.Exec(`UPDATE messages SET is_pending = 0 WHERE id = ?`, id)
	if err != nil {
		return
	}

	return
}

func (m *model) GetAcceptedMessages() ([]*Message, error) {
	messages := []*Message{}
	if err := m.db.Select(&messages, `SELECT * FROM messages WHERE is_pending = 0`); err != nil {
		return nil, err
	}

	return messages, nil
}