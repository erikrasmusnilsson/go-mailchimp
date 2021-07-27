package mailchimp

import "testing"

type testStruct struct {
	RequiredField    string `mc_validator:"required"`
	NonRequiredField string
}

func TestValidator_Validation(t *testing.T) {
	validTest := testStruct{
		RequiredField:    "Foo",
		NonRequiredField: "bar",
	}
	if _, valid := validate(validTest); !valid {
		t.Error("test struct with all required fields was considered invalid")
	}

	validTest = testStruct{
		RequiredField: "Foo",
	}
	if _, valid := validate(validTest); !valid {
		t.Error("test struct with all required fields was considered invalid")
	}

	invalidTest := testStruct{
		NonRequiredField: "bar",
	}
	if _, valid := validate(invalidTest); valid {
		t.Error("test struct without all required fields was considered valid")
	}
}

func TestValidator_InvalidFields(t *testing.T) {
	validTest := testStruct{
		RequiredField:    "Foo",
		NonRequiredField: "bar",
	}
	fields, _ := validate(validTest)
	if len(fields) > 0 {
		t.Errorf("expected number of invalid fields to be 0 but was %d", len(fields))
	}

	validTest = testStruct{
		RequiredField: "Foo",
	}
	fields, _ = validate(validTest)
	if len(fields) > 0 {
		t.Errorf("expected number of invalid fields to be 0 but was %d", len(fields))
	}

	invalidTest := testStruct{
		NonRequiredField: "bar",
	}
	fields, _ = validate(invalidTest)
	if len(fields) != 1 {
		t.Errorf("expected number of invalid fields to be 1 but was %d", len(fields))
	}
	if fields[0] != "RequiredField" {
		t.Errorf("expected invalid field to be 'RequiredField', but was '%s'", fields[0])
	}
}
