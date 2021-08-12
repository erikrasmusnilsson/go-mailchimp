package mailchimp

import (
	"fmt"
)

const (
	StatusSubscribed   = "subscribed"
	StatusUnsubscribed = "unsubscribed"
	StatusPending      = "pending"
	StatusCleaned      = "cleaned"
)

var NullMember = Member{}

type Member struct {
	EmailAddress string            `json:"email_address" mc_validator:"required"`
	EmailType    string            `json:"email_type"`
	Status       string            `json:"status"`
	MergeFields  map[string]string `json:"merge_fields"`
}

func (m Member) Subscribed() bool {
	return m.Status == StatusSubscribed
}

func (m Member) Unsubscribed() bool {
	return m.Status == StatusUnsubscribed
}

func (m Member) Pending() bool {
	return m.Status == StatusPending
}

func (m Member) Cleaned() bool {
	return m.Status == StatusCleaned
}

type MemberBuilder struct {
	obj Member
}

func (sb MemberBuilder) Build() (Member, error) {
	if invalidParams, valid := validate(sb.obj); valid {
		return sb.obj, nil
	} else {
		return NullMember, fmt.Errorf(
			"could not build member due to invalid parameters %v",
			invalidParams,
		)
	}
}

func (sb MemberBuilder) EmailAddress(emailAddress string) MemberBuilder {
	sb.obj.EmailAddress = emailAddress
	return sb
}

func (sb MemberBuilder) EmailType(emailType string) MemberBuilder {
	sb.obj.EmailType = emailType
	return sb
}

func (sb MemberBuilder) StatusPending() MemberBuilder {
	sb.obj.Status = StatusPending
	return sb
}

func (sb MemberBuilder) StatusSubscribed() MemberBuilder {
	sb.obj.Status = StatusSubscribed
	return sb
}

func (sb MemberBuilder) StatusUnsubscribed() MemberBuilder {
	sb.obj.Status = StatusUnsubscribed
	return sb
}

func (sb MemberBuilder) StatusCleaned() MemberBuilder {
	sb.obj.Status = StatusCleaned
	return sb
}

func (sb MemberBuilder) MergeField(name string, value string) MemberBuilder {
	if sb.obj.MergeFields == nil {
		sb.obj.MergeFields = make(map[string]string)
	}
	sb.obj.MergeFields[name] = value
	return sb
}
