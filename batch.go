package mailchimp

import (
	"encoding/json"
	"fmt"
	"strings"
)

var NullOperation = Operation{}

type Operation struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}

type OperationCollection []Operation

func NewTagsOperation(listID, memberEmail string, tags []Tag) (Operation, error) {
	payload := updateMemberTagsPayload{
		Tags:      tags,
		IsSyncing: false,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return NullOperation, nil
	}
	return Operation{
		Method: "POST",
		Path: fmt.Sprintf(
			"/lists/%s/members/%s/tags",
			listID,
			hashMd5(strings.ToLower(memberEmail)),
		),
		Body: string(body),
	}, nil
}
