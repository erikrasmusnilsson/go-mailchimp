package mailchimp

import (
	"encoding/json"
	"errors"
	"testing"
)

func equalEvents(a, b WebhookEvents) bool {
	return a.Campaign == b.Campaign &&
		a.Cleaned == b.Cleaned &&
		a.Profile == b.Profile &&
		a.Subscribe == b.Subscribe &&
		a.Unsubscribe == b.Unsubscribe &&
		a.UpEmail == b.UpEmail
}

func equalSources(a, b WebhookSources) bool {
	return a.API == b.API && a.Admin == b.Admin && a.User == b.User
}

func equalWebhooks(a, b Webhook) bool {
	return a.ListID == b.ListID &&
		a.URL == b.URL &&
		equalEvents(a.Events, b.Events) &&
		equalSources(a.Sources, b.Sources)
}

func TestWebhookClient_AddSuccess(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			raw, _ := json.Marshal(NullWebhook)
			return raw, nil
		},
	}
	client := NewCustomDependencyWebhookClient(&mock)
	wh, err := client.AddWebhook(AddWebhookParams{
		URL:     "https://test.com",
		ListID:  "1234",
		Events:  WebhookEvents{},
		Sources: WebhookSources{},
	})
	if err != nil {
		t.Errorf("unexpected error returned from AddWebhook: %s", err.Error())
	}
	if !equalWebhooks(wh, NullWebhook) {
		t.Errorf(
			"expected returned webhook to equal NullWebhook but was %+v",
			wh,
		)
	}
	if mock.PostCalls != 1 {
		t.Errorf(
			"expected 1 POST call to have been made but was %d",
			mock.PostCalls,
		)
	}
}

func TestWebhookClient_AddFailure(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, errors.New("mocked error")
		},
	}
	client := NewCustomDependencyWebhookClient(&mock)
	_, err := client.AddWebhook(AddWebhookParams{
		URL:     "https://test.com",
		ListID:  "1234",
		Events:  WebhookEvents{},
		Sources: WebhookSources{},
	})
	if err == nil {
		t.Error("expected error returned from AddWebhook but none was")
	}
	if mock.PostCalls != 1 {
		t.Errorf(
			"expected 1 POST call to have been made but was %d",
			mock.PostCalls,
		)
	}
}

func TestWebhookClient_GetSuccess(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			raw, _ := json.Marshal(NullWebhook)
			return raw, nil
		},
	}
	client := NewCustomDependencyWebhookClient(&mock)
	wh, err := client.GetWebhook("1234", "1234")
	if err != nil {
		t.Errorf("unexpected error returned from GetWebhook: %s", err.Error())
	}
	if !equalWebhooks(wh, NullWebhook) {
		t.Errorf(
			"expected returned webhook to equal NullWebhook but was %+v",
			wh,
		)
	}
	if mock.GetCalls != 1 {
		t.Errorf(
			"expected 1 GET call to have been made but was %d",
			mock.GetCalls,
		)
	}
}

func TestWebhookClient_GetFailure(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return nil, errors.New("mocked error")
		},
	}
	client := NewCustomDependencyWebhookClient(&mock)
	_, err := client.GetWebhook("1234", "1234")
	if err == nil {
		t.Error("expected error returned from GetWebhook but none was")
	}
	if mock.GetCalls != 1 {
		t.Errorf(
			"expected 1 GET call to have been made but was %d",
			mock.GetCalls,
		)
	}
}

func TestWebhookClient_DeleteSuccess(t *testing.T) {
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			raw, _ := json.Marshal(NullWebhook)
			return raw, nil
		},
	}
	client := NewCustomDependencyWebhookClient(&mock)
	err := client.DeleteWebhook("1234", "1234")
	if err != nil {
		t.Errorf("unexpected error returned from DeleteWebhook: %s", err.Error())
	}
	if mock.DeleteCalls != 1 {
		t.Errorf(
			"expected 1 DELETE call to have been made but was %d",
			mock.DeleteCalls,
		)
	}
}

func TestWebhookClient_DeleteFailure(t *testing.T) {
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			return nil, errors.New("mocked error")
		},
	}
	client := NewCustomDependencyWebhookClient(&mock)
	err := client.DeleteWebhook("1234", "1234")
	if err == nil {
		t.Error("expected error returned from DeleteWebhook but none was")
	}
	if mock.DeleteCalls != 1 {
		t.Errorf(
			"expected 1 DELETE call to have been made but was %d",
			mock.DeleteCalls,
		)
	}
}
