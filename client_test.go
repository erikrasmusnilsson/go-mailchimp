package mailchimp

import (
	"errors"
	"fmt"
	"testing"
)

func TestAthorization(t *testing.T) {
	key := "123456"
	expectedAuth := "Basic YW55c3RyaW5nOjEyMzQ1Ng=="
	actualAuth := authorization(key)
	if expectedAuth != actualAuth {
		t.Errorf(
			"expected auth to be '%s', but was '%s'",
			expectedAuth,
			actualAuth,
		)
	}
}

func TestClient_PingSuccess(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return []byte("{\"health_status\":\"Everything's Chimpy!\"}"), nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	err := client.Ping()
	if err != nil {
		t.Errorf(
			"expected no error to be returned from Ping invocation but '%s' was returned",
			err.Error(),
		)
	}
	if mock.GetCalls != 1 {
		t.Errorf(
			"expected 1 get call to provider, but there was %d",
			mock.GetCalls,
		)
	}
}

func TestClient_PingFailureToReadResponse(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return []byte("{\"status\":\"Everything's Chimpy!\"}"), nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	err := client.Ping()
	if err == nil {
		t.Error(
			"expected error to be returned from Ping invocation but none was returned",
		)
	}
}

func TestClient_PingMailChimpFailure(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return nil, errors.New("something went wrong")
		},
	}
	client := NewCustomDependencyClient(&mock)
	err := client.Ping()
	if err == nil {
		t.Error(
			"expected error to be returned from Ping invocation but none was returned",
		)
	}
}

