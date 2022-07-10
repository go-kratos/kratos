package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestContextHeader(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}},
		res:    nil,
		w:      responseWriter{},
	}
	h := w.Header()
	if !reflect.DeepEqual(h, http.Header{"name": {"kratos"}}) {
		t.Errorf("expected %v, got %v", http.Header{"name": {"kratos"}}, h)
	}
}

func TestContextForm(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}, Method: "POST"},
		res:    nil,
		w:      responseWriter{},
	}
	form := w.Form()
	if !reflect.DeepEqual(form, url.Values{}) {
		t.Errorf("expected %v, got %v", url.Values{}, form)
	}

	w = wrapper{
		router: nil,
		req:    &http.Request{Form: map[string][]string{"name": {"kratos"}}},
		res:    nil,
		w:      responseWriter{},
	}
	form = w.Form()
	if !reflect.DeepEqual(form, url.Values{"name": {"kratos"}}) {
		t.Errorf("expected %v, got %v", url.Values{"name": {"kratos"}}, form)
	}
}

func TestContextQuery(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "github.com", Path: "go-kratos/kratos", RawQuery: "page=1"}, Method: "POST"},
		res:    nil,
		w:      responseWriter{},
	}
	q := w.Query()
	if !reflect.DeepEqual(q, url.Values{"page": {"1"}}) {
		t.Errorf("expected %v, got %v", url.Values{"page": {"1"}}, q)
	}
}

func TestContextRequest(t *testing.T) {
	req := &http.Request{Method: "POST"}
	w := wrapper{
		router: nil,
		req:    req,
		res:    nil,
		w:      responseWriter{},
	}
	res := w.Request()
	if !reflect.DeepEqual(res, req) {
		t.Errorf("expected %v, got %v", req, res)
	}
}

func TestContextResponse(t *testing.T) {
	res := httptest.NewRecorder()
	w := wrapper{
		router: &Router{srv: &Server{enc: DefaultResponseEncoder}},
		req:    &http.Request{Method: "POST"},
		res:    res,
		w:      responseWriter{200, res},
	}
	if !reflect.DeepEqual(w.Response(), res) {
		t.Errorf("expected %v, got %v", res, w.Response())
	}
	err := w.Returns(map[string]string{}, nil)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	needErr := fmt.Errorf("some error")
	err = w.Returns(map[string]string{}, needErr)
	if !errors.Is(err, needErr) {
		t.Errorf("expected %v, got %v", needErr, err)
	}
}

func TestContextBindQuery(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "go-kratos-dev", RawQuery: "page=2"}},
		res:    nil,
		w:      responseWriter{},
	}
	type BindQuery struct {
		Page int `json:"page"`
	}
	b := BindQuery{}
	err := w.BindQuery(&b)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(b, BindQuery{Page: 2}) {
		t.Errorf("expected %v, got %v", BindQuery{Page: 2}, b)
	}
}

func TestContextBindForm(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "go-kratos-dev"}, Form: map[string][]string{"page": {"2"}}},
		res:    nil,
		w:      responseWriter{},
	}
	type BindForm struct {
		Page int `json:"page"`
	}
	b := BindForm{}
	err := w.BindForm(&b)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(b, BindForm{Page: 2}) {
		t.Errorf("expected %v, got %v", BindForm{Page: 2}, b)
	}
}

func TestContextResponseReturn(t *testing.T) {
	writer := httptest.NewRecorder()
	w := wrapper{
		router: nil,
		req:    nil,
		res:    writer,
		w:      responseWriter{},
	}
	err := w.JSON(200, "success")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	err = w.XML(200, "success")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	err = w.String(200, "success")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	err = w.Blob(200, "blob", []byte("success"))
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	err = w.Stream(200, "stream", bytes.NewBuffer([]byte("success")))
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
}

func TestContextCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &http.Request{Method: "POST"}
	req = req.WithContext(ctx)
	w := wrapper{
		router: &Router{srv: &Server{enc: DefaultResponseEncoder}},
		req:    req,
		res:    nil,
		w:      responseWriter{},
	}
	_, ok := w.Deadline()
	if !ok {
		t.Errorf("expected %v, got %v", true, ok)
	}
	done := w.Done()
	if done == nil {
		t.Errorf("expected %v, got %v", true, ok)
	}
	err := w.Err()
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	v := w.Value("test")
	if v != nil {
		t.Errorf("expected %v, got %v", nil, v)
	}

	w = wrapper{
		router: &Router{srv: &Server{enc: DefaultResponseEncoder}},
		req:    nil,
		res:    nil,
		w:      responseWriter{},
	}
	_, ok = w.Deadline()
	if ok {
		t.Errorf("expected %v, got %v", false, ok)
	}
	done = w.Done()
	if done != nil {
		t.Errorf("expected not nil, got %v", done)
	}
	err = w.Err()
	if err == nil {
		t.Errorf("expected not %v, got %v", nil, err)
	}
	v = w.Value("test")
	if v != nil {
		t.Errorf("expected %v, got %v", nil, v)
	}
}
