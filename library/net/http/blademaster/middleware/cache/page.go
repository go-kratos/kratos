package cache

import (
	"bytes"
	"crypto/sha1"
	"io"
	"net/http"
	"net/url"
	"sync"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/cache/store"

	proto "github.com/gogo/protobuf/proto"
)

// consts for blademaster cache
const (
	_pagePrefix = "bm.page"
)

// Page is used to cache common response
type Page struct {
	Expire int32
	pool   sync.Pool
}

type cachedWriter struct {
	ctx      *bm.Context
	response http.ResponseWriter
	store    store.Store
	status   int
	expire   int32
	key      string
}

var _ http.ResponseWriter = &cachedWriter{}

// NewPage will create a new page cache struct
func NewPage(expire int32) *Page {
	pc := &Page{
		Expire: expire,
	}
	pc.pool.New = func() interface{} {
		return &cachedWriter{}
	}
	return pc
}

// Key is used to identify response cache key in most key-value store
func (p *Page) Key(ctx *bm.Context) string {
	url := ctx.Request.URL
	key := urlEscape(_pagePrefix, url.RequestURI())
	return key
}

// Handler is used to execute cache service
func (p *Page) Handler(store store.Store) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		var (
			resp   *ResponseCache
			cached []byte
			err    error
		)
		key := p.Key(ctx)
		cached, err = store.Get(ctx, key)

		// if we did got the previous cache,
		// try to unmarshal it
		if err == nil && len(cached) > 0 {
			resp = new(ResponseCache)
			err = proto.Unmarshal(cached, resp)
		}

		// if we failed to fetch the cache or failed to parse cached data,
		// then consider try to cache this response
		if err != nil || resp == nil {
			writer := p.pool.Get().(*cachedWriter)
			writer.ctx = ctx
			writer.response = ctx.Writer
			writer.key = key
			writer.expire = p.Expire
			writer.store = store

			ctx.Writer = writer
			ctx.Next()

			p.pool.Put(writer)
			return
		}

		// write cached response
		headers := ctx.Writer.Header()
		for key, value := range resp.Header {
			headers[key] = value.Value
		}
		ctx.Writer.WriteHeader(int(resp.Status))
		ctx.Writer.Write(resp.Data)
		ctx.Abort()
	}
}

func (w *cachedWriter) Header() http.Header {
	return w.response.Header()
}

func (w *cachedWriter) WriteHeader(code int) {
	w.status = int(code)
	w.response.WriteHeader(code)
}

func (w *cachedWriter) Write(data []byte) (size int, err error) {
	var (
		origin []byte
		pdata  []byte
	)
	if size, err = w.response.Write(data); err != nil {
		return
	}

	store := w.store
	origin, err = store.Get(w.ctx, w.key)
	resp := new(ResponseCache)
	if err == nil || len(origin) > 0 {
		err1 := proto.Unmarshal(origin, resp)
		if err1 == nil {
			data = append(resp.Data, data...)
		}
	}

	resp.Status = int32(w.status)
	resp.Header = headerValues(w.Header())
	resp.Data = data
	if pdata, err = proto.Marshal(resp); err != nil {
		// cannot happen
		log.Error("Failed to marshal response to protobuf: %v", err)
		return
	}

	if err = store.Set(w.ctx, w.key, pdata, w.expire); err != nil {
		log.Error("Failed to set response cache: %v", err)
		return
	}

	return
}

func headerValues(headers http.Header) map[string]*HeaderValue {
	result := make(map[string]*HeaderValue, len(headers))
	for key, values := range headers {
		result[key] = &HeaderValue{
			Value: values,
		}
	}
	return result
}

func urlEscape(prefix string, u string) string {
	key := url.QueryEscape(u)
	if len(key) > 200 {
		h := sha1.New()
		io.WriteString(h, u)
		key = string(h.Sum(nil))
	}
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}
