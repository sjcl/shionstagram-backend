package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sjcl/shionstagram-backend/web/model"
)

type (
	ResID struct {
		ID string `json:"id"`
	}

	ResMessage struct {
		Message string `json:"message"`
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

	msg.Pending = true

	if msg.Image == "" {
		rand.Seed(time.Now().UnixNano())
		msg.Image = fmt.Sprintf("randomImages/%d.jpg", rand.Intn(6)+1)
	}

	id, err := h.Model.AddMessage(msg)
	if err != nil {
		return err
	}

	wh := BuildWebhookRequest(id, msg)

	whPayload, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	whRes, err := http.Post(os.Getenv("WEBHOOK_URL")+"?wait=true", "application/json", bytes.NewBuffer(whPayload))
	if err != nil {
		return err
	}
	defer whRes.Body.Close()

	res := &ResID{
		ID: id,
	}

	if whRes.StatusCode != http.StatusOK {
		fmt.Println("Failed to post webhook message: %d %s", id, msg.UUID)
		return c.NoContent(http.StatusInternalServerError)
	}

	var whResContent ResID
	err = json.NewDecoder(whRes.Body).Decode(&whResContent)
	if err != nil {
		fmt.Println("Failed to read webhook message response: %d %s", id, msg.UUID)
		return c.NoContent(http.StatusInternalServerError)
	}

	err = h.Model.SetDiscordMessageID(id, whResContent.ID)
	if err != nil {
		fmt.Println("Failed to store webhook message id: %d %s", id, msg.UUID)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, res)
}

func checkContentType(ct string) bool {
	fmt.Println(ct)
	return ct != "image/jpeg" && ct != "image/png" && ct != "image/gif"
}

func (h *handler) PostImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	ct := file.Header["Content-Type"][0]
	if checkContentType(ct) {
		return c.NoContent(http.StatusBadRequest)
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buf := make([]byte, 512)
	src.Read(buf)
	if checkContentType(http.DetectContentType(buf)) {
		return c.NoContent(http.StatusBadRequest)
	}

	src.Seek(0, 0)

	uuidObj, _ := uuid.NewRandom()
	uuid := uuidObj.String()

	ext := filepath.Ext(file.Filename)
	dst, err := os.Create(filepath.Join("/images", uuid+ext))
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
	whReq, err := http.NewRequest("PATCH", os.Getenv("WEBHOOK_URL")+"/messages/"+msg.DiscordMessageID, bytes.NewBuffer(whPayload))
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
	whReq, err := http.NewRequest("PATCH", os.Getenv("WEBHOOK_URL")+"/messages/"+msg.DiscordMessageID, bytes.NewBuffer(whPayload))
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
