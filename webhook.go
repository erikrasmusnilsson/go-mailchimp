package mailchimp

import "fmt"

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
