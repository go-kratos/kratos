package blademaster

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/metadata"
)

// Logger is logger  middleware
func Logger() HandlerFunc {
	const noUser = "no_user"
	return func(c *Context) {
		now := time.Now()
		ip := metadata.String(c, metadata.RemoteIP)
		req := c.Request
		path := req.URL.Path
		params := req.Form
		var quota float64
		if deadline, ok := c.Context.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		err := c.Error
		cerr := ecode.Cause(err)
		dt := time.Since(now)
		caller := metadata.String(c, metadata.Caller)
		if caller == "" {
			caller = noUser
		}

		stats.Incr(caller, path[1:], strconv.FormatInt(int64(cerr.Code()), 10))
		stats.Timing(caller, int64(dt/time.Millisecond), path[1:])

		lf := log.Infov
		errmsg := ""
		isSlow := dt >= (time.Millisecond * 500)
		if err != nil {
			errmsg = err.Error()
			lf = log.Errorv
			if cerr.Code() > 0 {
				lf = log.Warnv
			}
		} else {
			if isSlow {
				lf = log.Warnv
			}
		}
		lf(c,
			log.KVString("method", req.Method),
			log.KVString("ip", ip),
			log.KVString("user", caller),
			log.KVString("path", path),
			log.KVString("params", params.Encode()),
			log.KVInt("ret", cerr.Code()),
			log.KVString("msg", cerr.Message()),
			log.KVString("stack", fmt.Sprintf("%+v", err)),
			log.KVString("err", errmsg),
			log.KVFloat64("timeout_quota", quota),
			log.KVFloat64("ts", dt.Seconds()),
			log.KVString("source", "http-access-log"),
		)
	}
}
