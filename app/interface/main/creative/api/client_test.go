package v1

import (
	"context"
	"testing"
	"time"

	"go-common/library/log"
	"go-common/library/naming/discovery"
	"go-common/library/net/netutil/breaker"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/resolver"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func testInit() CreativeClient {
	log.Init(nil)
	conf := &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
	}
	wc := warden.NewClient(conf)
	resolver.Register(discovery.New(nil))
	conn, err := wc.Dial(context.TODO(), "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	return NewCreativeClient(conn)
}

//var client CreativeClient
//
//func init() {
//	var err error
//	client, err = NewClient(nil)
//	if err != nil {
//		panic(err)
//	}
//}

func TestFlowJudge(t *testing.T) {
	client := testInit()
	convey.Convey("TestFlowJudge", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			oids, err := client.FlowJudge(context.TODO(), &FlowRequest{Business: int64(4), Gid: int64(24), Oids: []int64{22, 333, 10110208, 10110119}})
			ctx.So(err, convey.ShouldBeNil)
			ctx.Printf("%+v\n", oids.Oids)
		})
		//ctx.Convey("When error", func(ctx convey.C) {
		//})
	})
}

func TestCheckTaskState(t *testing.T) {
	client := testInit()
	client.CheckTaskState(context.TODO(), &TaskRequest{Mid: int64(1), TaskId: int64(1)})
}
