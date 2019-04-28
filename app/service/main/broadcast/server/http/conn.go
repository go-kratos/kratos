package http

import (
	"io/ioutil"

	pb "go-common/app/service/main/broadcast/api/grpc/v1"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/gogo/protobuf/proto"
)

func connect(ctx *bm.Context) {
	query := ctx.Request.URL.Query()
	if token, err := c.Get("httpToken").String(); err != nil || token != query.Get("token") {
		ctx.Protobuf(nil, ecode.Unauthorized)
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	var req pb.ConnectReq
	if err = proto.Unmarshal(b, &req); err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	mid, key, room, platform, accepts, err := srv.Connect(ctx, req.Server, req.ServerKey, req.Cookie, req.Token)
	if err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	ctx.Protobuf(&pb.ConnectReply{Mid: mid, Key: key, RoomID: room, Accepts: accepts, Platform: platform}, nil)
}

func disconnect(ctx *bm.Context) {
	query := ctx.Request.URL.Query()
	if token, err := c.Get("httpToken").String(); err != nil || token != query.Get("token") {
		ctx.Protobuf(nil, ecode.Unauthorized)
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	var req pb.DisconnectReq
	if err = proto.Unmarshal(b, &req); err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	has, err := srv.Disconnect(ctx, req.Mid, req.Key, req.Server)
	if err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	ctx.Protobuf(&pb.DisconnectReply{Has: has}, nil)
}

func heartbeat(ctx *bm.Context) {
	query := ctx.Request.URL.Query()
	if token, err := c.Get("httpToken").String(); err != nil || token != query.Get("token") {
		ctx.Protobuf(nil, ecode.Unauthorized)
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	var req pb.HeartbeatReq
	if err = proto.Unmarshal(b, &req); err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	if err := srv.Heartbeat(ctx, req.Mid, req.Key, req.Server); err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	ctx.Protobuf(&pb.HeartbeatReply{}, nil)
}

func renewOnline(ctx *bm.Context) {
	query := ctx.Request.URL.Query()
	if token, err := c.Get("httpToken").String(); err != nil || token != query.Get("token") {
		ctx.Protobuf(nil, ecode.Unauthorized)
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	var req pb.OnlineReq
	if err = proto.Unmarshal(b, &req); err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	roomCount, err := srv.RenewOnline(ctx, req.Server, req.Sharding, req.RoomCount)
	if err != nil {
		ctx.Protobuf(nil, err)
		return
	}
	ctx.Protobuf(&pb.OnlineReply{RoomCount: roomCount}, nil)
}
