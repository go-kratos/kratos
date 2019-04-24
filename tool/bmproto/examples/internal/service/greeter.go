package service

import (
	pb "github.com/bilibili/kratos/tool/bmproto/examples/api"
	"context"
)

// GreeterService struct
type GreeterService struct {
	//	conf   *conf.Config
	// optionally add other properties here, such as dao
	// 值得注意的是，多个service公用一个Dao时，可以在外面New一个Dao对象传进来
	// 不必每个service都New一个Dao，以免造成资源浪费(例如mysql连接等)
	// dao *dao.Dao
}

//NewGreeterService init
func NewGreeterService(
// c *conf.Config
) (s *GreeterService) {
	s = &GreeterService{
		//		conf:   c,
	}
	return s
}

// SayHello implementation
// api 标题
// api 说明
func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {
	resp = &pb.HelloResponse{}
	return
}

// SayHelloCustomUrl implementation
func (s *GreeterService) SayHelloCustomUrl(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {
	resp = &pb.HelloResponse{}
	return
}

// SayHelloPost implementation
func (s *GreeterService) SayHelloPost(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {
	resp = &pb.HelloResponse{}
	return
}
