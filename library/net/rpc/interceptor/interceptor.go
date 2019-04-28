package interceptor

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/rpc/context"
	"go-common/library/stat"

	"golang.org/x/time/rate"
)

var stats = stat.RPCServer

const (
	_innerService = "inner"
)

// Interceptor is rpc interceptor
type Interceptor struct {
	// for limit rate
	rateLimits map[string]*rate.Limiter
	// for auth
	token string
}

// NewInterceptor new a interceptor
func NewInterceptor(token string) *Interceptor {
	in := &Interceptor{token: token}
	in.rateLimits = make(map[string]*rate.Limiter)
	return in
}

// Rate check the call is limit or not
func (i *Interceptor) Rate(c context.Context) error {
	limit, ok := i.rateLimits[c.ServiceMethod()]
	if ok && !limit.Allow() {
		return ecode.Degrade
	}
	return nil
}

// Stat add stat info to ops
func (i *Interceptor) Stat(c context.Context, args interface{}, err error) {
	const noUser = "no_user"
	var (
		user   = c.User()
		method = c.ServiceMethod()
		tmsub  = time.Since(c.Now())
	)
	if user == "" {
		user = noUser
	}
	stats.Timing(user, int64(tmsub/time.Millisecond), method)
	stats.Incr(user, method, strconv.Itoa((ecode.Cause(err).Code())))
	if err != nil {
		log.Errorv(c,
			log.KV("args", fmt.Sprintf("%v", args)),
			log.KV("method", method),
			log.KV("duration", tmsub),
			log.KV("error", fmt.Sprintf("%+v", err)),
		)
	} else {
		if !strings.HasPrefix(method, _innerService) || bool(log.V(2)) {
			log.Infov(c,
				log.KV("args", fmt.Sprintf("%v", args)),
				log.KV("method", method),
				log.KV("duration", tmsub),
				log.KV("error", fmt.Sprintf("%+v", err)),
			)
		}
	}
}

// Auth check token has auth
func (i *Interceptor) Auth(c context.Context, addr net.Addr, token string) error {
	if i.token != token {
		return ecode.RPCNoAuth
	}
	return nil
}
