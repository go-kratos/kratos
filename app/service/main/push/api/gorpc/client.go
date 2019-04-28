package gorpc

import (
	"context"

	"go-common/app/service/main/push/model"
	"go-common/library/net/rpc"
)

const (
	_addReport          = "RPC.AddReport"
	_delInvalidReports  = "RPC.DelInvalidReports"
	_delReport          = "RPC.DelReport"
	_addCallback        = "RPC.AddCallback"
	_addReportCache     = "RPC.AddReportCache"
	_addUserReportCache = "RPC.AddUserReportCache"
	_setting            = "RPC.Setting"
	_setSetting         = "RPC.SetSetting"
	_addMidProgress     = "RPC.AddMidProgress"
	_addTokenCache      = "RPC.AddTokenCache"
	_addTokensCache     = "RPC.AddTokensCache"
)

var (
	// _noArg   = &struct{}{}
	_noReply = &struct{}{}
	_appid   = "push.service"
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

// AddReport adds report.
func (s *Service) AddReport(c context.Context, arg *model.ArgReport) (err error) {
	err = s.client.Call(c, _addReport, arg, _noReply)
	return
}

// DelInvalidReports deletes invalid reports.
func (s *Service) DelInvalidReports(c context.Context, arg *model.ArgDelInvalidReport) (err error) {
	err = s.client.Call(c, _delInvalidReports, arg, _noReply)
	return
}

// DelReport deletes report.
func (s *Service) DelReport(c context.Context, arg *model.ArgReport) (err error) {
	err = s.client.Call(c, _delReport, arg, _noReply)
	return
}

// AddCallback adds callback data.
func (s *Service) AddCallback(c context.Context, arg *model.ArgCallback) (err error) {
	err = s.client.Call(c, _addCallback, arg, _noReply)
	return
}

// AddReportCache adds report.
func (s *Service) AddReportCache(c context.Context, arg *model.ArgReport) (err error) {
	err = s.client.Call(c, _addReportCache, arg, _noReply)
	return
}

// AddUserReportCache adds user report cache.
func (s *Service) AddUserReportCache(c context.Context, arg *model.ArgUserReports) (err error) {
	err = s.client.Call(c, _addUserReportCache, arg, _noReply)
	return
}

// Setting gets user push switch setting.
func (s *Service) Setting(c context.Context, arg *model.ArgMid) (res map[int]int, err error) {
	err = s.client.Call(c, _setting, arg, &res)
	return
}

// SetSetting sets user push switch setting.
func (s *Service) SetSetting(c context.Context, arg *model.ArgSetting) (err error) {
	err = s.client.Call(c, _setSetting, arg, _noReply)
	return
}

// AddMidProgress adds mid count number to task's progress field
func (s *Service) AddMidProgress(c context.Context, arg *model.ArgMidProgress) (err error) {
	err = s.client.Call(c, _addMidProgress, arg, _noReply)
	return
}

// AddTokenCache add token cache
func (s *Service) AddTokenCache(ctx context.Context, arg *model.ArgReport) (err error) {
	err = s.client.Call(ctx, _addTokenCache, arg, _noReply)
	return
}

// AddTokensCache add tokens cache
func (s *Service) AddTokensCache(ctx context.Context, arg *model.ArgReports) (err error) {
	err = s.client.Call(ctx, _addTokensCache, arg, _noReply)
	return
}
