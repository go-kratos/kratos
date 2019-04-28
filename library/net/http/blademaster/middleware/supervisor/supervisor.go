package supervisor

import (
	"time"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// Config supervisor conf.
type Config struct {
	On    bool      // all post/put/delete method off.
	Begin time.Time // begin time
	End   time.Time // end time
}

// Supervisor supervisor midleware.
type Supervisor struct {
	conf *Config
	on   bool
}

// New new and return supervisor midleware.
func New(c *Config) (s *Supervisor) {
	s = &Supervisor{
		conf: c,
	}
	s.Reload(c)
	return
}

// Reload reload supervisor conf.
func (s *Supervisor) Reload(c *Config) {
	if c == nil {
		return
	}
	s.on = c.On && c.Begin.Before(c.End)
	s.conf = c // NOTE datarace but no side effect.
}

func (s *Supervisor) ServeHTTP(c *bm.Context) {
	if s.on {
		now := time.Now()
		method := c.Request.Method
		if s.forbid(method, now) {
			c.JSON(nil, ecode.ServiceUpdate)
			c.Abort()
			return
		}
	}
}

func (s *Supervisor) forbid(method string, now time.Time) bool {
	// only allow GET request.
	return method != "GET" && now.Before(s.conf.End) && now.After(s.conf.Begin)
}
