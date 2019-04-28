package v1

import (
	"context"

	v1pb "go-common/app/service/live/grpc-demo/api/grpc/v1"
	"go-common/app/service/live/grpc-demo/conf"

	"google.golang.org/grpc/metadata"
)

// GreeterService struct
type GreeterService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewGreeterService init
func NewGreeterService(c *conf.Config) (s *GreeterService) {
	s = &GreeterService{
		conf: c,
	}
	return s
}

// SayHello implementation
func (s *GreeterService) SayHello(ctx context.Context, req *v1pb.GeeterReq) (resp *v1pb.GreeterResp, err error) {
	resp = &v1pb.GreeterResp{}
	metadata.FromIncomingContext(ctx)
	resp.Uid = req.Uid
	return
}

// SayHelloInternal implementation
// `method:"POST" internal:"true"`
func (s *GreeterService) SayHelloInternal(ctx context.Context, req *v1pb.GeeterReq) (resp *v1pb.GreeterResp, err error) {
	resp = &v1pb.GreeterResp{}
	return
}
