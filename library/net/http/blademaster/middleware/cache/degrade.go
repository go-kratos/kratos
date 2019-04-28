package cache

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/cache/store"
)

const (
	_degradeInterval = 60 * 10
	_degradePrefix   = "bm.degrade"
)

var (
	_degradeBytes = []byte(fmt.Sprintf("{\"code\":%d, \"message\":\"\"}", ecode.Degrade))
)

// Degrader is the common degrader instance.
type Degrader struct {
	lock sync.RWMutex
	urls map[string]*state

	expire int32
	ch     chan *result
	pool   sync.Pool // degradeWriter pool
}

// argsDegrader means the degrade will happened by args policy
type argsDegrader struct {
	*Degrader

	args []string
}

type degradeWriter struct {
	*Degrader

	ctx      *bm.Context
	response http.ResponseWriter
	store    store.Store
	key      string
	state    *state
}

type state struct {
	// FIXME(zhoujiahui): using transient map to avoid potential memory leak?
	// record last cached time
	sync.RWMutex
	gens map[string]*int64
}

type result struct {
	key   string
	value []byte
	store store.Store
}

var _ http.ResponseWriter = &degradeWriter{}
var _ Policy = &argsDegrader{}

// NewDegrader will create a new degrade struct
func NewDegrader(expire int32) (d *Degrader) {
	d = &Degrader{
		urls:   make(map[string]*state),
		ch:     make(chan *result, 1024),
		expire: expire,
	}
	d.pool.New = func() interface{} {
		return &degradeWriter{
			Degrader: d,
		}
	}

	go d.degradeproc()
	return
}

func (d *Degrader) degradeproc() {
	for {
		r := <-d.ch
		if err := r.store.Set(context.Background(), r.key, r.value, d.expire); err != nil {
			log.Error("store write key(%s) error(%v)", r.key, err)
		}
	}
}

// Args means this path will be degrade by specified args
func (d *Degrader) Args(args ...string) Policy {
	return &argsDegrader{
		Degrader: d,
		args:     args,
	}
}

func (d *Degrader) state(path string) *state {
	d.lock.RLock()
	s, ok := d.urls[path]
	d.lock.RUnlock()
	if !ok {
		s = &state{
			gens: make(map[string]*int64),
		}
		d.lock.Lock()
		d.urls[path] = s
		d.lock.Unlock()
	}
	return s
}

// Key is used to identify response cache key in most key-value store
func (ad *argsDegrader) Key(ctx *bm.Context) string {
	req := ctx.Request
	path := req.URL.Path
	params := req.Form

	vs := make([]string, 0, len(ad.args))
	for _, arg := range ad.args {
		vs = append(vs, params.Get(arg))
	}
	return fmt.Sprintf("%s:%s_%x", _degradePrefix, strings.Replace(path, "/", "_", -1), md5.Sum([]byte(strings.Join(vs, "-"))))
}

// Handler is used to execute degrade service
func (ad *argsDegrader) Handler(store store.Store) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		req := ctx.Request
		path := req.URL.Path

		writer := ad.pool.Get().(*degradeWriter)
		writer.response = ctx.Writer
		writer.ctx = ctx
		writer.store = store
		writer.state = ad.state(path)
		writer.key = ad.Key(ctx)

		ctx.Writer = writer // replace to degrade writer
		ctx.Next()

		ad.pool.Put(writer)
	}
}

func (w *degradeWriter) Header() http.Header  { return w.response.Header() }
func (w *degradeWriter) WriteHeader(code int) { w.response.WriteHeader(code) }

func (w *degradeWriter) Write(data []byte) (size int, err error) {
	e := w.ctx.Error
	// if an degrade error code is raised from upstream,
	// degrade this request directly
	if e != nil {
		if ec := ecode.Cause(e); ec.Code() == ecode.Degrade.Code() {
			return w.write()
		}
	}

	// write origin response
	if size, err = w.response.Write(data); err != nil {
		return
	}

	// error raised, this is a unsuccessful response
	if e != nil {
		return
	}

	// is required to cache
	if !w.state.required(w.key) {
		return
	}

	// async cache succeeded response for further degradation
	select {
	case w.ch <- &result{key: w.key, value: data, store: w.store}:
	default:
	}

	return
}

func (w *degradeWriter) write() (int, error) {
	data, err := w.store.Get(w.ctx, w.key)
	if err != nil || len(data) == 0 {
		// FIXME(zhoujiahui): The default response data should be respect to render type or content-type header
		data = _degradeBytes
	}
	return w.response.Write(data)
}

// check is required to cache response
// it depends on last cache time and _degradeInterval
func (st *state) required(key string) bool {
	now := time.Now().Unix()

	st.RLock()
	pLast, ok := st.gens[key]
	st.RUnlock()
	if !ok {
		st.Lock()
		pLast = new(int64)
		st.gens[key] = pLast
		st.Unlock()
	}

	last := atomic.LoadInt64(pLast)
	if now-last < _degradeInterval {
		return false
	}
	return atomic.CompareAndSwapInt64(pLast, last, now)
}
