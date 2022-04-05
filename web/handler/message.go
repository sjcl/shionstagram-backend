package handler

import (
	"net/http"

	"github.com/sjcl/shionstagram-backend/web/model"
	"github.com/labstack/echo/v4"
)

func (h *handler) PostMessage(c echo.Context) error {
	msg := new(model.Message)
	if err := c.Bind(msg); err != nil {
		return err
	}
	
	return c.JSON(http.StatusCreated, msg)
}