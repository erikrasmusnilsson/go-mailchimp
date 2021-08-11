package mailchimp

import (
	"fmt"
)

var (
	NullList      = List{}
	NullListSlice = []List{}
)

type Contact struct {
	Address1 string `json:"address1" mc_validator:"required"`
	Address2 string `json:"address2"`
	State    string `json:"state" mc_validator:"required"`
	Zip      string `json:"zip" mc_validator:"required"`
	Phone    string `json:"phone"`
	Company  string `json:"company" mc_validator:"required"`
	City     string `json:"city" mc_validator:"required"`
	Country  string `json:"country" mc_validator:"required"`
}

type CampaignDefaults struct {
	FromName  string `json:"from_name" mc_validator:"required"`
	FromEmail string `json:"from_email" mc_validator:"required"`
	Subject   string `json:"subject" mc_validator:"required"`
	Language  string `json:"language" mc_validator:"required"`
}

type List struct {
	ID                   string           `json:"id"`
	WebID                int              `json:"web_id"`
	Name                 string           `json:"name" mc_validator:"required"`
	Contact              Contact          `json:"contact"`
	PermissionReminder   string           `json:"permission_reminder" mc_validator:"required"`
	CampaignDefaults     CampaignDefaults `json:"campaign_defaults"`
	EmailTypeOption      bool             `json:"email_type_option"`
	UseArchiveBar        bool             `json:"use_archive_bar"`
	NotifyOnSubscribe    string           `json:"notify_on_subscribe"`
	NotifyOnUnsubscribe  string           `json:"notify_on_unsubscribe"`
	DoubleOptin          bool             `json:"double_optin"`
	MarketingPermissions bool             `json:"marketing_permissions"`
}

type listCollection struct {
	Lists []List `json:"lists"`
}

type ListBuilder struct {
	obj List
}

func (lb ListBuilder) Build() (List, error) {
	if invalidParams, valid := validate(lb.obj); !valid {
		return NullList, fmt.Errorf(
			"could not build list due to invalid parameters %v",
			invalidParams,
		)
	}
	if invalidParams, valid := validate(lb.obj.CampaignDefaults); !valid {
		return NullList, fmt.Errorf(
			"could not build list due to invalid campaign defaults parameters %v",
			invalidParams,
		)
	}
	if invalidParams, valid := validate(lb.obj.Contact); !valid {
		return NullList, fmt.Errorf(
			"could not build list due to invalid contact parameters %v",
			invalidParams,
		)
	}
	return lb.obj, nil
}

func (lb ListBuilder) Name(name string) ListBuilder {
	lb.obj.Name = name
	return lb
}

func (lb ListBuilder) PermissionReminder(permissionReminder string) ListBuilder {
	lb.obj.PermissionReminder = permissionReminder
	return lb
}

func (lb ListBuilder) EmailTypeOption(emailTypeOption bool) ListBuilder {
	lb.obj.EmailTypeOption = emailTypeOption
	return lb
}

func (lb ListBuilder) Contact(c Contact) ListBuilder {
	lb.obj.Contact = c
	return lb
}

func (lb ListBuilder) CampaignDefaults(cd CampaignDefaults) ListBuilder {
	lb.obj.CampaignDefaults = cd
	return lb
}
