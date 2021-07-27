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
