package spy

import (
	"context"

	"go-common/app/service/main/spy/model"
	"go-common/library/net/rpc"
)

const (
	_addEvent        = "RPC.HandleEvent"
	_userScore       = "RPC.UserScore"
	_reBuildPortrait = "RPC.ReBuildPortrait"
	//_resetManually    = "RPC.ResetManually"
	_updateBaseScore  = "RPC.UpdateBaseScore"
	_refreshBaseScore = "RPC.RefreshBaseScore"
	_updateEventScore = "RPC.UpdateEventScore"
	_userInfo         = "RPC.UserInfo"
	_clearReliveTimes = "RPC.ClearReliveTimes"
	_statByID         = "RPC.StatByID"
)

const (
	_appid = "account.service.spy"
)

var (
	_noRes = &struct{}{}
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

// New create instance of service and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// HandleEvent add spy event to user.
func (s *Service) HandleEvent(c context.Context, arg *model.ArgHandleEvent) (err error) {
	err = s.client.Call(c, _addEvent, arg, _noRes)
	return
}

// UserScore archive get user spy score.
func (s *Service) UserScore(c context.Context, arg *model.ArgUserScore) (res *model.UserScore, err error) {
	res = &model.UserScore{}
	err = s.client.Call(c, _userScore, arg, res)
	return
}

// ReBuildPortrait rebuild user risk portrait by task.
func (s *Service) ReBuildPortrait(c context.Context, arg *model.ArgReBuild) (err error) {
	err = s.client.Call(c, _reBuildPortrait, arg, _noRes)
	return
}

// UpdateBaseScore cli.
func (s *Service) UpdateBaseScore(c context.Context, arg *model.ArgReset) (err error) {
	err = s.client.Call(c, _updateBaseScore, arg, _noRes)
	return
}

// RefreshBaseScore cli.
func (s *Service) RefreshBaseScore(c context.Context, arg *model.ArgReset) (err error) {
	err = s.client.Call(c, _refreshBaseScore, arg, _noRes)
	return
}

// UpdateEventScore cli.
func (s *Service) UpdateEventScore(c context.Context, arg *model.ArgReset) (err error) {
	err = s.client.Call(c, _updateEventScore, arg, _noRes)
	return
}

// UserInfo cli.
func (s *Service) UserInfo(c context.Context, arg *model.ArgUser) (res *model.UserInfo, err error) {
	err = s.client.Call(c, _userInfo, arg, res)
	return
}

// ClearReliveTimes cli.
func (s *Service) ClearReliveTimes(c context.Context, arg *model.ArgReset) (err error) {
	err = s.client.Call(c, _clearReliveTimes, arg, _noRes)
	return
}

// StatByID cli.
func (s *Service) StatByID(c context.Context, arg *model.ArgStat) (stat []*model.Statistics, err error) {
	err = s.client.Call(c, _statByID, arg, &stat)
	return
}
