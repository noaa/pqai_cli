package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const defaultEndpoint = "https://api.projectpq.ai"

// Client is a thin HTTP client for the PQAI API.
type Client struct {
	Endpoint string
	Token    string
	HTTP     *http.Client
}

func NewClient() *Client {
	endpoint := os.Getenv("PQAI_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	return &Client{
		Endpoint: endpoint,
		Token:    os.Getenv("PQAI_API_KEY"),
		HTTP:     &http.Client{Timeout: 120 * time.Second},
	}
}

// Get sends a GET request to route with the given query parameters.
// If auth is true, the token is attached and its presence validated.
func (c *Client) Get(route string, params url.Values, auth bool) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	if auth {
		if c.Token == "" {
			return nil, fmt.Errorf("no API token found. Set the PQAI_API_KEY environment variable (see 'pqai config set-api-key')")
		}
		params.Set("token", c.Token)
	}
	u := c.Endpoint + route
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg := string(body)
		if len(msg) > 500 {
			msg = msg[:500] + "..."
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg)
	}
	return body, nil
}

// GetJSON runs Get and decodes the response into v.
func (c *Client) GetJSON(route string, params url.Values, auth bool, v any) error {
	body, err := c.Get(route, params, auth)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
