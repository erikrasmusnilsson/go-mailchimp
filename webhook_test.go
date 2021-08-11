package mailchimp

import "testing"

func TestWebhookBuilder_AddURL(t *testing.T) {
	testUrl := "https://test.com"
	builder := WebhookBuilder{}
	builder = builder.URL(testUrl)
	if builder.obj.URL != testUrl {
		t.Errorf(
			"expected URL to be '%s' but was '%s'",
			testUrl,
			builder.obj.URL,
		)
	}
}

func TestWebhookBuilder_AddListID(t *testing.T) {
	testID := "1234"
	builder := WebhookBuilder{}
	builder = builder.ListID(testID)
	if builder.obj.ListID != testID {
		t.Errorf(
			"expected ListID to be '%s' but was '%s'",
			testID,
			builder.obj.ListID,
		)
	}
}

func TestWebhookBuilder_AddEvents(t *testing.T) {
	testEvents := WebhookEvents{Subscribe: true}
	builder := WebhookBuilder{}
	builder = builder.Events(testEvents)
	if !equalEvents(builder.obj.Events, testEvents) {
		t.Errorf(
			"expected Events to be '%+v' but was '%+v'",
			testEvents,
			builder.obj.Events,
		)
	}
}

func TestWebhookBuilder_AddSources(t *testing.T) {
	testSources := WebhookSources{Admin: true}
	builder := WebhookBuilder{}
	builder = builder.Sources(testSources)
	if !equalSources(builder.obj.Sources, testSources) {
		t.Errorf(
			"expected Sources to be '%+v' but was '%+v'",
			testSources,
			builder.obj.Sources,
		)
	}
}

func TestWebhookBuilder_Build(t *testing.T) {
	builder := WebhookBuilder{}
	_, err := builder.Build()
	if err == nil {
		t.Error("expected Build to return error without URL but none was")
	}
	builder = builder.URL("test")
	_, err = builder.Build()
	if err == nil {
		t.Error("expected Build to return error without ListID but none was")
	}
	builder = builder.ListID("1234")
	_, err = builder.Build()
	if err != nil {
		t.Error("unexpected error returned from Build with Email and ListID defined")
	}
}
