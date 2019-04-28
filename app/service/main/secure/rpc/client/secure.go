package secure

import (
	"context"

	model "go-common/app/service/main/secure/model"
	"go-common/library/net/rpc"
)

const (
	_status      = "RPC.Status"
	_expect      = "RPC.ExpectionLoc"
	_addFeedback = "RPC.AddFeedBack"
	_closeNotify = "RPC.CloseNotify"
)

const (
	_appid = "account.service.secure"
)

var (
	_noRes = &struct{}{}
)

// Service rpc service.
type Service struct {
	client *rpc.Client2
}

// New new rpc service.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Status get the ip info.
func (s *Service) Status(c context.Context, arg *model.ArgSecure) (res *model.Msg, err error) {
	res = new(model.Msg)
	err = s.client.Call(c, _status, arg, &res)
	return
}

// CloseNotify clsoe notify.
func (s *Service) CloseNotify(c context.Context, arg *model.ArgSecure) (err error) {
	return s.client.Call(c, _closeNotify, arg, &_noRes)
}

// AddFeedBack  add expection feedback.
func (s *Service) AddFeedBack(c context.Context, arg *model.ArgFeedBack) (err error) {
	return s.client.Call(c, _addFeedback, arg, &_noRes)
}

// ExpectionLoc get expection loc.
func (s *Service) ExpectionLoc(c context.Context, arg *model.ArgSecure) (res []*model.Expection, err error) {
	err = s.client.Call(c, _expect, arg, &res)
	return
}
