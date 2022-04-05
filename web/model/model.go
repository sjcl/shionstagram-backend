package model

import (
	"github.com/jmoiron/sqlx"
)

type (
	model struct {
		db *sqlx.DB
	}

	Database interface {
		AddMessage(msg *Message) (err error)
	}
)

func NewModel(db *sqlx.DB) *model {
	return &model{
		db: db,
	}
}