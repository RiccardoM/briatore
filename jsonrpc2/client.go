package jsonrpc2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
)

// Client represents a JSON-RPC client
type Client struct {
	url        string
	httpClient *http.Client
}

// NewClient creates a new Client instance
func NewClient(url string, httpClient *http.Client) (*Client, error) {
	if _, err := urlpkg.Parse(url); err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	return &Client{
		url:        url,
		httpClient: httpClient,
	}, nil
}

// Call allows to perform an RPC call to the given method with the provided params
func (c *Client) Call(ctx context.Context, method string, params any, result any) error {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error while unmarshalling params: %w", err)
	}

	req := NewRequest(-1, method, paramsJSON)
	reqJSON, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(reqJSON))
	if err != nil {
		return fmt.Errorf("error while performing http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := c.httpClient.Do(httpReq)

	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}()

	var resp Response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return fmt.Errorf("error while unmarshalling response: status code %d: %w", httpResp.StatusCode, err)
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: status code %d: %w", httpResp.StatusCode, resp.Error)
	}

	err = json.Unmarshal(resp.Result, result)
	if err != nil {
		return fmt.Errorf("unmarshal result: %w", err)
	}

	return nil
}
