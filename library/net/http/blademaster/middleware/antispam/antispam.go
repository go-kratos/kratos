package antispam

import (
	"fmt"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_prefixRate  = "r_%d_%s_%d"
	_prefixTotal = "t_%d_%s_%d"
	// antispam
	_defSecond = 1
	_defHour   = 1
)

// Antispam is a antispam instance.
type Antispam struct {
	redis *redis.Pool
	conf  *Config
}

// Config antispam config.
type Config struct {
	On     bool // switch on/off
	Second int  // every N second allow N requests.
	N      int  // one unit allow N requests.
	Hour   int  // every N hour allow M requests.
	M      int  // one winodw allow M requests.

	Redis *redis.Config
}

func (c *Config) validate() error {
	if c == nil {
		return errors.New("antispam: empty config")
	}
	if c.Second < _defSecond {
		return errors.New("antispam: invalid Second")
	}
	if c.Hour < _defHour {
		return errors.New("antispam: invalid Hour")
	}
	return nil
}

// New new a antispam service.
func New(c *Config) (s *Antispam) {
	if err := c.validate(); err != nil {
		panic(err)
	}
	s = &Antispam{
		redis: redis.NewPool(c.Redis),
	}
	s.Reload(c)
	return s
}

// Reload reload antispam config.
func (s *Antispam) Reload(c *Config) {
	if err := c.validate(); err != nil {
		log.Error("Failed to reload antispam: %+v", err)
		return
	}
	s.conf = c
}

// Rate antispam by user + path.
func (s *Antispam) Rate(c *bm.Context, second, count int) (err error) {
	mid, ok := c.Get("mid")
	if !ok {
		return
	}
	curSecond := int(time.Now().Unix())
	burst := curSecond - curSecond%second
	key := rateKey(mid.(int64), c.Request.URL.Path, burst)
	return s.antispam(c, key, second, count)
}

// Total antispam by user + path.
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

func (s *Antispam) antispam(c *bm.Context, key string, interval, count int) error {
	conn := s.redis.Get(c)
	defer conn.Close()
	incred, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return nil
	}
	if incred == 1 {
		conn.Do("EXPIRE", key, interval)
	}
	if incred > int64(count) {
		return ecode.LimitExceed
	}
	return nil
}

func rateKey(mid int64, path string, burst int) string {
	return fmt.Sprintf(_prefixRate, mid, path, burst)
}

func totalKey(mid int64, path string, burst int) string {
	return fmt.Sprintf(_prefixTotal, mid, path, burst)
}

func (s *Antispam) ServeHTTP(ctx *bm.Context) {
	if err := s.Rate(ctx, s.conf.Second, s.conf.N); err != nil {
		ctx.JSON(nil, ecode.LimitExceed)
		ctx.Abort()
		return
	}
	if err := s.Total(ctx, s.conf.Hour, s.conf.M); err != nil {
		ctx.JSON(nil, ecode.LimitExceed)
		ctx.Abort()
		return
	}
}

// Handler is antispam handle.
func (s *Antispam) Handler() bm.HandlerFunc {
	return s.ServeHTTP
}
