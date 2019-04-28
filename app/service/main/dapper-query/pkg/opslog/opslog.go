// Package opslog provide ops-log api
package opslog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	_kbnVersion        = "5.4.3"
	_indexPrefix       = "billions-"
	_mqueryContentType = "application/x-ndjson"
	_ajsSessioID       = "_AJSESSIONID"
)

// Err errors
var (
	ErrOverRange = errors.New("search time over range")
)

type response struct {
	Hits struct {
		Hits []struct {
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// Record represent log record
type Record struct {
	Time    time.Time              `json:"timestamp"`
	Fields  map[string]interface{} `json:"fields"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
}

// Client query log from ops-log
type Client interface {
	Query(ctx context.Context, familys []string, traceID uint64, sessionID string, start, end int64, options ...Option) ([]*Record, error)
}

type option struct {
	traceField string
	size       int
	level      string
}

var _defaultOpt = option{
	traceField: "traceid",
	size:       100,
}

// Option for query
type Option func(opt *option)

// SetTraceField default "traceid"
func SetTraceField(traceField string) Option {
	return func(opt *option) {
		opt.traceField = traceField
	}
}

// SetSize default 100
func SetSize(size int) Option {
	return func(opt *option) {
		opt.size = size
	}
}

// SetLevel return all if level is empty
func SetLevel(level string) Option {
	return func(opt *option) {
		opt.level = level
	}
}

// New ops-log client
func New(searchAPI string, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &client{
		searchAPI:  searchAPI,
		httpclient: httpClient,
	}
}

type client struct {
	searchAPI  string
	httpclient *http.Client
}

func (c *client) Query(ctx context.Context, familys []string, traceID uint64, sessionID string, start, end int64, options ...Option) ([]*Record, error) {
	if start <= 0 || end <= 0 {
		return nil, ErrOverRange
	}
	if len(familys) == 0 {
		return make([]*Record, 0), nil
	}
	opt := _defaultOpt
	for _, fn := range options {
		fn(&opt)
	}
	req, err := c.newReq(familys, traceID, sessionID, start, end, &opt)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "send request to %s fail", c.searchAPI)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		return nil, errors.Errorf("ops-log response error: status_code: %d, body: %s", resp.StatusCode, buf[:n])
	}
	return decodeRecord(resp.Body)
}

func (c *client) newReq(familys []string, traceID uint64, sessionID string, start, end int64, opt *option) (*http.Request, error) {
	prefixTraceID := strconv.FormatUint(traceID, 16)
	leagcyTraceID := strconv.FormatUint(traceID, 10)
	startMillis := start * int64((time.Second / time.Millisecond))
	endMillis := end * int64((time.Second / time.Millisecond))
	body := &bytes.Buffer{}
	enc := json.NewEncoder(body)
	header := map[string]interface{}{"index": formatIndices(familys), "ignore_unavailable": true}
	if err := enc.Encode(header); err != nil {
		return nil, err
	}
	shoulds := []map[string]interface{}{
		{"prefix": map[string]interface{}{opt.traceField: prefixTraceID}},
		{"match": map[string]interface{}{opt.traceField: leagcyTraceID}},
	}
	traceQuery := map[string]interface{}{"bool": map[string]interface{}{"should": shoulds}}
	rangeQuery := map[string]interface{}{
		"range": map[string]interface{}{
			"@timestamp": map[string]interface{}{"gte": startMillis, "lte": endMillis, "format": "epoch_millis"},
		},
	}
	musts := []map[string]interface{}{traceQuery, rangeQuery}
	if opt.level != "" {
		musts = append(musts, map[string]interface{}{"match": map[string]interface{}{"level": opt.level}})
	}
	query := map[string]interface{}{
		"sort": map[string]interface{}{
			"@timestamp": map[string]interface{}{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{"must": musts},
		},
		"version": true,
		"size":    opt.size,
	}
	if err := enc.Encode(query); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.searchAPI, body)
	if err != nil {
		return nil, err
	}
	session := &http.Cookie{Name: _ajsSessioID, Value: sessionID}
	req.AddCookie(session)
	req.Header.Set("Content-Type", _mqueryContentType)
	req.Header.Set("kbn-version", _kbnVersion)
	return req, nil
}

func decodeRecord(src io.Reader) ([]*Record, error) {
	var resp struct {
		Responses []response `json:"responses"`
	}
	if err := json.NewDecoder(src).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "decode response error")
	}
	if len(resp.Responses) == 0 {
		return nil, nil
	}
	records := make([]*Record, 0, len(resp.Responses[0].Hits.Hits))
	for _, hit := range resp.Responses[0].Hits.Hits {
		record := &Record{
			Fields: make(map[string]interface{}),
		}
		for k, v := range hit.Source {
			switch k {
			case "@timestamp":
				s, _ := v.(string)
				record.Time, _ = time.Parse(time.RFC3339Nano, s)
			case "log":
				s, _ := v.(string)
				record.Message = s
			case "level":
				s, _ := v.(string)
				record.Level = s
			default:
				record.Fields[k] = v
			}
		}
		records = append(records, record)
	}
	return records, nil
}

func formatIndices(familys []string) []string {
	indices := make([]string, len(familys))
	for i := range familys {
		indices[i] = _indexPrefix + familys[i] + "*"
	}
	return indices
}
