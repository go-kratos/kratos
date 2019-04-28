package grpc

import (
	"context"

	pb "go-common/app/service/main/push/api/grpc/v1"
	"go-common/app/service/main/push/model"
	"go-common/app/service/main/push/service"
	"go-common/library/net/rpc/warden"
)

// New warden rpc server
func New(c *warden.ServerConfig, svc *service.Service) *warden.Server {
	svr := warden.NewServer(c)
	pb.RegisterPushServer(svr.Server(), &server{svc: svc})
	ws, err := svr.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svc *service.Service
}

// AddReport 上报
func (s *server) AddReport(ctx context.Context, req *pb.AddReportRequest) (reply *pb.AddReportReply, err error) {
	reply = new(pb.AddReportReply)
	if req.Report == nil {
		return
	}
	r := req.Report
	err = s.svc.AddReport(ctx, &model.Report{
		APPID:        int64(r.APPID),
		PlatformID:   int(r.PlatformID),
		Mid:          r.Mid,
		Buvid:        r.Buvid,
		DeviceToken:  r.DeviceToken,
		Build:        int(r.Build),
		TimeZone:     int(r.TimeZone),
		NotifySwitch: int(r.NotifySwitch),
		DeviceBrand:  r.DeviceBrand,
		DeviceModel:  r.DeviceModel,
		OSVersion:    r.OSVersion,
		Extra:        r.Extra,
	})
	return
}

// DelReport 删除上报
func (s *server) DelReport(ctx context.Context, req *pb.DelReportRequest) (reply *pb.DelReportReply, err error) {
	reply = new(pb.DelReportReply)
	err = s.svc.DelReport(ctx, int64(req.APPID), req.Mid, req.DeviceToken)
	return
}

// DelInvalidReports 删除无效token
func (s *server) DelInvalidReports(ctx context.Context, req *pb.DelInvalidReportsRequest) (reply *pb.DelInvalidReportsReply, err error) {
	reply = new(pb.DelInvalidReportsReply)
	err = s.svc.DelInvalidReports(ctx, int(req.Type))
	return
}

// AddReportCache 上报缓存
func (s *server) AddReportCache(ctx context.Context, req *pb.AddReportCacheRequest) (reply *pb.AddReportCacheReply, err error) {
	reply = new(pb.AddReportCacheReply)
	if req.Report == nil {
		return
	}
	r := req.Report
	err = s.svc.AddReportCache(ctx, &model.Report{
		APPID:        int64(r.APPID),
		PlatformID:   int(r.PlatformID),
		Mid:          r.Mid,
		Buvid:        r.Buvid,
		DeviceToken:  r.DeviceToken,
		Build:        int(r.Build),
		TimeZone:     int(r.TimeZone),
		NotifySwitch: int(r.NotifySwitch),
		DeviceBrand:  r.DeviceBrand,
		DeviceModel:  r.DeviceModel,
		OSVersion:    r.OSVersion,
		Extra:        r.Extra,
	})
	return
}

// AddUserReportCache 用户的上报缓存
func (s *server) AddUserReportCache(ctx context.Context, req *pb.AddUserReportCacheRequest) (reply *pb.AddUserReportCacheReply, err error) {
	reply = new(pb.AddUserReportCacheReply)
	var reports []*model.Report
	for _, r := range req.Reports {
		reports = append(reports, &model.Report{
			APPID:        int64(r.APPID),
			PlatformID:   int(r.PlatformID),
			Mid:          r.Mid,
			Buvid:        r.Buvid,
			DeviceToken:  r.DeviceToken,
			Build:        int(r.Build),
			TimeZone:     int(r.TimeZone),
			NotifySwitch: int(r.NotifySwitch),
			DeviceBrand:  r.DeviceBrand,
			DeviceModel:  r.DeviceModel,
			OSVersion:    r.OSVersion,
			Extra:        r.Extra,
		})
	}
	err = s.svc.AddUserReportCache(ctx, req.Mid, reports)
	return
}

