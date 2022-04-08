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

func BuildWebhookRequest(id string, msg *model.Message) *DiscordWebhook {
	var fields []Field

	fields = append(fields, Field{
		Name: "Name",
		Value: msg.Name,
		Inline: true,
	})
	fields = append(fields, Field{
		Name: "Twitter Name",
		Value: msg.TwitterName,
		Inline: true,
	})

	if msg.Location != "" {
		fields = append(fields, Field{
			Name: "Location",
			Value: msg.Location,
			Inline: true,
		})
	}

	fields = append(fields, Field{
		Name: "Avatar",
		Value: strconv.Itoa(msg.Avatar),
		Inline: true,
	})
	fields = append(fields, Field{
		Name: "Message",
		Value: msg.Message,
		Inline: false,
	})

	if msg.Pending {
		fields = append(fields, Field{
			Name: "Status",
			Value: "Pending",
			Inline: true,
		})
		fields = append(fields, Field{
			Name: "Action",
			Value: fmt.Sprintf("[Accept](%s/accept/%s?id=%s)", os.Getenv("API_BASE_URL"), id, msg.UUID),
			Inline: true,
		})
	} else {
		fields = append(fields, Field{
			Name: "Status",
			Value: "Approved",
			Inline: true,
		})	
		fields = append(fields, Field{
			Name: "Action",
			Value: fmt.Sprintf("[Remove](%s/remove/%s?id=%s)", os.Getenv("API_BASE_URL"), id, msg.UUID),
			Inline: true,
		})
	}

	return &DiscordWebhook{
		Username: "Shionstagram",
		AvatarUrl: os.Getenv("WEBHOOK_AVATAR_URL"),
		Embeds: []Embed {
			{
				Title: "New message posted!",
				Color: "10813695",
				Image: Image{
					URL: os.Getenv("API_BASE_URL") + "/images/" + msg.Image,
				},
				Fields: fields,
			},
		},
	}
}

func (h *handler) PostMessage(c echo.Context) error {
	msg := new(model.Message)
	if err := c.Bind(msg); err != nil {
		return err
	}

	msg.TwitterName = strings.TrimLeft(msg.TwitterName, "@")

	uuidObj, _ := uuid.NewRandom()
	msg.UUID = uuidObj.String()

	msg.Pending = true

	id, err := h.Model.AddMessage(msg)
	 if err != nil {
		return err
	}
	
	wh := BuildWebhookRequest(id, msg)

	whPayload, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	whRes, err := http.Post(os.Getenv("WEBHOOK_URL") + "?wait=true", "application/json", bytes.NewBuffer(whPayload))
	if err != nil {
		return err
	}
	defer whRes.Body.Close()

	res := &ResID{
		ID: id,
	}

	if whRes.StatusCode != http.StatusOK {
		fmt.Println("Failed to post webhook message: %d %s", id, msg.UUID)
		return c.JSON(http.StatusCreated, res)
	}

	var whResContent ResID
	err = json.NewDecoder(whRes.Body).Decode(&whResContent)
	if err != nil {
		fmt.Println("Failed to read webhook message response: %d %s", id, msg.UUID)
		return c.JSON(http.StatusCreated, res)
	}

	err = h.Model.SetDiscordMessageID(id, whResContent.ID)
	if err != nil {
		fmt.Println("Failed to store webhook message id: %d %s", id, msg.UUID)
		return c.JSON(http.StatusCreated, res)
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
	id := c.Param("id")

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

	if !msg.Pending {
		return c.JSON(http.StatusBadRequest, &ResMessage{
			Message: "This message is already accepted.",
		})
	}

	msg.Pending = false

	if err := h.Model.UpdateMessagePendingStatus(id, false); err != nil {
		return err
	}

	wh := BuildWebhookRequest(id, msg)

	whPayload, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	client := &http.Client{}
	whReq, err := http.NewRequest("PATCH", os.Getenv("WEBHOOK_URL") + "/messages/" + msg.DiscordMessageID, bytes.NewBuffer(whPayload))
	if err != nil {
		return err
	}

	whReq.Header.Add("Content-Type", "application/json")

	whRes, err := client.Do(whReq)
	if err != nil {
		return err
	}
	defer whRes.Body.Close()


	res := &ResMessage{
		Message: "Message accepted. You can close this tab now.",
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) RemoveMessage(c echo.Context) error {
	id := c.Param("id")

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

	if msg.Pending {
		return c.JSON(http.StatusBadRequest, &ResMessage{
			Message: "This message is already removed.",
		})
	}

	msg.Pending = true

	if err := h.Model.UpdateMessagePendingStatus(id, true); err != nil {
		return err
	}

	wh := BuildWebhookRequest(id, msg)

	whPayload, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	client := &http.Client{}
	whReq, err := http.NewRequest("PATCH", os.Getenv("WEBHOOK_URL") + "/messages/" + msg.DiscordMessageID, bytes.NewBuffer(whPayload))
	if err != nil {
		return err
	}

	whReq.Header.Add("Content-Type", "application/json")

	whRes, err := client.Do(whReq)
	if err != nil {
		return err
	}
	defer whRes.Body.Close()


	res := &ResMessage{
		Message: "Message removed. You can close this tab now.",
	}

	return c.JSON(http.StatusOK, res)

}

func (h *handler) GetMessages(c echo.Context) error {
	messages, err := h.Model.GetAcceptedMessages()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, messages)
}