package mailchimp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

type Client interface {
	Ping() (bool, error)
	CreateList(List) (List, error)
	FetchLists() (listCollection, error)
	FetchList(string) (List, error)
	UpdateList(string, List) (List, error)
	DeleteList(string) error
	Batch(string, []Member) error
	BatchWithUpdate(id string, members []Member) error
}

type client struct {
	key      string
	provider MailChimpProvider
}

type pingResponse struct {
	HealthStatus string `json:"health_status"`
}

type errorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func NewClient(key, region string) Client {
	return client{
		provider: mailChimpProvider{
			Region:        region,
			Authorization: authorization(key),
		},
	}
}

func NewMockClient(provider MailChimpProvider) Client {
	return client{
		provider: provider,
	}
}

func (c client) Ping() (bool, error) {
	var status pingResponse
	body, err := c.provider.Get("/ping")
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(body, &status); err != nil {
		return false, err
	}
	return status.HealthStatus == "Everything's Chimpy!", err
}

func (c client) CreateList(l List) (List, error) {
	list := List{}
	body, err := c.provider.Post("/lists", l)
	if err != nil {
		return List{}, err
	}
	if err := json.Unmarshal(body, &list); err != nil {
		return List{}, err
	}
	return list, nil
}

func (c client) FetchLists() (listCollection, error) {
	lists := listCollection{}
	body, err := c.provider.Get("/lists")
	if err != nil {
		return lists, err
	}
	if err := json.Unmarshal(body, &lists); err != nil {
		return lists, err
	}
	return lists, nil
}

func (c client) FetchList(id string) (List, error) {
	list := List{}
	body, err := c.provider.Get(fmt.Sprintf("/lists/%s", id))
	if err != nil {
		return List{}, err
	}
	if err := json.Unmarshal(body, &list); err != nil {
		return List{}, err
	}
	return list, nil
}

func (c client) UpdateList(id string, l List) (List, error) {
	body, err := c.provider.Patch(
		fmt.Sprintf("/lists/%s", id),
		l,
	)
	if err != nil {
		return List{}, err
	}
	var list List
	if err := json.Unmarshal(body, &list); err != nil {
		return list, err
	}
	return list, nil
}

func (c client) DeleteList(id string) error {
	_, err := c.provider.Delete(
		fmt.Sprintf("/lists/%s", id),
	)
	return err
}

type batchedMember struct {
	EmailAddress string            `json:"email_address"`
	Status       string            `json:"status"`
	MergeFields  map[string]string `json:"merge_fields"`
}

type batch struct {
	Members        []batchedMember `json:"members"`
	UpdateExisting bool            `json:"update_existing"`
}

func (c client) Batch(id string, members []Member) error {
	return c.batch(id, members, false)
}

func (c client) BatchWithUpdate(id string, members []Member) error {
	return c.batch(id, members, true)
}

func (c client) batch(id string, members []Member, update bool) error {
	if len(members) > 500 {
		return errors.New("batch operation only allows for a maximum of 500 members")
	}
	data := make([]batchedMember, 0)
	for _, member := range members {
		data = append(data, batchedMember{
			EmailAddress: member.EmailAddress,
			Status:       member.Status,
			MergeFields:  member.MergeFields,
		})
	}
	_, err := c.provider.Post(fmt.Sprintf("/lists/%s", id), batch{
		Members:        data,
		UpdateExisting: update,
	})
	return err
}

func authorization(key string) string {
	method := "Basic"
	k := base64.
		StdEncoding.
		EncodeToString(
			[]byte(
				fmt.Sprintf("anystring:%s", key),
			),
		)
	return fmt.Sprintf("%s %s", method, k)
}
