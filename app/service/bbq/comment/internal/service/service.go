package service

import (
	"context"
	"fmt"

	pb "go-common/app/service/bbq/comment/api"
	"go-common/app/service/bbq/comment/internal/dao"
	"go-common/library/conf/paladin"

	"github.com/golang/protobuf/ptypes/empty"
)

// Service service.
type Service struct {
	ac  *paladin.Map
	dao *dao.Dao
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}
	return s
}

// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("hello %s", req.Name)
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
