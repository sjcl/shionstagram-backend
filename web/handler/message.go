package handler

import (
	"net/http"
	"strings"
	"io"
	"os"
	"path/filepath"

	"github.com/sjcl/shionstagram-backend/web/model"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
)

type (
	ResPostImage struct {
		ID string `json:"id"`
	}
)

func (h *handler) PostMessage(c echo.Context) error {
	msg := new(model.Message)
	if err := c.Bind(msg); err != nil {
		return err
	}

	msg.TwitterName = strings.TrimLeft(msg.TwitterName, "@")

	uuidObj, _ := uuid.NewRandom()
	msg.UUID = uuidObj.String()

	if err := h.Model.AddMessage(msg); err != nil {
		return err
	}
	
	return c.NoContent(http.StatusCreated)
}

func (h *handler) PostImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	uuidObj, _ := uuid.NewRandom()
	uuid := uuidObj.String()

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	dst, err := os.Create(filepath.Join("/images", uuid + ext))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	res := &ResPostImage{
		ID: uuid + ext,
	}
	return c.JSON(http.StatusCreated, res)
}