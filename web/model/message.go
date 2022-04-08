package model

import (
	"time"
	"strconv"
)

type (
	Message struct {
		ID               string    `json:"id,string" db:"id"`
		UUID             string    `json:"-" db:"uuid"`
		TwitterName      string    `json:"twitter" db:"twitter_name"`
		Name             string    `json:"name" db:"name"`
		Location         string    `json:"location" db:"location"`
		Message          string    `json:"message" db:"message"`
		Image            string    `json:"image,omitempty" db:"image"`
		Avatar           int       `json:"pfp" db:"avatar"`
		Pending          bool      `json:"-" db:"is_pending"`
		CreatedAt        time.Time `json:"created_at,omitempty" db:"created_at"`
		DiscordMessageID string    `json:"-" db:"discord_message_id"`
	}
)

func (m *model) AddMessage(msg *Message) (string, error) {
	res, err := m.db.NamedExec(`INSERT INTO
						messages (uuid, twitter_name, name, location, message, avatar, image)
						VALUES (:uuid, :twitter_name, :name, :location, :message, :avatar, :image)`, msg)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	
	return strconv.FormatInt(id, 10), nil
}

func (m *model) GetMessage(id string) (*Message, error) {
	var msg Message
	if err := m.db.Get(&msg, `SELECT * FROM messages WHERE id = ?`, id); err != nil {
		return nil, err
	}
	
	return &msg, nil
}

func (m *model) AcceptMessage(id string) (err error) {
	_, err = m.db.Exec(`UPDATE messages SET is_pending = 0 WHERE id = ?`, id)
	if err != nil {
		return
	}

	return
}

func (m *model) GetAcceptedMessages() ([]*Message, error) {
	messages := []*Message{}
	if err := m.db.Select(&messages, `SELECT id, twitter_name, name, location, message, image, avatar, created_at FROM messages WHERE is_pending = 0`); err != nil {
		return nil, err
	}

	return messages, nil
}


func (m *model) SetDiscordMessageID(id string, discordMsgId string) (err error) {
	_, err = m.db.Exec(`UPDATE messages SET discord_message_id = ? WHERE id = ?`, discordMsgId, id)
	if err != nil {
		return
	}

	return
}