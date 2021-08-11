package mailchimp

type MailChimpProviderMock struct {
	PostMock    func(string, interface{}) ([]byte, error)
	PostCalls   int
	GetMock     func(string) ([]byte, error)
	GetCalls    int
	PatchMock   func(string, interface{}) ([]byte, error)
	PatchCalls  int
	DeleteMock  func(string) ([]byte, error)
	DeleteCalls int
}

func (mcpm *MailChimpProviderMock) Post(uri string, body interface{}) ([]byte, error) {
	mcpm.PostCalls++
	return mcpm.PostMock(uri, body)
}

func (mcpm *MailChimpProviderMock) Get(uri string) ([]byte, error) {
	mcpm.GetCalls++
	return mcpm.GetMock(uri)
}

func (mcpm *MailChimpProviderMock) Patch(uri string, body interface{}) ([]byte, error) {
	mcpm.PatchCalls++
	return mcpm.PatchMock(uri, body)
}

func (mcpm *MailChimpProviderMock) Delete(uri string) ([]byte, error) {
	mcpm.DeleteCalls++
	return mcpm.DeleteMock(uri)
}

type ClientMock struct {
	PingMock  func() error
	PingCalls int

	CreateListMock  func(List) (List, error)
	CreateListCalls int
	FetchListsMock  func() ([]List, error)
	FetchListsCalls int
	FetchListMock   func(string) (List, error)
	FetchListCalls  int
	UpdateListMock  func(string, List) (List, error)
	UpdateListCalls int
	DeleteListMock  func(string) error
	DeleteListCalls int

	BatchMock            func(string, []Member) error
	BatchCalls           int
	BatchWithUpdateMock  func(string, []Member) error
	BatchWithUpdateCalls int

	FetchMemberTagsMock       func(string, string) ([]Tag, error)
	FetchMemberTagsCalls      int
	UpdateMemberTagsMock      func(string, string, []Tag) error
	UpdateMemberTagsCalls     int
	UpdateMemberTagsSyncMock  func(string, string, []Tag) error
	UpdateMemberTagsSyncCalls int

	CreateWebhookMock  func(webhook Webhook) (Webhook, error)
	CreateWebhookCalls int
	FetchWebhookMock   func(listID string, webhookID string) (Webhook, error)
	FetchWebhookCalls  int
	DeleteWebhookMock  func(listID string, webhookID string) error
	DeleteWebhookCalls int
}

func (client *ClientMock) Ping() error {
	client.PingCalls++
	return client.PingMock()
}

func (client *ClientMock) CreateList(list List) (List, error) {
	client.CreateListCalls++
	return client.CreateListMock(list)
}

func (client *ClientMock) FetchLists() ([]List, error) {
	client.FetchListsCalls++
	return client.FetchListsMock()
}

func (client *ClientMock) FetchList(id string) (List, error) {
	client.FetchListCalls++
	return client.FetchListMock(id)
}

func (client *ClientMock) UpdateList(id string, list List) (List, error) {
	client.UpdateListCalls++
	return client.UpdateListMock(id, list)
}

func (client *ClientMock) DeleteList(id string) error {
	client.DeleteListCalls++
	return client.DeleteListMock(id)
}

func (client *ClientMock) Batch(id string, members []Member) error {
	client.BatchCalls++
	return client.BatchMock(id, members)
}

func (client *ClientMock) BatchWithUpdate(id string, members []Member) error {
	client.BatchWithUpdateCalls++
	return client.BatchWithUpdateMock(id, members)
}

func (client *ClientMock) FetchMemberTags(id, memberEmail string) ([]Tag, error) {
	client.FetchMemberTagsCalls++
	return client.FetchMemberTagsMock(id, memberEmail)
}

func (client *ClientMock) UpdateMemberTags(id, memberEmail string, tags []Tag) error {
	client.UpdateMemberTagsCalls++
	return client.UpdateMemberTagsMock(id, memberEmail, tags)
}

func (client *ClientMock) UpdateMemberTagsSync(id, memberEmail string, tags []Tag) error {
	client.UpdateMemberTagsSyncCalls++
	return client.UpdateMemberTagsSyncMock(id, memberEmail, tags)
}

func (mock *ClientMock) CreaterWebhook(webhook Webhook) (Webhook, error) {
	mock.CreateWebhookCalls++
	return mock.CreateWebhookMock(webhook)
}

func (mock *ClientMock) FetchWebhook(listID, webhookID string) (Webhook, error) {
	mock.FetchWebhookCalls++
	return mock.FetchWebhookMock(listID, webhookID)
}

func (mock *ClientMock) DeleteWebhook(listID, webhookID string) error {
	mock.DeleteWebhookCalls++
	return mock.DeleteWebhookMock(listID, webhookID)
}
