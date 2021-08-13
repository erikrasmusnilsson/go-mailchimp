package mailchimp

import (
	"fmt"
)

const (
	tagStatusActive   = "active"
	tagStatusInactive = "inactive"
)

type Tag struct {
	Name   string `json:"name" mc_validator:"required"`
	Status string `json:"status" mc_validator:"required"`
}

type TagBuilder struct {
	obj Tag
}

func (tb TagBuilder) Build() (Tag, error) {
	if invalidParams, valid := validate(tb.obj); valid {
		return tb.obj, nil
	} else {
		return tb.obj, fmt.Errorf(
			"could not build tag due to invalid parameters %v",
			invalidParams,
		)
	}
}

func (tb TagBuilder) Name(name string) TagBuilder {
	tb.obj.Name = name
	return tb
}

func (tb TagBuilder) StatusActive() TagBuilder {
	tb.obj.Status = tagStatusActive
	return tb
}

func (tb TagBuilder) StatusInactive() TagBuilder {
	tb.obj.Status = tagStatusInactive
	return tb
}
