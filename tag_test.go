package mailchimp

import "testing"

func TestTagBuilder_Build(t *testing.T) {
	_, err := TagBuilder{}.Name("test").StatusActive().Build()
	if err != nil {
		t.Errorf(
			"expected no error to be returned, but got '%s'",
			err.Error(),
		)
	}
	_, err = TagBuilder{}.StatusInactive().Build()
	if err == nil {
		t.Error("expected error to be returned, but none was")
	}
	_, err = TagBuilder{}.Build()
	if err == nil {
		t.Error("expected error to be returned, but none was")
	}
}

func TestTagBuilder_BuildShouldPass(t *testing.T) {
	_, err := TagBuilder{}.Name("test").StatusActive().Build()
	if err != nil {
		t.Errorf(
			"expected no error to be returned, but got '%s'",
			err.Error(),
		)
	}
}

func TestTagBuilder_BuildWithoutNameReturnsError(t *testing.T) {
	_, err := TagBuilder{}.StatusInactive().Build()
	if err == nil {
		t.Error("expected error to be returned, but none was")
	}
}

func TestTagBuilder_BuildWithoutStatusReturnsError(t *testing.T) {
	_, err := TagBuilder{}.Name("test-tag").Build()
	if err == nil {
		t.Error("expected error to be returned, but none was")
	}
}

func TestTagBuilder_AddName(t *testing.T) {
	testName := "test-tag"
	builder := TagBuilder{}
	builder = builder.Name(testName)
	if builder.obj.Name != testName {
		t.Errorf(
			"expected name to be '%s' but was '%s'",
			testName,
			builder.obj.Name,
		)
	}
}

func TestTagBuilder_AddStatusActive(t *testing.T) {
	builder := TagBuilder{}
	builder = builder.StatusActive()
	if builder.obj.Status != tagStatusActive {
		t.Errorf(
			"expected status to be '%s' but was '%s'",
			tagStatusActive,
			builder.obj.Status,
		)
	}
}

func TestTagBuilder_AddStatusInactive(t *testing.T) {
	builder := TagBuilder{}
	builder = builder.StatusInactive()
	if builder.obj.Status != tagStatusInactive {
		t.Errorf(
			"expected status to be '%s' but was '%s'",
			tagStatusInactive,
			builder.obj.Status,
		)
	}
}
