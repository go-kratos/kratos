// Package cltclient provide fetch and merge data from collector
package cltclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"go-common/library/sync/errgroup"
)

const (
	_jsonMime = "application/json"
)

// ClientStatusResp response clientstatus request just for debug
type ClientStatusResp struct {
	QueueLen int             `json:"queue_len"`
	Clients  []*ClientStatus `json:"clients"`
}

// ClientStatus client status
type ClientStatus struct {
	Addr     string `json:"addr"`
	UpTime   int64  `json:"up_time"`
	ErrCount int64  `json:"err_count"`
	Rate     int64  `json:"rate"`
}

// CltStatus collector status
type CltStatus struct {
	Node     string          `json:"node"`
	QueueLen int             `json:"queue_len"`
	Clients  []*ClientStatus `json:"clients"`
}

// New collector client
func New(nodes []string, httpclient *http.Client) (*Client, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no node provided")
	}
	if httpclient == nil {
		httpclient = http.DefaultClient
	}
	return &Client{nodes: nodes, httpclient: httpclient}, nil
}

// Client collector client
type Client struct {
	nodes      []string
	httpclient *http.Client
}

// Status return all collector status
func (c *Client) Status(ctx context.Context) ([]*CltStatus, error) {
	var mx sync.Mutex
	var g errgroup.Group
	results := make(map[string]*ClientStatusResp)
	for _, node := range c.nodes {
		node := node // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			resp, err := c.fetchStatus(ctx, node)
			if err != nil {
				return err
			}
			mx.Lock()
			results[node] = resp
			mx.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	clts := make([]*CltStatus, 0, len(results))
	for node, resp := range results {
		clts = append(clts, &CltStatus{
			Node:     node,
			QueueLen: resp.QueueLen,
			Clients:  resp.Clients,
		})
	}
	return clts, nil
}

func (c *Client) fetchStatus(ctx context.Context, node string) (*ClientStatusResp, error) {
	var wrapResp struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Data    ClientStatusResp `json:"data"`
	}
	reqURL := "http://" + node + "/x/internal/dapper-collector/client-status"
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", _jsonMime)
	req = req.WithContext(ctx)
	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		p := make([]byte, 2048)
		n, _ := resp.Body.Read(p)
		return nil, fmt.Errorf("request url: %s status code: %d, body: %s", reqURL, resp.StatusCode, p[:n])
	}
	err = json.NewDecoder(resp.Body).Decode(&wrapResp)
	return &wrapResp.Data, err
}
