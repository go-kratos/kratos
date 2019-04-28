package proxy

import (
	"bytes"
	"io"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go-common/library/conf/env"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

type endpoint struct {
	url       *url.URL
	proxy     *httputil.ReverseProxy
	condition func(ctx *bm.Context) bool
}
type logger struct{}

func (logger) Write(p []byte) (int, error) {
	log.Warn("%s", string(p))
	return len(p), nil
}

func newep(rawurl string, condition func(ctx *bm.Context) bool) *endpoint {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(errors.Errorf("Invalid URL: %s", rawurl))
	}
	e := &endpoint{
		url: u,
	}
	e.proxy = &httputil.ReverseProxy{
		Director: e.director,
		ErrorLog: stdlog.New(logger{}, "bm.proxy: ", stdlog.LstdFlags),
	}
	e.condition = condition
	return e
}

func (e *endpoint) director(req *http.Request) {
	req.URL.Scheme = e.url.Scheme
	req.URL.Host = e.url.Host
	// keep the origin request path
	if e.url.Path != "" {
		req.URL.Path = e.url.Path
	}

	body, length := rebuildBody(req)
	req.Body = body
	req.ContentLength = int64(length)
}

func (e *endpoint) ServeHTTP(ctx *bm.Context) {
	req := ctx.Request
	ip := metadata.String(ctx, metadata.RemoteIP)
	logArgs := []log.D{
		log.KV("method", req.Method),
		log.KV("ip", ip),
		log.KV("path", req.URL.Path),
		log.KV("params", req.Form.Encode()),
	}
	if !e.condition(ctx) {
		logArgs = append(logArgs, log.KV("proxied", "false"))
		log.Infov(ctx, logArgs...)
		return
	}
	logArgs = append(logArgs, log.KV("proxied", "true"))
	log.Infov(ctx, logArgs...)
	e.proxy.ServeHTTP(ctx.Writer, ctx.Request)
	ctx.Abort()
}

func rebuildBody(req *http.Request) (io.ReadCloser, int) {
	// GET request
	if req.Body == nil {
		return nil, 0
	}

	// Submit with form
	if len(req.PostForm) > 0 {
		br := bytes.NewReader([]byte(req.PostForm.Encode()))
		return ioutil.NopCloser(br), br.Len()
	}

	// copy the original body
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	br := bytes.NewReader(bodyBytes)
	return ioutil.NopCloser(br), br.Len()
}

func always(ctx *bm.Context) bool {
	return true
}

// NewZoneProxy is
func NewZoneProxy(matchZone, dst string) bm.HandlerFunc {
	ep := newep(dst, func(*bm.Context) bool {
		if env.Zone == matchZone {
			return true
		}
		return false
	})
	return ep.ServeHTTP
}

// NewAlways is
func NewAlways(dst string) bm.HandlerFunc {
	ep := newep(dst, always)
	return ep.ServeHTTP
}
