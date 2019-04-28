// to avoid recycle imports,
// the test must be in different package with liverpc...
// otherwise, test import => generated pb
// 			 generated pb => import liverpc (which includes the test)
package liverpc

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"go-common/library/net/rpc/liverpc"
	"go-common/library/net/rpc/liverpc/testdata"
	v1 "go-common/library/net/rpc/liverpc/testdata/v1"
	v2 "go-common/library/net/rpc/liverpc/testdata/v2"

	"github.com/pkg/errors"
)

func TestDialClient(t *testing.T) {
	cli := testdata.New(nil)
	var req = &v1.RoomGetInfoReq{Id: 1002}
	var hd = &liverpc.Header{
		Platform:    "ios",
		Src:         "test",
		Buvid:       "AUTO3315341311353015",
		TraceId:     "18abb1a2596c43ea:18abb1a2596c43ea:0:0",
		Uid:         10,
		Caller:      "live-api.rpc",
		UserIp:      "127.0.0.1",
		SourceGroup: "default",
	}
	var ctx = context.TODO()
	reply, err := cli.V1Room.GetInfo(ctx, req, &liverpc.HeaderOption{Header: hd})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("reply:%v %v %v", reply.Code, reply.Msg, reply.Data)

	_, err = cli.GetRawCli().CallRaw(context.TODO(), 2, "Room.get_by_ids", &liverpc.Args{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestCallRaw(t *testing.T) {
	var cli = liverpc.NewClient(&liverpc.ClientConfig{AppID: "live.room"})
	var hd = &liverpc.Header{
		Platform:    "ios",
		Src:         "test",
		Buvid:       "AUTO3315341311353015",
		TraceId:     "18abb1a2596c43ea:18abb1a2596c43ea:0:0",
		Uid:         10,
		Caller:      "live-api.rpc",
		UserIp:      "127.0.0.1",
		SourceGroup: "default",
	}
	var req = &liverpc.Args{Body: map[string]interface{}{"id": 1002}, Header: hd}
	hd = nil
	reply, err := cli.CallRaw(context.TODO(), 1, "Room.get_info", req, &liverpc.TimeoutOption{Timeout: 200 * time.Millisecond})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("reply:%v %v %v", reply.Code, reply.Message, string(reply.Data))

	_, err = cli.CallRaw(context.TODO(), 1, "Room.get_info", req, &liverpc.TimeoutOption{Timeout: 2 * time.Millisecond})
	if err == nil {
		t.Error(errors.New("should fail"))
		t.FailNow()
	}
}

func TestMap(t *testing.T) {
	client := liverpc.NewClient(&liverpc.ClientConfig{AppID: "live.room"})
	var rpcClient = v2.NewRoomRPCClient(client)
	var req = &v2.RoomGetByIdsReq{Ids: []int64{1002}}
	var header = &liverpc.Header{
		Platform:    "ios",
		Src:         "test",
		Buvid:       "AUTO3315341311353015",
		TraceId:     "18abb1a2596c43ea:18abb1a2596c43ea:0:0",
		Uid:         10,
		Caller:      "live-api.rpc",
		UserIp:      "127.0.0.1",
		SourceGroup: "default",
	}
	var ctx = context.TODO()
	reply, err := rpcClient.GetByIds(ctx, req, liverpc.HeaderOption{Header: header})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("reply:%v %v %v", reply.Code, reply.Msg, reply.Data)
}

func TestDiscoveryClient(t *testing.T) {
	conf := &liverpc.ClientConfig{
		AppID: "live.room",
	}
	cli := liverpc.NewClient(conf)
	arg := &v1.RoomGetInfoReq{Id: 1001}
	var rpcClient = v1.NewRoomRPCClient(cli)
	reply, err := rpcClient.GetInfo(context.TODO(), arg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("reply:%+v", reply)
}

func TestCancel(t *testing.T) {
	conf := &liverpc.ClientConfig{
		AppID: "live.room",
	}
	cli := liverpc.NewClient(conf)
	arg := &v1.RoomGetInfoReq{Id: 1001}
	var rpcClient = v1.NewRoomRPCClient(cli)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond)
	defer cancel()
	var err error
	_, err = rpcClient.GetInfo(ctx, arg)
	if err == nil || errors.Cause(err) != context.DeadlineExceeded {
		t.Error(err)
		t.FailNow()
	}
}

func TestCallRawCancel(t *testing.T) {
	var cli = liverpc.NewClient(&liverpc.ClientConfig{AppID: "live.room"})
	var hd = &liverpc.Header{
		Platform:    "ios",
		Src:         "test",
		Buvid:       "AUTO3315341311353015",
		TraceId:     "18abb1a2596c43ea:18abb1a2596c43ea:0:0",
		Uid:         10,
		Caller:      "live-api.rpc",
		UserIp:      "127.0.0.1",
		SourceGroup: "default",
	}
	var req = &liverpc.Args{Body: map[string]interface{}{"id": 1002}, Header: hd}
	hd = nil
	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond)
	defer cancel()
	_, err := cli.CallRaw(ctx, 1, "Room.get_info", req)
	if err == nil || errors.Cause(err) != context.DeadlineExceeded {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("err is +%v", err)
}

func BenchmarkDialClient(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	var header = &liverpc.Header{
		Platform:    "ios",
		Src:         "test",
		Buvid:       "AUTO3315341311353015",
		TraceId:     "18abb1a2596c43ea:18abb1a2596c43ea:0:0",
		Uid:         10,
		Caller:      "live-api.rpc",
		UserIp:      "127.0.0.1",
		SourceGroup: "default",
	}

	var ctx = context.TODO()
	for i := 0; i < b.N; i++ { //use b.N for looping
		id := rand.Intn(10000)
		arg := &v1.RoomGetInfoReq{Id: int64(id)}
		cli := liverpc.NewClient(&liverpc.ClientConfig{AppID: "live.room"})
		var rpcClient = v1.NewRoomRPCClient(cli)
		_, err := rpcClient.GetInfo(ctx, arg, liverpc.HeaderOption{Header: header})
		if err != nil {
			b.Errorf("%s %d", err, i)
			b.FailNow()
		}
	}
}
