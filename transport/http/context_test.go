package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextHeader(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}},
		res:    nil,
		w:      responseWriter{},
	}
	h := w.Header()
	assert.Equal(t, h, http.Header{"name": {"kratos"}})
}

func TestContextForm(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}, Method: "POST"},
		res:    nil,
		w:      responseWriter{},
	}
	form := w.Form()
	assert.Equal(t, form, url.Values{})

	w = wrapper{
		router: nil,
		req:    &http.Request{Form: map[string][]string{"name": {"kratos"}}},
		res:    nil,
		w:      responseWriter{},
	}
	form = w.Form()
	assert.Equal(t, form, url.Values{"name": []string{"kratos"}})
}

func TestContextQuery(t *testing.T) {
	w := wrapper{
		router: nil,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "github.com", Path: "go-kratos/kratos", RawQuery: "page=1"}, Method: "POST"},
		res:    nil,
		w:      responseWriter{},
	}
	q := w.Query()
	assert.Equal(t, q, url.Values{"page": []string{"1"}})
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
	assert.Equal(t, res, req)
}

func TestContextResponse(t *testing.T) {
	res := httptest.NewRecorder()
	w := wrapper{
		router: &Router{srv: &Server{enc: DefaultResponseEncoder}},
		req:    &http.Request{Method: "POST"},
		res:    res,
		w:      responseWriter{200, res},
	}
	assert.Equal(t, w.Response(), res)
	err := w.Returns(map[string]string{}, nil)
	assert.Nil(t, err)
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
	assert.Nil(t, err)
	assert.Equal(t, b, BindQuery{Page: 2})
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
	assert.Nil(t, err)
	assert.Equal(t, b, BindForm{Page: 2})
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
	assert.Nil(t, err)
	err = w.XML(200, "success")
	assert.Nil(t, err)
	err = w.String(200, "success")
	assert.Nil(t, err)
	err = w.Blob(200, "blob", []byte("success"))
	assert.Nil(t, err)
	err = w.Stream(200, "stream", bytes.NewBuffer([]byte("success")))
	assert.Nil(t, err)
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
	assert.Equal(t, ok, true)
	done := w.Done()
	assert.NotNil(t, done)
	err := w.Err()
	assert.Nil(t, err)
	v := w.Value("test")
	assert.Nil(t, v)

	w = wrapper{
		router: &Router{srv: &Server{enc: DefaultResponseEncoder}},
		req:    nil,
		res:    nil,
		w:      responseWriter{},
	}
	_, ok = w.Deadline()
	assert.Equal(t, ok, false)
	done = w.Done()
	assert.Nil(t, done)
	err = w.Err()
	assert.NotNil(t, err)
	v = w.Value("test")
	assert.Nil(t, v)
}
