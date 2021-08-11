package mailchimp

import (
	"fmt"
	"time"
)

var NullWebhook = Webhook{}

type Webhook struct {
	ID      string         `json:"id"`
	URL     string         `json:"url" mc_validator:"required"`
	Events  WebhookEvents  `json:"events"`
	Sources WebhookSources `json:"sources"`
	ListID  string         `json:"list_id" mc_validator:"required"`
}

type WebhookBuilder struct {
	obj Webhook
}

func (builder WebhookBuilder) Build() (Webhook, error) {
	if invalidParams, valid := validate(builder.obj); valid {
		return builder.obj, nil
	} else {
		return NullWebhook, fmt.Errorf(
			"could not build webhook due to invalid parameters %v",
			invalidParams,
		)
	}
}

func (builder WebhookBuilder) URL(url string) WebhookBuilder {
	builder.obj.URL = url
	return builder
}

func (builder WebhookBuilder) Events(events WebhookEvents) WebhookBuilder {
	builder.obj.Events = events
	return builder
}

func (builder WebhookBuilder) Sources(sources WebhookSources) WebhookBuilder {
	builder.obj.Sources = sources
	return builder
}

func (builder WebhookBuilder) ListID(listID string) WebhookBuilder {
	builder.obj.ListID = listID
	return builder
}

type SubscribeEvent struct {
	Type    string             `json:"type"`
	FiredAt time.Time          `json:"fired_at"`
	Data    SubscribeEventData `json:"data"`
}

type SubscribeEventData struct {
	ID          string            `json:"id"`
	ListID      string            `json:"list_id"`
	Email       string            `json:"email"`
	EmailType   string            `json:"email_type"`
	MergeFields map[string]string `json:"merges"`
}

type UnsubscribeEvent struct {
	Type    string    `json:"type"`
	FiredAt time.Time `json:"fired_at"`
}

type UnsubscribeEventData struct {
	Action      string            `json:"action"`
	Reason      string            `json:"reason"`
	ID          string            `json:"id"`
	ListID      string            `json:"list_id"`
	Email       string            `json:"email"`
	EmailType   string            `json:"email_type"`
	MergeFields map[string]string `json:"merges"`
}

type WebhookEvents struct {
	Subscribe   bool `json:"subscribe"`
	Unsubscribe bool `json:"unsubscribe"`
	Profile     bool `json:"profile"`
	Cleaned     bool `json:"cleaned"`
	UpEmail     bool `json:"upemail"`
	Campaign    bool `json:"campaign"`
}

type WebhookSources struct {
	User  bool `json:"user"`
	Admin bool `json:"admin"`
	API   bool `json:"api"`
}
