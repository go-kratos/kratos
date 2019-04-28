package livedemo

import (
	"context"

	pb "go-common/app/interface/live/live-demo/api/http"
	"go-common/app/interface/live/live-demo/conf"
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
func (s *FooService) UnameByUid(ctx context.Context, req *pb.Bar1Req) (resp *pb.Bar1Resp, err error) {
	resp = &pb.Bar1Resp{}
	return
}

// GetInfo implementation
// 获取房间信息
// `midware:"guest"`
func (s *FooService) GetInfo(ctx context.Context, req *pb.GetInfoReq) (resp *pb.GetInfoResp, err error) {
	resp = &pb.GetInfoResp{}
	return
}

// UnameByUid3 implementation
// 根据uid得到uname v3
func (s *FooService) UnameByUid3(ctx context.Context, req *pb.Bar1Req) (resp *pb.Bar1Resp, err error) {
	resp = &pb.Bar1Resp{}
	return
}

// UnameByUid4 implementation
// test comment
// `internal:"true"`
func (s *FooService) UnameByUid4(ctx context.Context, req *pb.Bar1Req) (resp *pb.Bar1Resp, err error) {
	resp = &pb.Bar1Resp{}
	return
}

// GetDynamic implementation
// `dynamic_resp:"true"`
func (s *FooService) GetDynamic(ctx context.Context, req *pb.Bar1Req) (resp interface{}, err error) {
	resp = &pb.Bar1Resp{}
	return
}

// Nointerface implementation
// `dynamic:"true"`
func (s *FooService) Nointerface(ctx context.Context, req *pb.Bar1Req) (resp *pb.Bar1Resp, err error) {
	resp = &pb.Bar1Resp{}
	return
}
