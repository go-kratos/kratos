package rate

import (
	"net/http"
	"sync/atomic"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"golang.org/x/time/rate"
)

const (
	_defBurst = 100
)

// Limiter controls how frequently events are allowed to happen.
type Limiter struct {
	apps atomic.Value
	urls atomic.Value
}

// Config limitter  conf.
type Config struct {
	Apps map[string]*Limit
	URLs map[string]*Limit
}

// Limit limit conf.
type Limit struct {
	Limit rate.Limit
	Burst int
}

// New return Limiter.
func New(conf *Config) (l *Limiter) {
	l = &Limiter{}
	l.apps.Store(make(map[string]*rate.Limiter))
	l.urls.Store(make(map[string]*rate.Limiter))
	if conf != nil {
		l.Reload(conf)
	}
	return
}

// Reload reload limit conf.
func (l *Limiter) Reload(c *Config) {
	if c == nil {
		return
	}
	var (
		ok  bool
		al  *rate.Limiter
		ul  *rate.Limiter
		as  map[string]*rate.Limiter
		nas map[string]*rate.Limiter
		us  map[string]*rate.Limiter
		nus map[string]*rate.Limiter
	)
	if as, ok = l.apps.Load().(map[string]*rate.Limiter); !ok {
		log.Error("apps limiter load map hava no data ")
		return
	}
	nas = make(map[string]*rate.Limiter, len(as))
	for k, v := range as {
		nas[k] = v
	}
	for k, v := range c.Apps {
		if al, ok = nas[k]; !ok || (al.Burst() != v.Burst || al.Limit() != v.Limit) {
			nas[k] = rate.NewLimiter(v.fix())
		}
	}
	l.apps.Store(nas)

	if us, ok = l.urls.Load().(map[string]*rate.Limiter); !ok {
		log.Error("urls limiter load map hava no data ")
		return
	}
	nus = make(map[string]*rate.Limiter, len(us))
	for k, v := range us {
		nus[k] = v
	}
	for k, v := range c.URLs {
		if ul, ok = nus[k]; !ok || (ul.Burst() != v.Burst || ul.Limit() != v.Limit) {
			nus[k] = rate.NewLimiter(v.fix())
		}
	}
	l.urls.Store(nus)
}

func (l *Limit) fix() (lim rate.Limit, b int) {
	lim = rate.Inf
	b = _defBurst
	if l.Limit <= 0 {
		lim = rate.Inf
	} else {
		lim = l.Limit
	}
	if l.Burst > 0 {
		b = l.Burst
	}
	return
}

// Allow reports whether event may happen at time now.
func (l *Limiter) Allow(appKey, path string) bool {
	if as, ok := l.apps.Load().(map[string]*rate.Limiter); ok {
		if lim, ok := as[appKey]; ok {
			if !lim.Allow() {
				return false
			}
		}
	}
	if us, ok := l.urls.Load().(map[string]*rate.Limiter); ok {
		if lim, ok := us[path]; ok {
			if !lim.Allow() {
				return false
			}
		}
	}
	return true
}

func (l *Limiter) ServeHTTP(c *bm.Context) {
	req := c.Request
	appkey := req.Form.Get("appkey")
	path := req.URL.Path
	if !l.Allow(appkey, path) {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
}

// Handler is router allow handle.
func (l *Limiter) Handler() bm.HandlerFunc {
	return l.ServeHTTP
}
