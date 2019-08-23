package blademaster

import (
	"fmt"
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
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

		var mid int64
		if v, ok := c.Get("mid"); ok {
			mid, _ = v.(int64)
		}
		err := c.Error
		cerr := ecode.Cause(err)
		dt := time.Since(now)
		caller := metadata.String(c, metadata.Caller)
		if caller == "" {
			caller = noUser
		}

		buvid := ""
		if dev, ok := c.Get("device"); ok {
			device, ok := dev.(*Device)
			if ok {
				buvid = device.Buvid
			}
		}

		if c.RoutePath != "" {
			MetricServerReqCodeTotal.Inc(c.RoutePath[1:], caller, strconv.FormatInt(int64(cerr.Code()), 10))
			MetricServerReqDur.Observe(int64(dt/time.Millisecond), c.RoutePath[1:], caller)
		}

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
			log.KVInt64("mid", mid),
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
			log.KVString("buvid", buvid),
			log.KVString("source", "http-access-log"),
		)
	}
}
