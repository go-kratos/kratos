package interceptor

import (
	ctx "context"
	"errors"
	"sync"
	"testing"

	"go-common/library/ecode"
	"go-common/library/net/rpc/context"

	"golang.org/x/time/rate"
)

var (
	once sync.Once
	i    *Interceptor
	c    context.Context
)

func interceptor() {
	i = NewInterceptor("test token")
	c = context.NewContext(ctx.TODO(), "testMethod", "test user", 0)
}

func TestRate(t *testing.T) {
	once.Do(interceptor)
	if err := i.Rate(c); err != nil {
		t.Errorf("TestRate error(%v)", err)
		t.FailNow()
	}
	i.rateLimits["testMethod"] = rate.NewLimiter(1, 0)
	if err := i.Rate(c); err != ecode.Degrade {
		t.Errorf("TestRate error(%v)", err)
		t.FailNow()
	}
}

func TestStat(t *testing.T) {
	once.Do(interceptor)
	i.Stat(c, nil, errors.New("test error"))
}

func TestAuth(t *testing.T) {
	once.Do(interceptor)
	if err := i.Auth(c, nil, "test token"); err != nil {
		t.Errorf("TestAuth error(%v)", err)
		t.FailNow()
	}
	if err := i.Auth(c, nil, "token"); err != ecode.RPCNoAuth {
		t.Errorf("TestAuth error(%v)", err)
		t.FailNow()
	}
}
