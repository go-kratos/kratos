package antispam

import (
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/main/upload/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

const (
	_prefixRate  = "r_%d_%s_%d"
	_prefixTotal = "t_%d_%s_%d"
)

// Antispam is a antispam instance.
type Antispam struct {
	redis     *redis.Pool
	limitFunc func(bucket, dir string) (model.DirRateConfig, bool)
	conf      *Config
}

// Config antispam config.
// On     bool // switch on/off
// Second int  // every N second allow N requests.
// N      int  // one unit allow N requests.
// Hour   int  // every N hour allow M requests.
// M      int  // one winodw allow M requests.
type Config struct {
	On     bool // switch on/off
	Second int  // every N second allow N requests.
	N      int  // one unit allow N requests.
	Hour   int  // every N hour allow M requests.
	M      int  // one winodw allow M requests.

	Redis *redis.Config
}

// New new a antispam service.
func New(c *Config, l func(bucket, dir string) (model.DirRateConfig, bool)) (s *Antispam) {
	if c == nil {
		panic("antispam config nil")
	}
	s = &Antispam{
		limitFunc: l,
		redis:     redis.NewPool(c.Redis),
	}
	s.conf = c
	return s
}

// NativeRate limit user + path second level
func (s *Antispam) NativeRate(c *bm.Context, path string, mid interface{}) (err error) {
	curSecond := int(time.Now().Unix())
	burst := curSecond - curSecond%s.conf.Second
	key := rateKey(mid.(int64), path, burst)
	return s.antispam(c, key, s.conf.Second, s.conf.N)
}

// Rate antispam by user + bucket + dir.
func (s *Antispam) Rate(c *bm.Context) (err error) {
	mid, ok := c.Get("mid")
	if !ok {
		return
	}
	ap := new(struct {
		Bucket string `form:"bucket" json:"bucket"`
		Dir    string `form:"dir" json:"dir"`
	})
	if err = c.BindWith(ap, binding.FormMultipart); err != nil {
		return s.NativeRate(c, c.Request.URL.Path, mid)
	}
	if ap.Bucket == "" || ap.Dir == "" { //not need dir limit
		return s.NativeRate(c, c.Request.URL.Path, mid)
	}
	limit, ok := s.limitFunc(ap.Bucket, ap.Dir)
	if !ok {
		return s.NativeRate(c, c.Request.URL.Path, mid)
	}
	if limit.SecondQPS == 0 || limit.CountQPS == 0 {
		return s.NativeRate(c, c.Request.URL.Path, mid)
	}
	path := strings.Join([]string{ap.Bucket, ap.Dir}, "_")
	curSecond := int(time.Now().Unix())
	burst := curSecond - curSecond%limit.SecondQPS
	key := rateKey(mid.(int64), path, burst)
	return s.antispam(c, key, limit.SecondQPS, limit.CountQPS)
}

func totalKey(mid int64, path string, burst int) string {
	return fmt.Sprintf(_prefixTotal, mid, path, burst)
}

// Total antispam by user + path hour level
func (s *Antispam) Total(c *bm.Context, hour, count int) (err error) {
	second := hour * 3600
	mid, ok := c.Get("mid")
	if !ok {
		return
	}
	curHour := int(time.Now().Unix() / 3600)
	burst := curHour - curHour%hour
	key := totalKey(mid.(int64), c.Request.URL.Path, burst)
	return s.antispam(c, key, second, count)
}

func (s *Antispam) antispam(c *bm.Context, key string, interval, count int) (err error) {
	conn := s.redis.Get(c)
	defer conn.Close()
	cur, err := redis.Int(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		err = nil
		return
	}
	if cur >= count {
		err = ecode.LimitExceed
		return
	}
	err = nil
	conn.Send("INCR", key)
	conn.Send("EXPIRE", key, interval)
	if err1 := conn.Flush(); err1 != nil {
		return
	}
	for i := 0; i < 2; i++ {
		if _, err1 := conn.Receive(); err1 != nil {
			return
		}
	}
	return
}

func rateKey(mid int64, path string, burst int) string {
	return fmt.Sprintf(_prefixRate, mid, path, burst)
}

func (s *Antispam) ServeHTTP(ctx *bm.Context) {
	// user + bucket + dir.
	if err := s.Rate(ctx); err != nil {
		ctx.JSON(nil, ecode.ServiceUnavailable)
		ctx.Abort()
		return
	}
	// user + path
	if err := s.Total(ctx, s.conf.Hour, s.conf.M); err != nil {
		ctx.JSON(nil, ecode.ServiceUnavailable)
		ctx.Abort()
		return
	}
}

// Handler is antispam handle.
func (s *Antispam) Handler() bm.HandlerFunc {
	return s.ServeHTTP
}
