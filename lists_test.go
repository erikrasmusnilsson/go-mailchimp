package mailchimp

import "testing"

func TestListBuilder_Name(t *testing.T) {
	testName := "Foo bar"
	builder := ListBuilder{}
	builder = builder.Name(testName)
	if builder.obj.Name != testName {
		t.Errorf(
			"expected list name to be '%s', but was '%s'",
			testName,
			builder.obj.Name,
		)
	}
}

func TestListBuilder_PermissionReminder(t *testing.T) {
	testPermissionReminder := "Foo bar"
	builder := ListBuilder{}
	builder = builder.PermissionReminder(testPermissionReminder)
	if builder.obj.PermissionReminder != testPermissionReminder {
		t.Errorf(
			"expected list permission reminder to be '%s', but was '%s'",
			testPermissionReminder,
			builder.obj.PermissionReminder,
		)
	}
}

func TestListBuilder_EmailTypeOption(t *testing.T) {
	testEmailTypeOption := true
	builder := ListBuilder{}
	builder = builder.EmailTypeOption(testEmailTypeOption)
	if builder.obj.EmailTypeOption != testEmailTypeOption {
		t.Errorf(
			"expected list permission reminder to be '%v', but was '%v'",
			testEmailTypeOption,
			builder.obj.EmailTypeOption,
		)
	}
}

func TestListBuilder_Contact(t *testing.T) {
	testContact := Contact{
		Address1: "Testvägen 1",
		State:    "Testmanland",
		Zip:      "123 45",
		Company:  "Testbolaget",
		City:     "Testköping",
		Country:  "Sweden",
	}
	builder := ListBuilder{}
	builder = builder.Contact(testContact)
	if builder.obj.Contact.Address1 != testContact.Address1 {
		t.Errorf(
			"expected address1 to be '%s', but was '%s'",
			testContact.Address1,
			builder.obj.Contact.Address1,
		)
	}
	if builder.obj.Contact.State != testContact.State {
		t.Errorf(
			"expected state to be '%s', but was '%s'",
			testContact.State,
			builder.obj.Contact.State,
		)
	}
	if builder.obj.Contact.Zip != testContact.Zip {
		t.Errorf(
			"expected zip to be '%s', but was '%s'",
			testContact.Zip,
			builder.obj.Contact.Zip,
		)
	}
	if builder.obj.Contact.Company != testContact.Company {
		t.Errorf(
			"expected company to be '%s', but was '%s'",
			testContact.Company,
			builder.obj.Contact.Company,
		)
	}
	if builder.obj.Contact.City != testContact.City {
		t.Errorf(
			"expected city to be '%s', but was '%s'",
			testContact.City,
			builder.obj.Contact.City,
		)
	}
	if builder.obj.Contact.Country != testContact.Country {
		t.Errorf(
			"expected country to be '%s', but was '%s'",
			testContact.Country,
			builder.obj.Contact.Country,
		)
	}
}

func TestListBuilder_CampaignDefaults(t *testing.T) {
	testCampaign := CampaignDefaults{
		FromName:  "Test",
		FromEmail: "test@test.com",
		Subject:   "Test",
		Language:  "SWEDISH",
	}
	builder := ListBuilder{}
	builder = builder.CampaignDefaults(testCampaign)
	if builder.obj.CampaignDefaults.FromName != testCampaign.FromName {
		t.Errorf(
			"expected FromName to be '%s', but was '%s'",
			testCampaign.FromName,
			builder.obj.CampaignDefaults.FromName,
		)
	}
	if builder.obj.CampaignDefaults.FromEmail != testCampaign.FromEmail {
		t.Errorf(
			"expected FromEmail to be '%s', but was '%s'",
			testCampaign.FromEmail,
			builder.obj.CampaignDefaults.FromEmail,
		)
	}
	if builder.obj.CampaignDefaults.Subject != testCampaign.Subject {
		t.Errorf(
			"expected Subject to be '%s', but was '%s'",
			testCampaign.Subject,
			builder.obj.CampaignDefaults.Subject,
		)
	}
	if builder.obj.CampaignDefaults.Language != testCampaign.Language {
		t.Errorf(
			"expected Language to be '%s', but was '%s'",
			testCampaign.Language,
			builder.obj.CampaignDefaults.Language,
		)
	}
}
