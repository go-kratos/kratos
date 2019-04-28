package gorpc

import (
	"go-common/app/service/main/push/conf"
	"go-common/app/service/main/push/model"
	"go-common/app/service/main/push/service"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/context"
)

// RPC rpc.
type RPC struct {
	s *service.Service
}

// New .
func New(c *conf.Config, s *service.Service) (svc *rpc.Server) {
	r := &RPC{s: s}
	svc = rpc.NewServer(c.RPCServer)
	if err := svc.Register(r); err != nil {
		panic(err)
	}
	return
}

// Ping checks connection success.
func (r *RPC) Ping(c context.Context, arg *struct{}, res *struct{}) (err error) {
	return
}

// Auth check connection success.
func (r *RPC) Auth(c context.Context, arg *rpc.Auth, res *struct{}) (err error) {
	return
}

// AddReport adds report by mid.
func (r *RPC) AddReport(c context.Context, arg *model.ArgReport, res *struct{}) (err error) {
	report := &model.Report{
		APPID:        arg.APPID,
		PlatformID:   arg.PlatformID,
		Mid:          arg.Mid,
		Buvid:        arg.Buvid,
		DeviceToken:  arg.DeviceToken,
		Build:        arg.Build,
		TimeZone:     arg.TimeZone,
		NotifySwitch: arg.NotifySwitch,
		DeviceBrand:  arg.DeviceBrand,
		DeviceModel:  arg.DeviceModel,
		OSVersion:    arg.OSVersion,
		Extra:        arg.Extra,
	}
	err = r.s.AddReport(c, report)
	return
}

// DelInvalidReports deletes invalid reports.
func (r *RPC) DelInvalidReports(c context.Context, arg *model.ArgDelInvalidReport, res *struct{}) (err error) {
	err = r.s.DelInvalidReports(c, arg.Type)
	return
}

// DelReport deletes report.
func (r *RPC) DelReport(c context.Context, arg *model.ArgReport, res *struct{}) (err error) {
	err = r.s.DelReport(c, arg.APPID, arg.Mid, arg.DeviceToken)
	return
}

// AddCallback adds callback data.
func (r *RPC) AddCallback(c context.Context, arg *model.ArgCallback, res *struct{}) (err error) {
	cb := &model.Callback{
		Task:     arg.Task,
		APP:      arg.APP,
		Platform: arg.Platform,
		Mid:      arg.Mid,
		Pid:      arg.Pid,
		Token:    arg.Token,
		Buvid:    arg.Buvid,
		Click:    arg.Click,
		Extra:    arg.Extra,
	}
	err = r.s.AddCallback(c, cb)
	return
}

// AddReportCache adds report cache.
func (r *RPC) AddReportCache(c context.Context, arg *model.ArgReport, res *struct{}) (err error) {
	report := &model.Report{
		ID:           arg.ID,
		APPID:        arg.APPID,
		PlatformID:   arg.PlatformID,
		Mid:          arg.Mid,
		Buvid:        arg.Buvid,
		DeviceToken:  arg.DeviceToken,
		Build:        arg.Build,
		TimeZone:     arg.TimeZone,
		NotifySwitch: arg.NotifySwitch,
		DeviceBrand:  arg.DeviceBrand,
		DeviceModel:  arg.DeviceModel,
		OSVersion:    arg.OSVersion,
		Extra:        arg.Extra,
	}
	err = r.s.AddReportCache(c, report)
	return
}

// AddUserReportCache adds user report cache.
func (r *RPC) AddUserReportCache(c context.Context, arg *model.ArgUserReports, res *struct{}) (err error) {
	err = r.s.AddUserReportCache(c, arg.Mid, arg.Reports)
	return
}

// Setting gets user push switch setting.
func (r *RPC) Setting(c context.Context, arg *model.ArgMid, res *map[int]int) (err error) {
	*res, err = r.s.Setting(c, arg.Mid)
	return
}

// SetSetting sets user push switch setting.
func (r *RPC) SetSetting(c context.Context, arg *model.ArgSetting, res *struct{}) (err error) {
	err = r.s.SetSetting(c, arg.Mid, arg.Type, arg.Value)
	return
}

// AddMidProgress add mid count number to task progress field
func (r *RPC) AddMidProgress(c context.Context, arg *model.ArgMidProgress, res *struct{}) (err error) {
	err = r.s.AddMidProgress(c, arg.Task, arg.MidTotal, arg.MidValid)
	return
}

// AddTokenCache add token cache
func (r *RPC) AddTokenCache(ctx context.Context, arg *model.ArgReport, res *struct{}) (err error) {
	report := &model.Report{
		APPID:        arg.APPID,
		PlatformID:   arg.PlatformID,
		Mid:          arg.Mid,
		Buvid:        arg.Buvid,
		DeviceToken:  arg.DeviceToken,
		Build:        arg.Build,
		TimeZone:     arg.TimeZone,
		NotifySwitch: arg.NotifySwitch,
		DeviceBrand:  arg.DeviceBrand,
		DeviceModel:  arg.DeviceModel,
		OSVersion:    arg.OSVersion,
		Extra:        arg.Extra,
	}
	err = r.s.AddTokenCache(ctx, report)
	return
}

// AddTokensCache add token cache
func (r *RPC) AddTokensCache(ctx context.Context, arg *model.ArgReports, res *struct{}) (err error) {
	rs := make(map[string]*model.Report, len(arg.Reports))
	for _, v := range arg.Reports {
		rs[v.DeviceToken] = v
	}
	err = r.s.AddTokensCache(ctx, rs)
	return
}
