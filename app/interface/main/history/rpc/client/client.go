package client

import (
	"context"

	"go-common/app/interface/main/history/model"
	"go-common/library/net/rpc"
)

const (
	_progress      = "RPC.Progress"
	_position      = "RPC.Position"
	_add           = "RPC.Add"
	_delete        = "RPC.Delete"
	_history       = "RPC.History"
	_historyCursor = "RPC.HistoryCursor"
	_clear         = "RPC.Clear"
)

var (
	_noRes = &struct{}{}
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

const (
	_appid = "community.service.history"
)

// New create instance of service and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Progress return map[mid]*history.
func (s *Service) Progress(c context.Context, arg *model.ArgPro) (res map[int64]*model.History, err error) {
	res = make(map[int64]*model.History)
	err = s.client.Call(c, _progress, arg, &res)
	return
}

// Position return map[mid]*history.
func (s *Service) Position(c context.Context, arg *model.ArgPos) (res *model.History, err error) {
	res = &model.History{}
	err = s.client.Call(c, _position, arg, res)
	return
}

// Add add history .
func (s *Service) Add(c context.Context, arg *model.ArgHistory) (err error) {
	err = s.client.Call(c, _add, arg, _noRes)
	return
}

// Delete add history .
func (s *Service) Delete(c context.Context, arg *model.ArgDelete) (err error) {
	err = s.client.Call(c, _delete, arg, _noRes)
	return
}

// History return all histories .
func (s *Service) History(c context.Context, arg *model.ArgHistories) (res []*model.Resource, err error) {
	err = s.client.Call(c, _history, arg, &res)
	return
}

// HistoryCursor return all histories .
func (s *Service) HistoryCursor(c context.Context, arg *model.ArgCursor) (res []*model.Resource, err error) {
	err = s.client.Call(c, _historyCursor, arg, &res)
	return
}

// Clear clear history
func (s *Service) Clear(c context.Context, arg *model.ArgClear) (err error) {
	err = s.client.Call(c, _clear, arg, _noRes)
	return
}
