package model

import (
	"github.com/jmoiron/sqlx"
)

type (
	model struct {
		db *sqlx.DB
	}

	Database interface {
		AddMessage(msg *Message) (string, error)
		GetMessage(id string) (*Message, error)
		AcceptMessage(id string) (err error)
		GetAcceptedMessages() ([]*Message, error)
		SetDiscordMessageID(id string, discordMsgId string) (err error)
	}
)

func NewModel(db *sqlx.DB) *model {
	return &model{
		db: db,
	}
}