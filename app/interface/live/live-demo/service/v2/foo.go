package v2

import (
	"context"

	v2pb "go-common/app/interface/live/live-demo/api/http/v2"
	"go-common/app/interface/live/live-demo/conf"
	"go-common/app/interface/live/live-demo/dao"
	"go-common/app/service/live/room/api/liverpc/v2"
	"go-common/library/ecode"
	"go-common/library/ecode/pb"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// FooService struct
type FooService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewFooService init
func NewFooService(c *conf.Config) (s *FooService) {
	s = &FooService{
		conf: c,
	}
	return s
}

// Foo 相关服务

// UnameByUid implementation
// 根据uid得到uname
// `method:"post" midware:"auth,verify"`
//
// 这是详细说明
func (s *FooService) UnameByUid(ctx context.Context, req *v2pb.Bar1Req) (resp *v2pb.Bar1Resp, err error) {
	resp = &v2pb.Bar1Resp{}
	return
}

// GetInfo implementation
// 获取房间信息
func (s *FooService) GetInfo(ctx context.Context, req *v2pb.GetInfoReq) (resp *v2pb.GetInfoResp, err error) {
	//msg = "hello"
	reply, err := dao.RoomAPI.V2Room.GetByIds(ctx, &v2.RoomGetByIdsReq{Ids: []int64{int64(req.RoomId)}})
	if err != nil {
		err = errors.Wrap(&pb.Error{ErrCode: 1, ErrMessage: "call room error"}, err.Error())
		//msg = "Call Room Err"
		return
	}
	log.Info("req is %v\n", req.RoomId)
	room, ok := reply.Data[req.RoomId]
	if !ok {
		err = ecode.RoomNotFound
		return
	}
	resp = &v2pb.GetInfoResp{}
	resp.Roomid = room.Roomid
	resp.Uname = room.Uname
	resp.Amap = map[int32]string{123: "world"}
	mid, _ := metadata.Value(ctx, "mid").(int64)
	resp.Mid = mid
	//resp.LiveTime = room.LiveTime
	return
}

// UnameByUid3 implementation
// 根据uid得到uname v3
func (s *FooService) UnameByUid3(ctx context.Context, req *v2pb.Bar1Req) (resp *v2pb.Bar1Resp, err error) {
	resp = &v2pb.Bar1Resp{}
	return
}

// UnameByUid4 implementation
// test comment
func (s *FooService) UnameByUid4(ctx context.Context, req *v2pb.Bar1Req) (resp *v2pb.Bar1Resp, err error) {
	resp = &v2pb.Bar1Resp{}
	return
}

// GetDynamic implementation
// `dynamic_resp:"true"`
func (s *FooService) GetDynamic(ctx context.Context, req *v2pb.Bar1Req) (resp interface{}, err error) {
	resp = &v2pb.Bar1Resp{Uname: "hehe"}
	return
}

// Nointerface implementation
// `dynamic:"true"`
func (s *FooService) Nointerface(ctx context.Context, req *v2pb.Bar1Req) (resp *v2pb.Bar1Resp, err error) {
	resp = &v2pb.Bar1Resp{}
	return
}

// JsonReq implementation
func (s *FooService) JsonReq(ctx context.Context, req *v2pb.JsonReq) (resp *v2pb.JsonResp, err error) {
	resp = &v2pb.JsonResp{P1: req.P1}
	return
}