func TestClient_PingCallsProviderWithCorrectParams(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			if s != "/ping" {
				t.Errorf(
					"expected uri to be '/ping', but was '%s'",
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.Ping()
}

func TestClient_Batch500Limit(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	members := make([]Member, 0, 501)
	for i := 0; i < 501; i++ {
		member := Member{EmailAddress: "test@test.com"}
		members = append(members, member)
	}
	err := client.Batch("test-id", members)
	if err == nil {
		t.Error(
			"expected error to be returned with more than 500 members in batch, no error was returned",
		)
	}
	if mock.PostCalls > 0 {
		t.Errorf(
			"expected mailchimp provider post() to not have been called, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_BatchWithUpdate500Limit(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	members := make([]Member, 0, 501)
	for i := 0; i < 501; i++ {
		member := Member{EmailAddress: "test@test.com"}
		members = append(members, member)
	}
	err := client.BatchWithUpdate("test-id", members)
	if err == nil {
		t.Error(
			"expected error to be returned with more than 500 members in batch, no error was returned",
		)
	}
	if mock.PostCalls > 0 {
		t.Errorf(
			"expected mailchimp provider post() to not have been called, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_BatchTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	members := make([]Member, 0, 2)
	for i := 0; i < 2; i++ {
		member := Member{EmailAddress: "test@test.com"}
		members = append(members, member)
	}
	err := client.Batch("test-id", members)
	if err != nil {
		t.Error("expected no error to be returned, but one was")
	}
	if mock.PostCalls != 1 {
		t.Errorf(
			"expected provider Post() to have been called once, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_CreateListTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	list := List{
		Name:               "Test",
		PermissionReminder: "Permission reminder",
	}
	client.CreateList(list)
	if mock.PostCalls != 1 {
		t.Errorf(
			"expected provider Post() to have been called once, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_FetchListsTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.FetchLists()
	if mock.GetCalls != 1 {
		t.Errorf(
			"expected provider Get() to have been called once, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_FetchListTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.FetchList("test-id")
	if mock.GetCalls != 1 {
		t.Errorf(
			"expected provider Get() to have been called once, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_UpdateListTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		PatchMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	list := List{
		Name:               "Test",
		PermissionReminder: "Permission reminder",
	}
	client.UpdateList("test-id", list)
	if mock.PatchCalls != 1 {
		t.Errorf(
			"expected provider Patch() to have been called once, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_DeleteListTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.DeleteList("test-id")
	if mock.DeleteCalls != 1 {
		t.Errorf(
			"expected provider Delete() to have been called once, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_BatchWithUpdateTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	members := make([]Member, 0, 2)
	for i := 0; i < 2; i++ {
		member := Member{EmailAddress: "test@test.com"}
		members = append(members, member)
	}
	err := client.BatchWithUpdate("test-id", members)
	if err != nil {
		t.Error("expected no error to be returned, but one was")
	}
	if mock.PostCalls != 1 {
		t.Errorf(
			"expected provider Post() to have been called once, was called %d times",
			mock.PostCalls,
		)
	}
}

func TestClient_BatchCallsProviderWithCorrectParams(t *testing.T) {
	testListId := "test-id"
	members := make([]Member, 0, 2)
	for i := 0; i < 2; i++ {
		member := Member{EmailAddress: "test@test.com"}
		members = append(members, member)
	}
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s", testListId) {
				t.Errorf(
					"expected uri to be /lists/%s, but was %s",
					testListId,
					s,
				)
			}
			payload := i.(batch)
			if len(payload.Members) != 2 {
				t.Errorf(
					"expected payload length to be 2, but was %d",
					len(payload.Members),
				)
			}
			if payload.UpdateExisting == true {
				t.Error("expected update_existing to be false, was true")
			}
			for _, member := range payload.Members {
				if member.EmailAddress != "test@test.com" {
					t.Errorf(
						"expected all members to have email adress test@test.com, but found %s",
						member.EmailAddress,
					)
				}
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.Batch(testListId, members)
}

func TestClient_BatchWithUpdateCallsProviderWithCorrectParams(t *testing.T) {
	testListId := "test-id"
	members := make([]Member, 0, 2)
	for i := 0; i < 2; i++ {
		member := Member{EmailAddress: "test@test.com"}
		members = append(members, member)
	}
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s", testListId) {
				t.Errorf(
					"expected uri to be /lists/%s, but was %s",
					testListId,
					s,
				)
			}
			payload := i.(batch)
			if len(payload.Members) != 2 {
				t.Errorf(
					"expected payload length to be 2, but was %d",
					len(payload.Members),
				)
			}
			if payload.UpdateExisting == false {
				t.Error(
					"expected update_existing to be true, but was false",
				)
			}
			for _, member := range payload.Members {
				if member.EmailAddress != "test@test.com" {
					t.Errorf(
						"expected all members to have email adress test@test.com, but found %s",
						member.EmailAddress,
					)
				}
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.BatchWithUpdate(testListId, members)
}

func TestClient_DeleteListCallsProviderWithCorrectParams(t *testing.T) {
	testListId := "test-id"
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s", testListId) {
				t.Errorf(
					"expected uri to be /lists/%s, but was %s",
					testListId,
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.DeleteList(testListId)
}

func TestClient_UpdateListCallsProviderWithCorrectParams(t *testing.T) {
	testListId := "test-id"
	mock := MailChimpProviderMock{
		PatchMock: func(s string, i interface{}) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s", testListId) {
				t.Errorf(
					"expected uri to be /lists/%s, but was %s",
					testListId,
					s,
				)
			}
			payload := i.(List)
			if payload.Name != "Test" {
				t.Errorf(
					"expected list name to be 'Test', but was '%s'",
					payload.Name,
				)
			}
			if payload.PermissionReminder != "This is a test" {
				t.Errorf(
					"expected list permission reminder to be 'This is a test', but was '%s'",
					payload.PermissionReminder,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	list := List{
		Name:               "Test",
		PermissionReminder: "This is a test",
	}
	client.UpdateList(testListId, list)
}

func TestClient_FetchListCallsProviderWithCorrectParams(t *testing.T) {
	testListId := "test-id"
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s", testListId) {
				t.Errorf(
					"expected uri to be /lists/%s, but was %s",
					testListId,
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.FetchList(testListId)
}

func TestClient_FetchListsCallsProviderWithCorrectParams(t *testing.T) {
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			if s != "/lists" {
				t.Errorf(
					"expected uri to be /lists, but was %s",
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.FetchLists()
}

func TestClient_CreateListCallsProviderWithCorrectParams(t *testing.T) {
	list := List{
		Name:               "Test list",
		PermissionReminder: "permission reminder",
	}
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			if s != "/lists" {
				t.Errorf(
					"expected uri to be /lists, but was %s",
					s,
				)
			}
			payload := i.(List)
			if payload.Name != list.Name {
				t.Errorf(
					"expected list name to be '%s', but was '%s'",
					list.Name,
					payload.Name,
				)
			}
			if payload.PermissionReminder != list.PermissionReminder {
				t.Errorf(
					"expected list permission reminder to be '%s', but was '%s'",
					list.PermissionReminder,
					payload.PermissionReminder,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.CreateList(list)
}

func TestClient_FetchMemberTagsCallsProviderWithCorrectParams(t *testing.T) {
	expectedMemberID := hashMd5("test@test.com")
	expectedListID := "list-id"
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s/members/%s/tags", expectedListID, expectedMemberID) {
				t.Errorf(
					"expected uri to be /lists/%s/members/%s/tags, but was %s",
					expectedListID,
					expectedMemberID,
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.FetchMemberTags("list-id", "test@test.com")
}

func TestClient_UpdateMemberTagsCallsProviderWithCorrectParams(t *testing.T) {
	expectedMemberID := hashMd5("test@test.com")
	expectedListID := "list-id"
	tag := Tag{
		Name:   "Test",
		Status: "active",
	}
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s/members/%s/tags", expectedListID, expectedMemberID) {
				t.Errorf(
					"expected uri to be /lists/%s/members/%s/tags, but was %s",
					expectedListID,
					expectedMemberID,
					s,
				)
			}
			payload := i.(updateMemberTagsPayload)
			if payload.IsSyncing {
				t.Error("expected is_syncing to be false, but was true")
			}
			if len(payload.Tags) != 1 {
				t.Errorf(
					"expected number of tags to be 1, but was %d",
					len(payload.Tags),
				)
			}
			if payload.Tags[0].Name != tag.Name {
				t.Errorf(
					"expected tag name to be '%s', but was '%s'",
					tag.Name,
					payload.Tags[0].Name,
				)
			}
			if payload.Tags[0].Status != tag.Status {
				t.Errorf(
					"expected tag status to be '%s', but was '%s'",
					tag.Status,
					payload.Tags[0].Status,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.UpdateMemberTags("list-id", "test@test.com", []Tag{tag})
}

func TestClient_UpdateMemberTagsSyncCallsProviderWithCorrectParams(t *testing.T) {
	expectedMemberID := hashMd5("test@test.com")
	expectedListID := "list-id"
	tag := Tag{
		Name:   "Test",
		Status: "active",
	}
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s/members/%s/tags", expectedListID, expectedMemberID) {
				t.Errorf(
					"expected uri to be /lists/%s/members/%s/tags, but was %s",
					expectedListID,
					expectedMemberID,
					s,
				)
			}
			payload := i.(updateMemberTagsPayload)
			if !payload.IsSyncing {
				t.Error("expected is_syncing to be true, but was false")
			}
			if len(payload.Tags) != 1 {
				t.Errorf(
					"expected number of tags to be 1, but was %d",
					len(payload.Tags),
				)
			}
			if payload.Tags[0].Name != tag.Name {
				t.Errorf(
					"expected tag name to be '%s', but was '%s'",
					tag.Name,
					payload.Tags[0].Name,
				)
			}
			if payload.Tags[0].Status != tag.Status {
				t.Errorf(
					"expected tag status to be '%s', but was '%s'",
					tag.Status,
					payload.Tags[0].Status,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.UpdateMemberTagsSync("list-id", "test@test.com", []Tag{tag})
}

func TestClient_ArchiveMemberCallsDelete(t *testing.T) {
	expectedListID := "list-id"
	expectedMemberEmail := "test@test.com"
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.ArchiveMember(expectedListID, expectedMemberEmail)
	if mock.DeleteCalls != 1 {
		t.Errorf(
			"expected Delete to have been called once but found %d",
			mock.DeleteCalls,
		)
	}
}

func TestClient_ArchiveMemberCallsDeleteWithCorrectURI(t *testing.T) {
	expectedListID := "list-id"
	expectedMemberEmail := "test@test.com"
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s/members/%s", expectedListID, hashMd5(expectedMemberEmail)) {
				t.Errorf(
					"expected uri to be /lists/%s/members/%s but was %s",
					expectedListID,
					hashMd5(expectedMemberEmail),
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.ArchiveMember(expectedListID, expectedMemberEmail)
}

func TestClient_CreateWebhookCallsProviderWithCorrectParams(t *testing.T) {
	expectedListID := "list-id"
	webhook := Webhook{
		URL:    "http://test.com",
		ListID: expectedListID,
	}
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s/webhooks", expectedListID) {
				t.Errorf(
					"expected uri to be /lists/%s/webhooks, but was %s",
					expectedListID,
					s,
				)
			}
			payload := i.(CreateWebhookRequestPayload)
			if payload.URL != webhook.URL {
				t.Errorf(
					"expected Webhook URL to be %s but was %s",
					webhook.URL,
					payload.URL,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.CreateWebhook(webhook)
}

func TestClient_CreateWebhookReturnsErrorIfProviderFails(t *testing.T) {
	expectedListID := "list-id"
	webhook := Webhook{
		URL:    "http://test.com",
		ListID: expectedListID,
	}
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, errors.New("mocked error")
		},
	}
	client := NewCustomDependencyClient(&mock)
	_, err := client.CreateWebhook(webhook)
	if err == nil {
		t.Error("expected error to be returned but none was")
	}
}

func TestClient_FetchWebhooksCallsProviderWithCorrectParams(t *testing.T) {
	expectedListID := "list-id"
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			if s != fmt.Sprintf("/lists/%s/webhooks", expectedListID) {
				t.Errorf(
					"expected uri to be /lists/%s/webhooks, but was %s",
					expectedListID,
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.FetchWebhooks(expectedListID)
}

func TestClient_FetchWebhooksReturnsErrorIfProviderFails(t *testing.T) {
	expectedListID := "list-id"
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return nil, errors.New("mocked error")
		},
	}
	client := NewCustomDependencyClient(&mock)
	_, err := client.FetchWebhooks(expectedListID)
	if err == nil {
		t.Error("expected error to be returned but none was")
	}
}

func TestClient_FetchWebhookCallsProviderWithCorrectParams(t *testing.T) {
	expectedListID := "list-id"
	expectedWebhookID := "webhook-id"
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			if s != fmt.Sprintf(
				"/lists/%s/webhooks/%s",
				expectedListID,
				expectedWebhookID,
			) {
				t.Errorf(
					"expected uri to be /lists/%s/webhooks/%s, but was %s",
					expectedListID,
					expectedWebhookID,
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.FetchWebhook(expectedListID, expectedWebhookID)
}

func TestClient_FetchWebhookReturnsErrorIfProviderFails(t *testing.T) {
	expectedListID := "list-id"
	expectedWebhookID := "webhook-id"
	mock := MailChimpProviderMock{
		GetMock: func(s string) ([]byte, error) {
			return nil, errors.New("mocked error")
		},
	}
	client := NewCustomDependencyClient(&mock)
	_, err := client.FetchWebhook(expectedListID, expectedWebhookID)
	if err == nil {
		t.Error("expected error to be returned but none was")
	}
}

func TestClient_DeleteWebhookCallsProviderWithCorrectParams(t *testing.T) {
	expectedListID := "list-id"
	expectedWebhookID := "webhook-id"
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			if s != fmt.Sprintf(
				"/lists/%s/webhooks/%s",
				expectedListID,
				expectedWebhookID,
			) {
				t.Errorf(
					"expected uri to be /lists/%s/webhooks/%s, but was %s",
					expectedListID,
					expectedWebhookID,
					s,
				)
			}
			return nil, nil
		},
	}
	client := NewCustomDependencyClient(&mock)
	client.DeleteWebhook(expectedListID, expectedWebhookID)
}

func TestClient_DeleteWebhookReturnsErrorIfProviderFails(t *testing.T) {
	expectedListID := "list-id"
	expectedWebhookID := "webhook-id"
	mock := MailChimpProviderMock{
		DeleteMock: func(s string) ([]byte, error) {
			return nil, errors.New("mocked error")
		},
	}
	client := NewCustomDependencyClient(&mock)
	err := client.DeleteWebhook(expectedListID, expectedWebhookID)
	if err == nil {
		t.Error("expected error to be returned but none was")
	}
}
