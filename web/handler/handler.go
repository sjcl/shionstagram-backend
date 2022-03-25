package handler

import (
	"github.com/sjcl/shionstagram-backend/web/model"
)

type handler struct {
	Model model.Database
}

func NewHandler(d model.Database) *handler {
	return &handler{
		Model: d,
	}
}