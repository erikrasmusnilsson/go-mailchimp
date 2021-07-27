package mailchimp

import (
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

func TestClient_Batch500Limit(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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

func TestClient_BatchWithUpdateTestProviderCall(t *testing.T) {
	mock := MailChimpProviderMock{
		PostMock: func(s string, i interface{}) ([]byte, error) {
			return nil, nil
		},
	}
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
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
	client := NewMockClient(&mock)
	client.CreateList(list)
}
