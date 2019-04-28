package v1

import (
	"context"

	v1pb "go-common/app/service/live/xrewardcenter/api/grpc/v1"
	"go-common/app/service/live/xrewardcenter/conf"
	"go-common/app/service/live/xrewardcenter/dao/anchorReward"
	model "go-common/app/service/live/xrewardcenter/model/anchorTask"
	"go-common/library/ecode"
	"go-common/library/log"

	"fmt"
	"time"

	"go-common/app/service/live/gift/api/liverpc/v0"
	"go-common/app/service/live/room/api/liverpc/v2"
	"go-common/app/service/live/xrewardcenter/dao"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

// consts .
const (
	PageSize = int64(20)
	// 任意门道具id
	GiftIDRenYiMen = int64(30035)
	// 默认奖励有效期
	DefaultLifespan = int64(24)
)

// AnchorRewardService struct
type AnchorRewardService struct {
	conf                 *conf.Config
	dao                  *anchorReward.Dao
	cron                 *cron.Cron
	ExpireCountFrequency string
	SetExpireFrequency   string
}

//NewAnchorTaskService init
func NewAnchorTaskService(c *conf.Config) (s *AnchorRewardService) {
	//spew.Dump(c)
	s = &AnchorRewardService{
		conf:                 c,
		cron:                 cron.New(),
		ExpireCountFrequency: c.Cfg.ExpireCountFrequency,
		SetExpireFrequency:   c.Cfg.SetExpireFrequency,
		dao:                  anchorReward.New(c),
	}
	s.addCrontab()
	s.cron.Start()
	return s
}

func (s *AnchorRewardService) addCrontab() (err error) {
	//spew.Dump(s.ExpireCountFrequency)
	err = s.cron.AddFunc(s.ExpireCountFrequency, s.CountExpire)

	if err != nil {
		log.Error("anchorTask.CountExpire AddFunc CountExpire error(%v)", err)
	}

	err = s.cron.AddFunc(s.SetExpireFrequency, s.SetExpire)

	if err != nil {
		log.Error("anchorTask.SetExpire AddFunc SetExpireFrequency error(%v)", err)
	}

	return
}

// MyReward implementation
// * (主播侧)-我的主播奖励(登录态)
//
func (s *AnchorRewardService) MyReward(ctx context.Context, req *v1pb.AnchorTaskMyRewardReq) (resp *v1pb.AnchorTaskMyRewardResp, err error) {

	var (
		iPage = req.Page
	)

	uid := req.GetUid()

	// 获取 uid
	if 0 == uid {
		err = ecode.NoLogin
		//errors.WithMessage(ecode.NoLogin, "账号未登录")
		return
	}

	if iPage <= 0 {
		iPage = 1
	}

	resp = &v1pb.AnchorTaskMyRewardResp{}

	pager, list, err := s.dao.GetByUidPage(ctx, uid, iPage, PageSize, []int64{model.RewardUnUsed})
	if nil != err {
		log.Error("myReward(%+v) error(%v)", req, err)
		err = nil
	}

	for _, v := range list {
		resp.Data = append(resp.Data, &v1pb.AnchorTaskMyRewardResp_RewardObj{
			Id:          v.Id,
			RewardType:  v.RewardType,
			Status:      v.Status,
			RewardId:    v.RewardId,
			Name:        v.Name,
			Icon:        v.Icon,
			AchieveTime: v.AchieveTime,
			ExpireTime:  v.ExpireTime,
			Source:      v.Source,
			RewardIntro: v.RewardIntro,
		})
	}

	resp.Page = &v1pb.AnchorTaskMyRewardResp_Page{
		Page:       pager.Page,
		PageSize:   pager.PageSize,
		TotalPage:  pager.TotalPage,
		TotalCount: pager.TotalPage,
	}

	resp.ExpireCount, _ = s.dao.GetExpireCountCache(ctx, fmt.Sprintf(model.CountExpireUserKey, uid))

	if err := s.dao.SetNewReward(ctx, uid, int64(0)); err != nil {
		log.Error("hasNewRewardMc(%v) error(%v)", uid, err)
	}

	if err := s.dao.ClearExpireCountCache(ctx, fmt.Sprintf(model.CountExpireUserKey, uid)); err != nil {
		log.Error("ClearExpireCountCache(%v) error(%v)", uid, err)
	}

	return
}

// SetExpire changes status if is expired.
func (s *AnchorRewardService) SetExpire() {
	var (
		c        = context.TODO()
		now      = time.Now()
		interval = int64(60)
	)

	log.Info("SetExpire start (%v)", now)

	if b, err := s.dao.SetNxLock(c, model.SetExpireLockKey, interval-1); !b || err != nil {
		log.Info("SetExpire had run (%v,%v)", b, err)
		return
	}

	s.dao.SetExpire(now)

	s.dao.DelLockCache(c, model.SetExpireLockKey)
}

// CountExpire .
func (s *AnchorRewardService) CountExpire() {
	var (
		c        = context.TODO()
		now      = time.Now()
		interval = int64(60)
	)

	log.Info("expireCount start (%v)", now)
	if b, err := s.dao.SetNxLock(c, model.CountExpireLockKey, interval-1); !b || err != nil {
		log.Info("expireCount had run (%v,%v)", b, err)
		return
	}

	// do Count
	s.dao.CountExpire(interval, now)

	s.dao.DelLockCache(c, model.CountExpireLockKey)
	log.Info("expireCount end (%v)", time.Now())
}

// UseRecord implementation
// * (主播侧)-奖励使用记录(登录态)
func (s *AnchorRewardService) UseRecord(ctx context.Context, req *v1pb.AnchorTaskUseRecordReq) (resp *v1pb.AnchorTaskUseRecordResp, err error) {
	resp = &v1pb.AnchorTaskUseRecordResp{}
	resp.Data = []*v1pb.AnchorTaskUseRecordResp_RewardObj{}
	resp.Page = &v1pb.AnchorTaskUseRecordResp_Page{}

	uid := req.GetUid()
	if 0 == uid {
		err = ecode.NoLogin
		return
	}

	iPage := req.GetPage()
	if iPage <= 0 {
		iPage = 1
	}

	pager, list, err := s.dao.GetByUidPage(ctx, uid, iPage, PageSize, []int64{model.RewardExpired, model.RewardUsed})
	//spew.Dump(pager)
	if nil != err {
		log.Error("%v\n", err)
		err = ecode.Error(ecode.ServerErr, "内部错误4")
		return
	}

	if pager == nil || list == nil {
		log.Error("xrc.GetByUidPage nil pager (%v) list (%v)", pager, list)
	}

	//spew.Dump(list)
	for _, v := range list {
		resp.Data = append(resp.Data, &v1pb.AnchorTaskUseRecordResp_RewardObj{
			Id:          v.Id,
			RewardId:    v.RewardId,
			Status:      v.Status,
			Name:        v.Name,
			Icon:        v.Icon,
			AchieveTime: v.AchieveTime,
			UseTime:     v.UseTime,
			ExpireTime:  v.ExpireTime,
			Source:      v.Source,
			RewardIntro: v.RewardIntro,
		})
	}

	resp.Page = &v1pb.AnchorTaskUseRecordResp_Page{
		Page:       pager.Page,
		PageSize:   pager.PageSize,
		TotalPage:  pager.TotalPage,
		TotalCount: pager.TotalCount,
	}
	//spew.Dump(resp)
	return
}

// UseReward implementation
// * (主播侧)-使用奖励(登录态)
//
func (s *AnchorRewardService) UseReward(ctx context.Context, req *v1pb.AnchorTaskUseRewardReq) (resp *v1pb.AnchorTaskUseRewardResp, err error) {
	resp = &v1pb.AnchorTaskUseRewardResp{}
	resp.Result = int64(0)

	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.NoLogin
		return
	}

	id := req.GetId()
	if id <= 0 {
		err = ecode.Error(1, "id不合法")
		return
	}

	aReward, err := s.dao.GetById(id)
	log.Info("useReward reward(%+v), id(%v), uid(%v)", aReward, uid, id)

	if err != nil || aReward == nil {
		err = ecode.Error(5, "奖励不存在")
		return
	}

	if uid != aReward.Uid {
		err = ecode.Error(1, "参数错误")
		return
	}

	if aReward.Status == model.RewardExpired {
		err = ecode.Error(2, "这个奖励已经过期了呢")
		return
	}

	if aReward.Status == model.RewardUsed {
		err = ecode.Error(3, "这个奖励已经被你使用啦~")
		return
	}

	if aReward.Status != model.RewardUnUsed {
		err = ecode.Error(1, "奖励不可用")
		return
	}

	reply, err := dao.RoomAPI.V2Room.GetByIds(ctx, &v2.RoomGetByIdsReq{Ids: []int64{int64(aReward.Roomid)}})

	if err != nil {
		log.Error("call V2Room.GetByIds(%v) error(%v)", aReward.Roomid, err)
		err = errors.Wrap(ecode.Error(-500, "内部错误"), err.Error())
		return
	}
	log.Info("call V2Room.GetByIds (%+v) succ (%+v)", aReward.Roomid, reply)

	if nil == reply || nil == reply.GetData() {
		log.Error("call V2Room.GetByIds roomid(%v) return nil(%+v)", aReward.Roomid, reply)
		err = ecode.Error(-500, "内部错误")
		return
	}

	if reply.GetData()[aReward.Roomid].LiveStatus != 1 {
		err = ecode.Error(4, "为了更好的使用体验，请在开播状态下使用【任意门】哦~")
		return
	}

	rst, err := s.dao.UseReward(id, req.GetUsePlat())

	if err != nil || !rst {
		err = ecode.Error(-500, "内部错误")
		return
	}

	smallTvReq := &v0.SmalltvStartReq{
		Uid:     uid,
		Roomid:  aReward.Roomid,
		GiftId:  GiftIDRenYiMen,
		Num:     1,
		Tid:     0,
		StyleId: 5,
	}
	// 小电视抽奖
	replyGift, err := dao.GiftAPI.V0Smalltv.Start(ctx, smallTvReq)

	if err != nil {
		log.Error("call V0Smalltv(%v) error(%v)", smallTvReq, err)
	} else {
		log.Info("call V0Smalltv (%+v) succ (%+v)", smallTvReq, replyGift)
	}

	resp.Result = int64(1)

	return
}

