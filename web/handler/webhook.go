package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/sjcl/shionstagram-backend/web/model"
)

type (
	DiscordWebhook struct {
		Embeds    []Embed `json:"embeds"`
		Username  string  `json:"username"`
		AvatarUrl string  `json:"avatar_url"`
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
		Title     string  `json:"title"`
		Color     string  `json:"color"`
		Fields    []Field `json:"fields"`
		Image     Image   `json:"image"`
		Timestamp string  `json:"timestamp"`
		Footer    Footer  `json:"footer"`
	}

	Footer struct {
		Text    string `json:"text"`
		IconUrl string `json:"icon_url"`
	}
)

func BuildWebhookRequest(id string, msg *model.Message) *DiscordWebhook {
	var fields []Field

	fields = append(fields, Field{
		Name:   "Name",
		Value:  msg.Name,
		Inline: true,
	})
	fields = append(fields, Field{
		Name:   "Twitter Name",
		Value:  msg.TwitterName,
		Inline: true,
	})

	if msg.Location != "" {
		fields = append(fields, Field{
			Name:   "Location",
			Value:  msg.Location,
			Inline: true,
		})
	}

	fields = append(fields, Field{
		Name:   "Message",
		Value:  msg.Message,
		Inline: false,
	})

	if msg.Pending {
		fields = append(fields, Field{
			Name:   "Status",
			Value:  "Pending",
			Inline: true,
		})
		fields = append(fields, Field{
			Name:   "Action",
			Value:  fmt.Sprintf("[Accept](%s/accept/%s?id=%s)", os.Getenv("API_BASE_URL"), id, msg.UUID),
			Inline: true,
		})
	} else {
		fields = append(fields, Field{
			Name:   "Status",
			Value:  "Approved",
			Inline: true,
		})
		fields = append(fields, Field{
			Name:   "Action",
			Value:  fmt.Sprintf("[Revert to pending](%s/remove/%s?id=%s)", os.Getenv("API_BASE_URL"), id, msg.UUID),
			Inline: true,
		})
	}

	timestamp := msg.CreatedAt
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	var embeds []Embed
	if msg.Pending {
		embeds = append(embeds, Embed{
			Title: "New message posted!",
			Color: "15844367",
			Image: Image{
				URL: os.Getenv("API_BASE_URL") + "/images/" + msg.Image,
			},
			Fields:    fields,
			Timestamp: timestamp.Format(time.RFC3339),
			Footer: Footer{
				Text:    fmt.Sprintf("Avatar%d", msg.Avatar),
				IconUrl: fmt.Sprintf("%s/images/pfp/%d.png", os.Getenv("API_BASE_URL"), msg.Avatar),
			},
		})
	} else {
		embeds = append(embeds, Embed{
			Title: "Approved message",
			Color: "10813695",
			Image: Image{
				URL: os.Getenv("API_BASE_URL") + "/images/" + msg.Image,
			},
			Fields:    fields,
			Timestamp: timestamp.Format(time.RFC3339),
			Footer: Footer{
				Text:    fmt.Sprintf("Avatar%d", msg.Avatar),
				IconUrl: fmt.Sprintf("%s/images/pfp/%d.png", os.Getenv("API_BASE_URL"), msg.Avatar),
			},
		})
	}

	return &DiscordWebhook{
		Username:  "Shiongram",
		AvatarUrl: os.Getenv("WEBHOOK_AVATAR_URL"),
		Embeds:    embeds,
	}
}
