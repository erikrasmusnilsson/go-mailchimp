package mailchimp

import (
	"errors"
	"fmt"
)

const (
	StatusSubscribed   = "subscribed"
	StatusUnsubscribed = "unsubscribed"
	StatusPending      = "pending"
	StatusCleaned      = "cleaned"
)

type Member struct {
	EmailAddress string            `json:"email_address" mc_validator:"required"`
	EmailType    string            `json:"email_type"`
	Status       string            `json:"status"`
	MergeFields  map[string]string `json:"merge_fields"`
}

type MemberBuilder struct {
	obj Member
}

func (sb MemberBuilder) Build() (Member, error) {
	if invalidParams, valid := validate(sb.obj); valid {
		return sb.obj, nil
	} else {
		return sb.obj, errors.New(fmt.Sprintf(
			"could not build member due to invalid parameters %v",
			invalidParams,
		))
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
