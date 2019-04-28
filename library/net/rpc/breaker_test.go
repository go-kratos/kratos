package rpc

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	rcontext "go-common/library/net/rpc/context"
)

func TestBreaker(t *testing.T) {
	env.DeployEnv = "dev"
	env.Zone = "testzone"
	log.Init(&log.Config{Stdout: false})
	d := discovery.New(nil)
	cancel, err := d.Register(context.TODO(), &naming.Instance{
		Zone:     "testzone",
		Env:      "dev",
		AppID:    "test.appid",
		Hostname: "test.host",
		Addrs:    []string{"gorpc://127.0.0.1:9000"},
		Version:  "1",
		LastTs:   time.Now().UnixNano(),
		Metadata: map[string]string{"weight": "100"},
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		cancel()
		time.Sleep(time.Second)
	}()
	svr := NewServer(&ServerConfig{Proto: "tcp", Addr: "127.0.0.1:9000"})
	err = svr.Register(&BreakerRPC{})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 200)
	RunCli()
}

type BreakerReq struct {
	Name string
}

type BreakerReply struct {
	Success bool
}

type BreakerRPC struct{}

// Ping check connection success.
func (r *BreakerRPC) Ping(c rcontext.Context, arg *BreakerReq, res *BreakerReply) (err error) {
	if rand.Int31n(100) < 40 {
		return ecode.ServerErr
	}
	res.Success = true
	return
}

func RunCli() {
	var success int64
	var su int64
	var se int64
	var other int64
	cli := NewDiscoveryCli("test.appid", nil)
	for i := 0; i < 1000; i++ {
		var res BreakerReply
		err := cli.Call(
			context.Background(),
			"BreakerRPC.Ping",
			&BreakerReq{Name: "test"},
			&res,
		)
		if err == nil || ecode.OK.Equal(err) {
			atomic.AddInt64(&success, 1)
		} else if ecode.ServiceUnavailable.Equal(err) {
			atomic.AddInt64(&su, 1)
		} else if ecode.ServerErr.Equal(err) {
			atomic.AddInt64(&se, 1)
		} else {
			atomic.AddInt64(&other, 1)
		}
		time.Sleep(time.Millisecond * 9)
	}
	fmt.Printf("success:%d su:%d se:%d other:%d \n", success, su, se, other)
}
