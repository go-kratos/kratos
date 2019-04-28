package gorpc

import (
	"context"

	"go-common/app/service/main/sms/model"
	"go-common/library/net/rpc"
)

const (
	_appid = "sms.service"

	_send      = "RPC.send"
	_sendBatch = "RPC.sendBatch"
)

var (
	_noRes = &struct{}{}
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

// New new service instance and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Send send sms
func (s *Service) Send(c context.Context, arg *model.ArgSend) (err error) {
	err = s.client.Call(c, _send, arg, _noRes)
	return
}

// SendBatch batch send sms
func (s *Service) SendBatch(c context.Context, arg *model.ArgSendBatch) (err error) {
	err = s.client.Call(c, _sendBatch, arg, _noRes)
	return
}
