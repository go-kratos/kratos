package livedemo

import (
	"context"

	pb "go-common/app/interface/live/live-demo/api/http"
	"go-common/app/interface/live/live-demo/conf"
)

// Foo2Service struct
type Foo2Service struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewFoo2Service init
func NewFoo2Service(c *conf.Config) (s *Foo2Service) {
	s = &Foo2Service{
		conf: c,
	}
	return s
}

// Hello implementation
func (s *Foo2Service) Hello(ctx context.Context, req *pb.Bar1Req) (resp *pb.Bar1Resp, err error) {
	resp = &pb.Bar1Resp{}
	return
}
