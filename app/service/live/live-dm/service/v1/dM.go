package v1

import (
	"context"
	"regexp"
	"time"

	v1pb "go-common/app/service/live/live-dm/api/grpc/v1"
	"go-common/app/service/live/live-dm/conf"
	"go-common/app/service/live/live-dm/dao"
	"go-common/library/ecode"
	"go-common/library/log"
)

// DMService struct
type DMService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//SendMsg 发送弹幕逻辑参数
type SendMsg struct {
	SendMsgReq   *v1pb.SendMsgReq
	Dmservice    *DMService
	SendMsgResp  *v1pb.SendMsgResp
	UserInfo     *dao.UserInfo
	RoomConf     *dao.RoomConf
	UserBindInfo *dao.UserBindInfo
	DMconf       *dao.DMConf
	TitleConf    *dao.CommentTitle
	UserScore    *dao.UserScore
	LimitConf    *dao.LimitConf
}

//
var (
	reMsg   = regexp.MustCompile(`#\s+#`)
	msgMust = regexp.MustCompile(`(/n)|(\n)`)
)

//NewDMService init
func NewDMService(c *conf.Config) (s *DMService) {
	s = &DMService{
		conf: c,
		dao:  dao.New(c),
	}
	return s
}

// SendMsg implementation
func (s *DMService) SendMsg(ctx context.Context, req *v1pb.SendMsgReq) (resp *v1pb.SendMsgResp, err error) {
	sdm := &SendMsg{
		SendMsgReq: req,
		Dmservice:  s,
		UserInfo: &dao.UserInfo{
			MedalInfo: &dao.FansMedalInfo{},
		},
		RoomConf:     &dao.RoomConf{},
		UserBindInfo: &dao.UserBindInfo{},
		TitleConf:    &dao.CommentTitle{},
		DMconf:       &dao.DMConf{},
		UserScore:    &dao.UserScore{},
		LimitConf: &dao.LimitConf{
			AreaLimit:        s.conf.DmRules.AreaLimit,
			AllUserLimit:     s.conf.DmRules.AllUserLimit,
			LevelLimitStatus: s.conf.DmRules.LevelLimitStatus,
			LevelLimit:       s.conf.DmRules.LevelLimit,
			RealName:         s.conf.DmRules.RealName,
			PhoneLimit:       s.conf.DmRules.PhoneLimit,
			MsgLength:        s.conf.DmRules.MsgLength,
			DMPercent:        s.conf.DmRules.DmPercent,
			DmNum:            s.conf.DmRules.DmNum,
			DMwhitelist:      s.conf.DmRules.DMwhitelist,
			DMwhitelistID:    s.conf.DmRules.DMwhiteListID,
		},
	}

	sdm.LimitConf.GetDMCheckConf()

	resp = &v1pb.SendMsgResp{}
	if req.GetLancer() == nil {
		req.Lancer = &v1pb.Lancer{}
	}

	//限制字体大小为25
	sdm.SendMsgReq.Fontsize = 25

	//特殊字符处理
	req.Msg = reMsg.ReplaceAllString(req.Msg, " ")
	sdm.SendMsgReq.Msg = msgMust.ReplaceAllString(req.Msg, "")

	// 发送弹幕的频率控制
	if err = rateLimit(ctx, sdm); err != nil {
		if perr, ok := err.(*ecode.Status); ok {
			resp.Code = int32(perr.Code())
			resp.LimitMsg = perr.Message()
			resp.IsLimit = true
			err = nil

			log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]", perr.Message(), req.Uid, req.Roomid, req.Msg)
			dao.InfocDMErr.Info(req.Roomid, req.Uid, req.Msg, perr.Message(),
				time.Now().Unix(), req.Platform, req.Lancer.Build, req.Lancer.Buvid,
				req.Lancer.UserAgent, req.Lancer.Refer)
			return
		}
		log.Error("DM limit err:%+v", err)
	}

	// 获取弹幕禁言检查依赖数据
	if err = getCheckMsg(ctx, sdm); err != nil {
		if perr, ok := err.(*ecode.Status); ok {
			resp.Code = int32(perr.Code())
			resp.LimitMsg = perr.Message()
			resp.IsLimit = false
			err = nil

			return
		}
		log.Error("DM get check message err:%+v", err)

		dao.InfocDMErr.Info(req.Roomid, req.Uid, req.Msg, "弹幕禁言检查依赖获取失败",
			time.Now().Unix(), req.Platform, req.Lancer.Build, req.Lancer.Buvid,
			req.Lancer.UserAgent, req.Lancer.Refer)
	}
	// 弹幕禁言检查
	if err = checkLegitimacy(ctx, sdm); err != nil {
		if perr, ok := err.(*ecode.Status); ok {
			resp.Code = int32(perr.Code())
			resp.LimitMsg = perr.Message()
			resp.IsLimit = true
			err = nil

			return
		}
		dao.InfocDMErr.Info(req.Roomid, req.Uid, req.Msg, "弹幕禁言检查失败",
			time.Now().Unix(), req.Platform, req.Lancer.Build, req.Lancer.Buvid,
			req.Lancer.UserAgent, req.Lancer.Refer)
		return
	}

	// 获取弹幕发送依赖数据
	if err = getDMconfig(ctx, sdm); err != nil {
		if perr, ok := err.(*ecode.Status); ok {
			resp.Code = int32(perr.Code())
			resp.LimitMsg = perr.Message()
			resp.IsLimit = false
			err = nil

			return
		}
		dao.InfocDMErr.Info(req.Roomid, req.Uid, req.Msg, "获取弹幕发送依赖数据失败",
			time.Now().Unix(), req.Platform, req.Lancer.Build, req.Lancer.Buvid,
			req.Lancer.UserAgent, req.Lancer.Refer)
		return
	}

	// 发送弹幕
	if err = send(ctx, sdm); err != nil {
		if perr, ok := err.(*ecode.Status); ok {
			resp.Code = int32(perr.Code())
			resp.LimitMsg = perr.Message()
			resp.IsLimit = false
			err = nil

			return
		}
		dao.InfocDMErr.Info(req.Roomid, req.Uid, req.Msg, "弹幕广播失败",
			time.Now().Unix(), req.Platform, req.Lancer.Build, req.Lancer.Buvid,
			req.Lancer.UserAgent, req.Lancer.Refer)
		return
	}

	log.Info("LIMITDM, reason: %s, [detail: uid:%d, roomid:%d, msg:%s]", "弹幕发送成功", req.Uid, req.Roomid, req.Msg)
	dao.InfocDMSend.Info(req.Roomid, req.Uid, req.Msg, req.Ip, time.Now().Unix(),
		sdm.DMconf.Color, req.Mode, req.Platform, req.Lancer.Build,
		req.Lancer.Buvid, req.Lancer.UserAgent, req.Lancer.Refer, req.Lancer.Cookie,
		req.Msgtype)
	return
}

// GetHistory implementation
func (s *DMService) GetHistory(ctx context.Context, req *v1pb.HistoryReq) (resp *v1pb.HistoryResp, err error) {
	resp = &v1pb.HistoryResp{}
	rest, err := s.dao.GetHistoryData(ctx, req.Roomid)
	if err != nil {
		return resp, err
	}
	resp.Admin = make([]string, 0, 10)
	resp.Room = make([]string, 0, 10)
	resp.Admin = append(resp.Admin, rest["admin"]...)
	resp.Room = append(resp.Room, rest["room"]...)
	return
}

func lancer(sdm *SendMsg, reason string) {
	req := sdm.SendMsgReq
	dao.InfocDMErr.Info(req.Roomid, req.Uid, req.Msg, reason,
		time.Now().Unix(), req.Platform, req.Lancer.Build, req.Lancer.Buvid,
		req.Lancer.UserAgent, req.Lancer.Refer)
}
