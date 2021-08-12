package mailchimp

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Client interface {
	Ping() error

	CreateList(List) (List, error)
	FetchLists() ([]List, error)
	FetchList(string) (List, error)
	UpdateList(string, List) (List, error)
	DeleteList(string) error

	Batch(string, []Member) error
	BatchWithUpdate(id string, members []Member) error

	FetchMemberTags(listID, memberEmail string) ([]Tag, error)
	UpdateMemberTags(listID, memberEmail string, tags []Tag) error
	UpdateMemberTagsSync(listID, memberEmail string, tags []Tag) error

	CreateWebhook(webhook Webhook) (Webhook, error)
	FetchWebhooks(listID string) ([]Webhook, error)
	FetchWebhook(listID string, webookID string) (Webhook, error)
	DeleteWebhook(listID string, webhookID string) error
}

type client struct {
	provider MailChimpProvider
}

type pingResponse struct {
	HealthStatus string `json:"health_status"`
}

func NewClient(key, region string) Client {
	return client{
		provider: mailChimpProvider{
			Region:        region,
			Authorization: authorization(key),
		},
	}
}

func NewCustomDependencyClient(provider MailChimpProvider) Client {
	return client{
		provider: provider,
	}
}

func (c client) Ping() error {
	var status pingResponse
	body, err := c.provider.Get("/ping")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &status); err != nil {
		return err
	}
	if status.HealthStatus != "Everything's Chimpy!" {
		return errors.New("unexpected pong response from MailChimp API")
	}
	return nil
}

func (c client) CreateList(l List) (List, error) {
	body, err := c.provider.Post("/lists", l)
	if err != nil {
		return NullList, err
	}
	list := List{}
	if err := json.Unmarshal(body, &list); err != nil {
		return NullList, err
	}
	return list, nil
}

func (c client) FetchLists() ([]List, error) {
	body, err := c.provider.Get("/lists")
	if err != nil {
		return NullListSlice, err
	}
	lists := listCollection{}
	if err := json.Unmarshal(body, &lists); err != nil {
		return NullListSlice, err
	}
	return lists.Lists, nil
}

func (c client) FetchList(id string) (List, error) {
	body, err := c.provider.Get(fmt.Sprintf("/lists/%s", id))
	if err != nil {
		return NullList, err
	}
	list := List{}
	if err := json.Unmarshal(body, &list); err != nil {
		return NullList, err
	}
	return list, nil
}

func (c client) UpdateList(id string, l List) (List, error) {
	body, err := c.provider.Patch(
		fmt.Sprintf("/lists/%s", id),
		l,
	)
	if err != nil {
		return NullList, err
	}
	list := List{}
	if err := json.Unmarshal(body, &list); err != nil {
		return NullList, err
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

type memberTagsResponse struct {
	Tags []Tag `json:"tags"`
}

func (c client) FetchMemberTags(listID, memberEmail string) ([]Tag, error) {
	tags := memberTagsResponse{}
	body, err := c.provider.Get(
		fmt.Sprintf(
			"/lists/%s/members/%s/tags",
			listID,
			hashMd5(strings.ToLower(memberEmail)),
		),
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, err
	}
	for i := range tags.Tags {
		tags.Tags[i].Status = tagStatusActive
	}
	return tags.Tags, nil
}

type updateMemberTagsPayload struct {
	Tags      []Tag `json:"tags"`
	IsSyncing bool  `json:"is_syncing"`
}

func (c client) UpdateMemberTags(listID, memberEmail string, tags []Tag) error {
	_, err := c.provider.Post(
		fmt.Sprintf(
			"/lists/%s/members/%s/tags",
			listID,
			hashMd5(strings.ToLower(memberEmail)),
		),
		updateMemberTagsPayload{
			Tags:      tags,
			IsSyncing: false,
		},
	)
	return err
}

func (c client) UpdateMemberTagsSync(listID, memberEmail string, tags []Tag) error {
	_, err := c.provider.Post(
		fmt.Sprintf(
			"/lists/%s/members/%s/tags",
			listID,
			hashMd5(strings.ToLower(memberEmail)),
		),
		updateMemberTagsPayload{
			Tags:      tags,
			IsSyncing: true,
		},
	)
	return err
}

type CreateWebhookRequestPayload struct {
	URL     string         `json:"url"`
	Events  WebhookEvents  `json:"events"`
	Sources WebhookSources `json:"sources"`
}

func (c client) CreateWebhook(webhook Webhook) (Webhook, error) {
	body, err := c.provider.Post(
		fmt.Sprintf("/lists/%s/webhooks", webhook.ListID),
		CreateWebhookRequestPayload{
			URL:     webhook.URL,
			Events:  webhook.Events,
			Sources: webhook.Sources,
		},
	)
	if err != nil {
		return NullWebhook, err
	}
	createdWebhook := Webhook{}
	if err := json.Unmarshal(body, &createdWebhook); err != nil {
		return NullWebhook, err
	}
	return createdWebhook, nil
}

func (c client) FetchWebhooks(listID string) ([]Webhook, error) {
	body, err := c.provider.Get(
		fmt.Sprintf("/lists/%s/webhooks", listID),
	)
	if err != nil {
		return nil, err
	}
	collection := webhookCollection{}
	if err := json.Unmarshal(body, &collection); err != nil {
		return nil, err
	}
	return collection.Webhooks, nil
}

func (c client) FetchWebhook(listID, webhookID string) (Webhook, error) {
	body, err := c.provider.Get(
		fmt.Sprintf(
			"/lists/%s/webhooks/%s",
			listID,
			webhookID,
		),
	)
	if err != nil {
		return NullWebhook, err
	}
	webhook := Webhook{}
	if err := json.Unmarshal(body, &webhook); err != nil {
		return NullWebhook, err
	}
	return webhook, nil
}

func (c client) DeleteWebhook(listID, webhookID string) error {
	_, err := c.provider.Delete(
		fmt.Sprintf(
			"/lists/%s/webhooks/%s",
			listID,
			webhookID,
		),
	)
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

func hashMd5(data string) string {
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}
