package v2

import (
	"context"

	v2pb "go-common/app/interface/live/live-demo/api/http/v2"
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
func (s *Foo2Service) Hello(ctx context.Context, req *v2pb.Bar1Req) (resp *v2pb.Bar1Resp, err error) {
	resp = &v2pb.Bar1Resp{}
	return
}
