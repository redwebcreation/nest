package cloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/build"
	"io"
	"net/http"
)

type Client struct {
	client      *http.Client
	endpoint    string
	id          string
	accessToken string
}

var (
	// ErrResourceNotFound is returned when the requested resource is not found
	ErrResourceNotFound = errors.New("resource not found")
)

func (c Client) Request(method string, url string, params map[string]any, v any) error {
	request, err := http.NewRequest(method, c.endpoint+url, nil)
	if err != nil {
		return err
	}

	// send params as JSON in request body
	request.Header.Set("Accept", "application/json")

	if params != nil {
		request.Header.Set("Content-Type", "application/json")

		data, err := json.Marshal(params)
		if err != nil {
			return err
		}

		request.Body = io.NopCloser(bytes.NewReader(data))
	}

	response, err := c.client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusNotFound {
		return ErrResourceNotFound
	} else if response.StatusCode < 200 || response.StatusCode > 400 {
		return fmt.Errorf("invalid nest cloud response: status=%d path=%s", response.StatusCode, url)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &v)
}

func (c *Client) Ping() error {
	return c.Request("GET", "/servers/"+c.id+"/ping", nil, nil)
}

func NewClient(id string, accessToken string) *Client {
	return &Client{
		client:      &http.Client{},
		endpoint:    build.Endpoint,
		id:          id,
		accessToken: accessToken,
	}
}
