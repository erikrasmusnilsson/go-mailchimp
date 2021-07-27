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
