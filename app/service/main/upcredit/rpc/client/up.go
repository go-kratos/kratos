package client

import (
	"go-common/library/net/rpc"
)

const (
	_appid = "archive.service.upcredit"
)

//Service rpc service
type Service struct {
	client *rpc.Client2
}

//RPC interface
type RPC interface {
}

//New create
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{
		client: rpc.NewDiscoveryCli(_appid, c),
	}
	return
}
