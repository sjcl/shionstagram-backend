package model

import (
	"github.com/jmoiron/sqlx"
)

type (
	model struct {
		db *sqlx.DB
	}

	Database interface {
		AddMessage(msg *Message) (int64, error)
		GetMessage(id uint64) (*Message, error)
		AcceptMessage(id uint64) (err error)
		GetAcceptedMessages() ([]*Message, error)
		SetDiscordMessageID(id int64, discordMsgId string) (err error)
	}
)

func NewModel(db *sqlx.DB) *model {
	return &model{
		db: db,
	}
}