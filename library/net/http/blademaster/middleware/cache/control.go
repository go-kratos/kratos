package cache

import (
	fmt "fmt"
	"net/http"
	"sync"
	"time"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/cache/store"
)

const (
	_maxMaxAge = 60 * 5 // 5 minutes
)

// Control is used to work as client side cache orchestrator
type Control struct {
	MaxAge int32
	pool   sync.Pool
}

type controlWriter struct {
	*Control

	ctx      *bm.Context
	response http.ResponseWriter
}

var _ http.ResponseWriter = &controlWriter{}

// NewControl will create a new control cache struct
func NewControl(maxAge int32) *Control {
	if maxAge > _maxMaxAge {
		panic("MaxAge should be less than 300 seconds")
	}
	ctl := &Control{
		MaxAge: maxAge,
	}
	ctl.pool.New = func() interface{} {
		return &controlWriter{}
	}
	return ctl
}

// Key method is not needed in this situation
func (ctl *Control) Key(ctx *bm.Context) string { return "" }

// Handler is used to execute cache service
func (ctl *Control) Handler(_ store.Store) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		writer := ctl.pool.Get().(*controlWriter)
		writer.Control = ctl
		writer.ctx = ctx
		writer.response = ctx.Writer

		ctx.Writer = writer
		ctx.Next()

		ctl.pool.Put(writer)
	}
}

func (w *controlWriter) Header() http.Header                     { return w.response.Header() }
func (w *controlWriter) Write(data []byte) (size int, err error) { return w.response.Write(data) }
func (w *controlWriter) WriteHeader(code int) {
	// do not inject header if this is an error response
	if w.ctx.Error == nil {
		headers := w.Header()
		headers.Set("Expires", time.Now().UTC().Add(time.Duration(w.MaxAge)*time.Second).Format(http.TimeFormat))
		headers.Set("Cache-Control", fmt.Sprintf("max-age=%d", w.MaxAge))
	}
	w.response.WriteHeader(code)
}
