package mailchimp

import "testing"

func TestMemberBuilder_AddEmailAddress(t *testing.T) {
	testEmailAddress := "test@testsson.com"
	builder := MemberBuilder{}
	builder = builder.EmailAddress(testEmailAddress)
	if builder.obj.EmailAddress != testEmailAddress {
		t.Errorf(
			"expected email address to be '%s' but was '%s'",
			testEmailAddress,
			builder.obj.EmailAddress,
		)
	}
}

func TestMemberBuilder_AddEmailType(t *testing.T) {
	testEmailType := "html"
	builder := MemberBuilder{}
	builder = builder.EmailType(testEmailType)
	if builder.obj.EmailType != testEmailType {
		t.Errorf(
			"expected email type to be '%s' but was '%s'",
			testEmailType,
			builder.obj.EmailAddress,
		)
	}
}

func TestMemberBuilder_AddStatus(t *testing.T) {
	builder := MemberBuilder{}
	builder = builder.StatusSubscribed()
	if builder.obj.Status != "subscribed" {
		t.Errorf(
			"expected status to be 'subscribed' but was '%s'",
			builder.obj.Status,
		)
	}
	builder = builder.StatusPending()
	if builder.obj.Status != "pending" {
		t.Errorf(
			"expected status to be 'pending' but was '%s'",
			builder.obj.Status,
		)
	}
	builder = builder.StatusUnsubscribed()
	if builder.obj.Status != "unsubscribed" {
		t.Errorf(
			"expected status to be 'unsubscribed' but was '%s'",
			builder.obj.Status,
		)
	}
	builder = builder.StatusCleaned()
	if builder.obj.Status != "cleaned" {
		t.Errorf(
			"expected status to be 'cleaned' but was '%s'",
			builder.obj.Status,
		)
	}
}

func TestMemberBuilder_AddMergeField(t *testing.T) {
	key := "TEST_KEY"
	value := "TEST_VALUE"
	builder := MemberBuilder{}
	builder = builder.MergeField(key, value)
	if builder.obj.MergeFields[key] != value {
		t.Errorf(
			"expected merge field %s to be '%s' but was '%s'",
			key,
			value,
			builder.obj.MergeFields[key],
		)
	}
}

func TestMemberBuilder_Build(t *testing.T) {
	_, err := MemberBuilder{}.EmailAddress("").Build()
	if err == nil {
		t.Error("expected error for build() on invalid member")
	}

	_, err = MemberBuilder{}.EmailAddress("foo@bar.com").Build()
	if err != nil {
		t.Error("unexpected error on build() for valid member")
	}
}

func TestMember_StatusSubscribed(t *testing.T) {
	member := Member{Status: "subscribed"}
	if !member.Subscribed() {
		t.Error("expected member.Subscribed to return true but returned false")
	}
	if member.Unsubscribed() {
		t.Error("expected member.Unsubscribed to return false but returned true")
	}
	if member.Cleaned() {
		t.Error("expected member.Cleaned to return false but returned true")
	}
	if member.Pending() {
		t.Error("expected member.Pending to return false but returned true")
	}
}

func TestMember_StatusUnsubscribed(t *testing.T) {
	member := Member{Status: "unsubscribed"}
	if member.Subscribed() {
		t.Error("expected member.Subscribed to return false but returned true")
	}
	if !member.Unsubscribed() {
		t.Error("expected member.Unsubscribed to return true but returned false")
	}
	if member.Cleaned() {
		t.Error("expected member.Cleaned to return false but returned true")
	}
	if member.Pending() {
		t.Error("expected member.Pending to return false but returned true")
	}
}

func TestMember_StatusCleaned(t *testing.T) {
	member := Member{Status: "cleaned"}
	if member.Subscribed() {
		t.Error("expected member.Subscribed to return false but returned true")
	}
	if member.Unsubscribed() {
		t.Error("expected member.Unsubscribed to return false but returned true")
	}
	if !member.Cleaned() {
		t.Error("expected member.Cleaned to return true but returned false")
	}
	if member.Pending() {
		t.Error("expected member.Pending to return false but returned true")
	}
}

func TestMember_StatusPending(t *testing.T) {
	member := Member{Status: "pending"}
	if member.Subscribed() {
		t.Error("expected member.Subscribed to return false but returned true")
	}
	if member.Unsubscribed() {
		t.Error("expected member.Unsubscribed to return false but returned true")
	}
	if member.Cleaned() {
		t.Error("expected member.Cleaned to return false but returned true")
	}
	if !member.Pending() {
		t.Error("expected member.Pending to return true but returned false")
	}
}
