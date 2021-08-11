package mailchimp

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventProfile     = "profile"
	EventCleaned     = "cleaned"
	EventUpEmail     = "upemail"
	EventCampaign    = "campaign"

	SourceUser  = "user"
	SourceAdmin = "admin"
	SourceAPI   = "api"
)

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

// WebhookClient is an interface allowing the caller to manage webhooks
// for a given list.
type WebhookClient interface {
	AddWebhook(params AddWebhookParams) (Webhook, error)
	GetWebhook(listID string, webookID string) (Webhook, error)
	DeleteWebhook(listID string, webhookID string) error
}

// NewWebhookClient returns a concrete implementation of the WebhookClient
// interface.
func NewWebhookClient(key, region string) WebhookClient {
	return webhookClient{
		provider: mailChimpProvider{
			Region:        region,
			Authorization: authorization(key),
		},
	}
}

func NewCustomDependencyWebhookClient(provider MailChimpProvider) WebhookClient {
	return webhookClient{provider: provider}
}

// webhookClient is a concrete implementation of the WebhookClient interface.
type webhookClient struct {
	provider MailChimpProvider
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

type AddWebhookParams struct {
	URL     string
	ListID  string
	Events  WebhookEvents
	Sources WebhookSources
}

type AddWebhookRequestPayload struct {
	URL     string         `json:"url"`
	Events  WebhookEvents  `json:"events"`
	Sources WebhookSources `json:"sources"`
}

func (c webhookClient) AddWebhook(params AddWebhookParams) (Webhook, error) {
	body, err := c.provider.Post(
		fmt.Sprintf("/lists/%s/webhooks", params.ListID),
		AddWebhookRequestPayload{
			URL:     params.URL,
			Events:  params.Events,
			Sources: params.Sources,
		},
	)
	if err != nil {
		return NullWebhook, err
	}
	webhook := Webhook{}
	if err := json.Unmarshal(body, &webhook); err != nil {
		return NullWebhook, err
	}
	return webhook, nil
}

func (c webhookClient) GetWebhook(listID, webhookID string) (Webhook, error) {
	body, err := c.provider.Get(
		fmt.Sprintf(
			"/lists/%s/webhooks/%s",
			listID,
			webhookID,
		),
	)
	if err != nil {
		return NullWebhook, err
	}
	webhook := Webhook{}
	if err := json.Unmarshal(body, &webhook); err != nil {
		return NullWebhook, err
	}
	return webhook, nil
}

func (c webhookClient) DeleteWebhook(listID, webhookID string) error {
	_, err := c.provider.Delete(
		fmt.Sprintf(
			"/lists/%s/webhooks/%s",
			listID,
			webhookID,
		),
	)
	return err
}
