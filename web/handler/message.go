package handler

import (
	"net/http"
	"strings"
	"io"
	"os"
	"path/filepath"
	"encoding/json"
	"bytes"
	"strconv"
	"fmt"

	"github.com/sjcl/shionstagram-backend/web/model"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
)

type (
	ResID struct {
		ID string `json:"id"`
	}

	ResMessage struct {
		Message string `json:"message"`
	}
)

type (
	DiscordWebhook struct {
		Embeds    []Embed    `json:"embeds"`
		Username  string     `json:"username"`
		AvatarUrl string     `json:"avatar_url"`
	}

	Field struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Inline bool   `json:"inline,omitempty"`
	}

	Image struct {
		URL string `json:"url,omitempty"`
	}
	
	Embed struct {
		Title  string      `json:"title"`
		Color  string      `json:"color"`
		Fields []Field     `json:"fields"`
		Image  Image       `json:"image"`
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

	id, err := h.Model.AddMessage(msg)
	 if err != nil {
		return err
	}

	var whFields []Field 
	if msg.Location != "" {
		whFields = []Field {
			{
				Name: "Name",
				Value: msg.Name,
				Inline: true,
			},
			{
				Name: "Twitter Name",
				Value: msg.TwitterName,
				Inline: true,
			},
			{
				Name: "Location",
				Value: msg.Location,
				Inline: true,
			},
			{
				Name: "Avatar",
				Value: strconv.Itoa(msg.Avatar),
				Inline: true,
			},
			{
				Name: "Message",
				Value: msg.Message,
				Inline: false,
			},
			{
				Name: "Action",
				Value: fmt.Sprintf("[Accept](%s/accept/%d?id=%s)", os.Getenv("API_BASE_URL"), id, msg.UUID),
				Inline: false,
			},
		}
	} else {
		whFields = []Field {
			{
				Name: "Name",
				Value: msg.Name,
				Inline: true,
			},
			{
				Name: "Twitter Name",
				Value: msg.TwitterName,
				Inline: true,
			},
			{
				Name: "Avatar",
				Value: strconv.Itoa(msg.Avatar),
				Inline: true,
			},
			{
				Name: "Message",
				Value: msg.Message,
				Inline: false,
			},
			{
				Name: "Action",
				Value: fmt.Sprintf("[Accept](%s/accept/%d?id=%s)", os.Getenv("API_BASE_URL"), id, msg.UUID),
				Inline: false,
			},
		}
	}
	

	wh := &DiscordWebhook{
		Username: "Shionstagram",
		AvatarUrl: os.Getenv("WEBHOOK_AVATAR_URL"),
		Embeds: []Embed {
			{
				Title: "New message posted!",
				Color: "10813695",
				Image: Image{
					URL: os.Getenv("API_BASE_URL") + "/images/" + msg.Image,
				},
				Fields: whFields,
			},
		},
	}

	whPayload, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	_, err = http.Post(os.Getenv("WEBHOOK_URL"), "application/json", bytes.NewBuffer(whPayload))
	if err != nil {
		return err
	}

	res := &ResID{
		ID: strconv.FormatInt(id, 10),
	}

	return c.JSON(http.StatusCreated, res)
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

	res := &ResID{
		ID: uuid + ext,
	}
	return c.JSON(http.StatusCreated, res)
}

func (h *handler) AcceptMessage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	uuid := c.QueryParam("id")
	if uuid == "" {
		return c.NoContent(http.StatusNotFound)
	}

	msg, err := h.Model.GetMessage(id)
	if err != nil {
		return err
	}

	if msg.UUID != uuid {
		return c.NoContent(http.StatusNotFound)
	}

	if err := h.Model.AcceptMessage(id); err != nil {
		return err
	}

	res := &ResMessage{
		Message: "Message accepted.",
	}

	return c.JSON(http.StatusOK, res)
}