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
	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpClient := http.DefaultClient
	req, err := http.NewRequest(
		"POST",
		mcp.url(uri),
		bytes.NewBuffer(rawBody),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", mcp.Authorization)
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
		errResponse := errorResponse{}
		err = json.Unmarshal(bytes, &errResponse)
		return nil, errors.New(
			fmt.Sprintf(
				"request was not successful '%s', status: %d",
				errResponse.Detail,
				errResponse.Status,
			),
		)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) Get(uri string) ([]byte, error) {
	httpClient := http.DefaultClient
	req, err := http.NewRequest(
		"GET",
		mcp.url(uri),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", mcp.Authorization)
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
		errResponse := errorResponse{}
		err = json.Unmarshal(bytes, &errResponse)
		return nil, errors.New(
			fmt.Sprintf(
				"request was not successful '%s', status: %d",
				errResponse.Detail,
				errResponse.Status,
			),
		)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) Patch(uri string, body interface{}) ([]byte, error) {
	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpClient := http.DefaultClient
	req, err := http.NewRequest(
		"PATCH",
		mcp.url(uri),
		bytes.NewBuffer(rawBody),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", mcp.Authorization)
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
		errResponse := errorResponse{}
		err = json.Unmarshal(bytes, &errResponse)
		return nil, errors.New(
			fmt.Sprintf(
				"request was not successful '%s', status: %d",
				errResponse.Detail,
				errResponse.Status,
			),
		)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) Delete(uri string) ([]byte, error) {
	httpClient := http.DefaultClient
	req, err := http.NewRequest(
		"DELETE",
		mcp.url(uri),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", mcp.Authorization)
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
		errResponse := errorResponse{}
		err = json.Unmarshal(bytes, &errResponse)
		return nil, errors.New(
			fmt.Sprintf(
				"request was not successful '%s', status: %d",
				errResponse.Detail,
				errResponse.Status,
			),
		)
	}
	return bytes, nil
}

func (mcp mailChimpProvider) url(uri string) string {
	return fmt.Sprintf(
		"https://%s.api.mailchimp.com/3.0%s",
		mcp.Region,
		uri,
	)
}
