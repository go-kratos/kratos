package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var testRouter = &Router{srv: NewServer()}

func TestContextHeader(t *testing.T) {
	w := wrapper{
		router: testRouter,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}},
		res:    nil,
	}
	h := w.Header()
	if !reflect.DeepEqual(h, http.Header{"name": {"kratos"}}) {
		t.Errorf("expected %v, got %v", http.Header{"name": {"kratos"}}, h)
	}
}

func TestContextForm(t *testing.T) {
	w := wrapper{
		router: testRouter,
		req:    &http.Request{Header: map[string][]string{"name": {"kratos"}}, Method: http.MethodPost},
		res:    nil,
	}
	form := w.Form()
	if !reflect.DeepEqual(form, url.Values{}) {
		t.Errorf("expected %v, got %v", url.Values{}, form)
	}

	w = wrapper{
		router: testRouter,
		req:    &http.Request{Form: map[string][]string{"name": {"kratos"}}},
		res:    nil,
	}
	form = w.Form()
	if !reflect.DeepEqual(form, url.Values{"name": {"kratos"}}) {
		t.Errorf("expected %v, got %v", url.Values{"name": {"kratos"}}, form)
	}
}

func TestContextQuery(t *testing.T) {
	w := wrapper{
		router: testRouter,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "github.com", Path: "go-kratos/kratos", RawQuery: "page=1"}, Method: http.MethodPost},
		res:    nil,
	}
	q := w.Query()
	if !reflect.DeepEqual(q, url.Values{"page": {"1"}}) {
		t.Errorf("expected %v, got %v", url.Values{"page": {"1"}}, q)
	}
}

func TestContextRequest(t *testing.T) {
	req := &http.Request{Method: http.MethodPost}
	w := wrapper{
		router: testRouter,
		req:    req,
		res:    nil,
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
		req:    &http.Request{Method: http.MethodPost},
		res:    res,
	}
	if !reflect.DeepEqual(w.Response(), res) {
		t.Errorf("expected %v, got %v", res, w.Response())
	}
	err := w.Returns(map[string]string{}, nil)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	needErr := errors.New("some error")
	err = w.Returns(map[string]string{}, needErr)
	if !errors.Is(err, needErr) {
		t.Errorf("expected %v, got %v", needErr, err)
	}
}

func TestContextResult(t *testing.T) {
	testCases := []struct {
		name   string
		enc    EncodeResponseFunc
		code   int
		header string
	}{
		{
			name: "normal",
			enc: func(rw http.ResponseWriter, r *http.Request, v interface{}) error {
				rw.Header().Set("X-Foo", "foo")
				_, err := rw.Write([]byte(v.(string)))
				return err
			},
			code:   400,
			header: "foo",
		},
		{
			name: "writeHeader",
			enc: func(rw http.ResponseWriter, r *http.Request, v interface{}) error {
				rw.Header().Set("X-Foo", "foo")
				rw.WriteHeader(500)
				_, err := rw.Write([]byte(v.(string)))
				return err
			},
			code:   500,
			header: "foo",
		},
		{
			name: "unwrap",
			enc: func(rw http.ResponseWriter, r *http.Request, v interface{}) error {
				u, ok := rw.(interface{ Unwrap() http.ResponseWriter })
				if !ok {
					t.Fatal("can not wrap http.ResponseWriter")
				}

				w := u.Unwrap()
				w.Header().Set("X-Foo", "foo")
				w.WriteHeader(500)

				_, err := w.Write([]byte(v.(string)))
				return err
			},
			code:   500,
			header: "foo",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			w := wrapper{
				router: &Router{srv: &Server{enc: tc.enc}},
				req:    nil,
				res:    res,
			}
			err := w.Result(400, "body")
			if err != nil {
				t.Fatalf("expected %v, got %v", nil, err)
			}

			resp := res.Result()
			if resp.StatusCode != tc.code {
				t.Fatalf("expected %d, got %d", tc.code, resp.StatusCode)
			}
			if s := resp.Header.Get("X-Foo"); s != tc.header {
				t.Fatalf("expected %q, got %q", tc.header, s)
			}
			if bs, _ := io.ReadAll(res.Body); string(bs) != "body" {
				t.Fatalf("expected %s, got: %s", "body", bs)
			}
		})
	}
}

func TestContextBindQuery(t *testing.T) {
	w := wrapper{
		router: testRouter,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "go-kratos-dev", RawQuery: "page=2"}},
		res:    nil,
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
		router: testRouter,
		req:    &http.Request{URL: &url.URL{Scheme: "https", Host: "go-kratos-dev"}, Form: map[string][]string{"page": {"2"}}},
		res:    nil,
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
		router: testRouter,
		req:    nil,
		res:    writer,
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
	req := &http.Request{Method: http.MethodPost}
	req = req.WithContext(ctx)
	w := wrapper{
		router: testRouter,
		req:    req,
		res:    nil,
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