// AddTokenCache token缓存
func (s *server) AddTokenCache(ctx context.Context, req *pb.AddTokenCacheRequest) (reply *pb.AddTokenCacheReply, err error) {
	reply = new(pb.AddTokenCacheReply)
	if req.Report == nil {
		return
	}
	r := req.Report
	err = s.svc.AddTokenCache(ctx, &model.Report{
		APPID:        int64(r.APPID),
		PlatformID:   int(r.PlatformID),
		Mid:          r.Mid,
		Buvid:        r.Buvid,
		DeviceToken:  r.DeviceToken,
		Build:        int(r.Build),
		TimeZone:     int(r.TimeZone),
		NotifySwitch: int(r.NotifySwitch),
		DeviceBrand:  r.DeviceBrand,
		DeviceModel:  r.DeviceModel,
		OSVersion:    r.OSVersion,
		Extra:        r.Extra,
	})
	return
}

// AddTokensCache 批量token缓存
func (s *server) AddTokensCache(ctx context.Context, req *pb.AddTokensCacheRequest) (reply *pb.AddTokensCacheReply, err error) {
	reply = new(pb.AddTokensCacheReply)
	if len(req.Reports) == 0 {
		return
	}
	reports := make(map[string]*model.Report, len(req.Reports))
	for _, v := range req.Reports {
		reports[v.DeviceToken] = &model.Report{
			APPID:        int64(v.APPID),
			PlatformID:   int(v.PlatformID),
			Mid:          v.Mid,
			Buvid:        v.Buvid,
			DeviceToken:  v.DeviceToken,
			Build:        int(v.Build),
			TimeZone:     int(v.TimeZone),
			NotifySwitch: int(v.NotifySwitch),
			DeviceBrand:  v.DeviceBrand,
			DeviceModel:  v.DeviceModel,
			OSVersion:    v.OSVersion,
			Extra:        v.Extra,
		}
	}
	err = s.svc.AddTokensCache(ctx, reports)
	return
}

// AddCallback 回执
func (s *server) AddCallback(ctx context.Context, req *pb.AddCallbackRequest) (reply *pb.AddCallbackReply, err error) {
	reply = new(pb.AddCallbackReply)
	cb := &model.Callback{
		Task:     req.Task,
		APP:      req.APP,
		Platform: int(req.Platform),
		Mid:      int64(req.Mid),
		Pid:      int(req.Pid),
		Token:    req.Token,
		Buvid:    req.Buvid,
		Click:    uint8(req.Click),
	}
	if req.Extra != nil {
		cb.Extra = &model.CallbackExtra{
			Status:  int(req.Extra.Status),
			Channel: int(req.Extra.Channel),
		}
	}
	err = s.svc.AddCallback(ctx, cb)
	return
}

// AddMidProgress 记录任务中mid数
func (s *server) AddMidProgress(ctx context.Context, req *pb.AddMidProgressRequest) (reply *pb.AddMidProgressReply, err error) {
	reply = new(pb.AddMidProgressReply)
	err = s.svc.AddMidProgress(ctx, req.Task, req.MidTotal, req.MidValid)
	return
}

// Setting 获取用户业务开关配置
func (s *server) Setting(ctx context.Context, req *pb.SettingRequest) (reply *pb.SettingReply, err error) {
	reply = new(pb.SettingReply)
	res, err := s.svc.Setting(ctx, req.Mid)
	if err != nil || res == nil {
		return
	}
	reply.Settings = make(map[int32]int32, len(res))
	for k, v := range res {
		reply.Settings[int32(k)] = int32(v)
	}
	return
}

// SetSetting 设置用户业务开关
func (s *server) SetSetting(ctx context.Context, req *pb.SetSettingRequest) (reply *pb.SetSettingReply, err error) {
	reply = new(pb.SetSettingReply)
	err = s.svc.SetSetting(ctx, req.Mid, int(req.Type), int(req.Value))
	return
}