// IsViewed implementation
// * (主播侧)-奖励和任务红点(登录态)
//
func (s *AnchorRewardService) IsViewed(ctx context.Context, req *v1pb.AnchorTaskIsViewedReq) (resp *v1pb.AnchorTaskIsViewedResp, err error) {
	resp = &v1pb.AnchorTaskIsViewedResp{}
	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.Error(1, "uid不合法")
		return
	}

	resp.RewardShouldNotice, _ = s.dao.HasNewReward(ctx, uid)
	resp.ShowRewardEntry, _ = s.dao.HasReward(ctx, uid)
	resp.Url = "https://live.bilibili.com/p/html/live-app-award/index.html?is_live_webview=1"
	return
}

// AddReward implementation
// (主播侧)-添加主播奖励(内部接口)
// `method:"POST" internal:"true"`
func (s *AnchorRewardService) AddReward(ctx context.Context, req *v1pb.AnchorTaskAddRewardReq) (resp *v1pb.AnchorTaskAddRewardResp, err error) {
	resp = &v1pb.AnchorTaskAddRewardResp{}

	lifespan := req.GetLifespan()
	if lifespan <= 0 {
		lifespan = DefaultLifespan
	}

	exist, err := s.dao.CheckOrderID(ctx, req.GetOrderId())

	if err != nil {
		log.Error("addReward(%+v) error(%v)", req, err)
		resp.Result = int64(0)
		err = nil
		return
	}

	if exist == model.RewardExists {
		err = ecode.Error(1,
			"order already exists!",
		)
		return
	}

	// 红点, 入口
	s.dao.AddReward(ctx, req.GetRewardId(), req.GetUid(), req.GetSource(), req.GetRoomid(), lifespan)

	// 广播
	err = s.dao.SendBroadcastV2(ctx, req.GetUid(), req.GetRoomid(), req.GetRewardId())
	if err != nil {
		log.Error("SendBroadcast(%v) error(%v)", req, err)
	}

	if err != nil {
		resp.Result = int64(0)
	} else {
		resp.Result = int64(1)
		s.dao.SaveOrderID(ctx, req.GetOrderId())
	}

	return
}
