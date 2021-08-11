package mailchimp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ResponseStatusSuccess      = 2 // for 200, 201, 204 etc.
	ResponseStatusFailedClient = 4 // for 400, 401, 404 etc.
	ResponseStatusFailedServer = 5 // for 500
)

type errorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

type MailChimpProvider interface {
	Post(uri string, body interface{}) ([]byte, error)
	Get(uri string) ([]byte, error)
	Patch(uri string, body interface{}) ([]byte, error)
	Delete(uri string) ([]byte, error)
}

type mailChimpProvider struct {
	Authorization string
	Region        string
}

func (mcp mailChimpProvider) Post(uri string, body interface{}) ([]byte, error) {
	httpClient := http.DefaultClient
	req, err := mcp.createBodyRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode / 100 // Get the first digit of the status code
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if status != ResponseStatusSuccess {
		return nil, mcp.handleFailedRequest(bytes)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) Get(uri string) ([]byte, error) {
	httpClient := http.DefaultClient
	req, err := mcp.createBodylessRequest("GET", uri)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode / 100 // Get the first digit of the status code
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if status != ResponseStatusSuccess {
		return nil, mcp.handleFailedRequest(bytes)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) Patch(uri string, body interface{}) ([]byte, error) {
	httpClient := http.DefaultClient
	req, err := mcp.createBodyRequest("PATCH", uri, body)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode / 100 // Get the first digit of the status code
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if status != ResponseStatusSuccess {
		return nil, mcp.handleFailedRequest(bytes)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) Delete(uri string) ([]byte, error) {
	httpClient := http.DefaultClient
	req, err := mcp.createBodylessRequest("DELETE", uri)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode / 100 // Get the first digit of the status code
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if status != ResponseStatusSuccess {
		return nil, mcp.handleFailedRequest(bytes)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) createBodylessRequest(method, uri string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		mcp.url(uri),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", mcp.Authorization)
	return req, nil
}

func (mcp mailChimpProvider) createBodyRequest(method, uri string, body interface{}) (*http.Request, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		mcp.url(uri),
		bytes.NewBuffer(raw),
	)
	fmt.Println(string(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", mcp.Authorization)
	return req, nil
}

func (mcp mailChimpProvider) handleFailedRequest(body []byte) error {
	fmt.Println(string(body))
	errResponse := errorResponse{}
	err := json.Unmarshal(body, &errResponse)
	if err != nil {
		return errors.New(
			"request was not successful, and could not unmarshal error response",
		)
	}
	return fmt.Errorf(
		"request was not successful: %s Status %d",
		errResponse.Detail,
		errResponse.Status,
	)
}

func (mcp mailChimpProvider) url(uri string) string {
	return fmt.Sprintf(
		"https://%s.api.mailchimp.com/3.0%s",
		mcp.Region,
		uri,
	)
}
