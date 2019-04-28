package gorpc

import (
	"go-common/library/net/rpc"
)

const (
	_appid = "archive.service"
)

var (
	_noArg = &struct{}{}
)

// Service2 is archive rpc client.
type Service2 struct {
	client *rpc.Client2
}

// New2 new a archive rpc client.
func New2(c *rpc.ClientConfig) (s *Service2) {
	s = &Service2{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}
