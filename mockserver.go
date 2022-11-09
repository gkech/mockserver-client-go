package mockserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gkech/mockserver-client-go/create"
	"github.com/gkech/mockserver-client-go/verify"
)

// HttpClient is a HTTP client interface.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a mockserver for mockserver.
type Client struct {
	host   string
	client HttpClient
}

// NewClient creates a new mockserver.
func NewClient(address string) *Client {
	if !strings.HasPrefix(address, "http") {
		address = "http://" + address
	}
	return &Client{
		host:   address,
		client: http.DefaultClient,
	}
}

// CreateExpectation creates a new expectation (request/response)
// in the mockserver. If the expectation is received successfully,
// the response will be a 201 HTTP status code.
func (c *Client) CreateExpectation(expectation create.Expectation) error {
	body, err := json.Marshal(expectation)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/mockserver/expectation", c.host)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"error creating mockserver expectation. error code: %d and body: %s", resp.StatusCode, resp.Body)
	}

	return nil
}

// VerifyRequest to verify a specific request was made to the mockserver.
// If the matching request was received the specified number of times,
// the response will be a 202 HTTP status code. If the request was not
// received the designated times, a 406 HTTP status code is issued.
func (c *Client) VerifyRequest(expectation verify.Expectation) error {
	body, err := json.Marshal(expectation)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/mockserver/verify", c.host)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"error verifying mockserver request. error code: %d and body: %s", resp.StatusCode, resp.Body)
	}

	return nil
}

// Reset to reset all request expectations and logs from the mockserver.
func (c *Client) Reset() error {
	url := fmt.Sprintf("%s/mockserver/reset", c.host)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"error resetting mockserver. error code: %d and body: %s", resp.StatusCode, resp.Body)
	}

	return nil
}

type mockserverResponse struct {
	Body       string
	StatusCode int
}

func (c *Client) sendRequest(req *http.Request) (*mockserverResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		fmt.Printf("ERROR: %v", err)
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &mockserverResponse{
		Body:       string(b),
		StatusCode: resp.StatusCode,
	}, err
}
