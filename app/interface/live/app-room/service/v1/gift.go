package v1

import (
	"context"
	v1pb "go-common/app/interface/live/app-room/api/http/v1"
	"go-common/app/interface/live/app-room/conf"
	"go-common/app/interface/live/app-room/dao"
	"go-common/app/interface/live/app-room/model"
	"go-common/app/service/live/gift/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"strconv"
	"time"
)

// GiftService struct
type GiftService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	client v1.GiftClient
	dao    *dao.Dao
}

const (
	// HasShow 已经展示
	HasShow = 302
	// StopPush 停止推送
	StopPush = 403
)

//NewGiftService init
func NewGiftService(c *conf.Config) (s *GiftService) {
	s = &GiftService{
		conf: c,
		dao:  dao.New(c),
	}
	cli, err := v1.NewClient(conf.Conf.Warden)
	if err != nil {
		panic(err)
	}
	s.client = cli
	return s
}

// DailyBag implementation
// `method:"GET" midware:"guest"`
func (s *GiftService) DailyBag(ctx context.Context, req *v1pb.DailyBagReq) (resp *v1pb.DailyBagResp, err error) {
	resp = &v1pb.DailyBagResp{}
	mid, _ := metadata.Value(ctx, metadata.Mid).(int64)
	if mid <= 0 {
		err = ecode.NoLogin
		return
	}
	//mid = 88895029
	ret, err := s.client.DailyBag(ctx, &v1.DailyBagReq{
		Uid: mid,
	})
	if err != nil {
		err = ecode.Error(-1, "系统错误")
		return
	}
	resp = &v1pb.DailyBagResp{
		BagStatus:       ret.BagStatus,
		BagExpireStatus: ret.BagExpireStatus,
		BagToast: &v1pb.DailyBagResp_BagToast{
			ToastStatus:  ret.BagToast.ToastStatus,
			ToastMessage: ret.BagToast.ToastMessage,
		},
		BagList: make([]*v1pb.DailyBagResp_BagList, 0),
	}
	for _, v := range ret.BagList {
		s := &v1pb.DailyBagResp_BagList_Source{}
		if v.Source != nil {
			s = &v1pb.DailyBagResp_BagList_Source{
				MedalId:   v.Source.MedalId,
				MedalName: v.Source.MedalName,
				Level:     v.Source.Level,
				UserLevel: v.Source.UserLevel,
			}
		}
		tmp := &v1pb.DailyBagResp_BagList{
			Type:    v.Type,
			BagName: v.BagName,
			Source:  s,
		}
		for _, v2 := range v.GiftList {
			tmp.GiftList = append(tmp.GiftList, &v1pb.DailyBagResp_BagList_GiftList{
				GiftId:   v2.GiftId,
				GiftNum:  v2.GiftNum,
				ExpireAt: v2.ExpireAt,
			})
		}
		resp.BagList = append(resp.BagList, tmp)
	}
	return
}

// 实际充值的金瓜子数 最小单位是1000
func realRechargeGold(gold int64) (realGold int64) {
	if gold%1000 == 0 {
		realGold = gold
	} else {
		realGold = gold + 1000 - gold%1000
	}
	return
}
func day(t time.Time) int {
	s := t.Format("2")
	d, _ := strconv.Atoi(s)
	return d
}
func yearMonthNum(t time.Time) int64 {
	yearMonth, _ := strconv.ParseInt(t.Format("200601"), 10, 64)
	return yearMonth
}

// 银瓜子是否提醒
func (s *GiftService) silverNeedTipRecharge(ctx context.Context, mid int64, req *v1pb.NeedTipRechargeReq) (resp *v1pb.NeedTipRechargeResp, err error) {
	resp = &v1pb.NeedTipRechargeResp{}
	d := day(time.Now())
	dHit := false
	for _, v := range s.conf.Gift.RechargeTip.SilverTipDays {
		if v == d {
			dHit = true
			break
		}
	}
	if !dHit { // 日期
		log.Info("silver not show because day not hit %d mid:%d", d, mid)
		return
	}
	yearMonth := yearMonthNum(time.Now())
	defer func() {
		if resp.Show == 1 {
			s.dao.AsyncSetUserConf(ctx, mid, model.SilverTarget, yearMonth)
		}
	}()
	// 每月只推一次
	userConf, err := s.dao.GetUserConf(ctx, mid, model.SilverTarget, []int64{yearMonth, StopPush})
	if err != nil {
		err = ecode.ServerErr
		return
	}
	if userConf.IsSet(yearMonth) || userConf.IsSet(StopPush) {
		log.Info("silver not show because userConf yearMonth:%v stopPush:%v mid:%d", userConf.IsSet(yearMonth), userConf.IsSet(StopPush), mid)
		return
	}

	w, err := s.dao.PayCenterWallet(ctx, mid, req.Platform)
	if err != nil {
		return
	}
	if w.CouponBalance < 1 { // bp coupon >= 1
		log.Info("silver not show because couponBalance less than 1 %f mid:%d", w.CouponBalance, mid)
		return
	}
	resp.Bp = w.BcoinBalance
	resp.BpCoupon = w.CouponBalance
	resp.Show = 1
	resp.RechargeGold = 0
	log.Info("SilverShow mid:%d coupon:%f", mid, w.CouponBalance)
	return
}

// 金瓜子是否提醒
func (s *GiftService) goldNeedTipRecharge(ctx context.Context, mid int64, req *v1pb.NeedTipRechargeReq) (resp *v1pb.NeedTipRechargeResp, err error) {
	resp = &v1pb.NeedTipRechargeResp{}
	defer func() {
		if resp.Show == 1 {
			s.dao.AsyncSetUserConf(ctx, mid, model.GoldTarget, HasShow) // 设置已经展示过
		}
	}()
	userConf, err := s.dao.GetUserConf(ctx, mid, model.GoldTarget, []int64{HasShow}) // 是否已经展示过
	if err != nil {
		err = ecode.ServerErr
		return
	}
	if userConf.IsSet(HasShow) {
		log.Info("gold not show because has show mid:%d", mid)
		return
	}

	eg, errCtx := errgroup.WithContext(ctx)
	var payCenterWallet *model.Wallet
	var liveWallet *model.LiveWallet
	eg.Go(func() error {
		var pErr error
		payCenterWallet, pErr = s.dao.PayCenterWallet(errCtx, mid, req.Platform)
		return pErr
	})
	eg.Go(func() error {
		var lErr error
		liveWallet, lErr = s.dao.LiveWallet(errCtx, mid, req.Platform)
		return lErr
	})
	err = eg.Wait()
	if err != nil {
		return
	}

	if liveWallet.GoldPayCnt > 0 { // 历史上没有消费过金瓜子
		log.Info("gold not show because gold pay cnt lt 0  %d mid:%d", liveWallet.GoldPayCnt, mid)
		return
	}

	if payCenterWallet.BcoinBalance < 1 { // bp余额大于1
		log.Info("gold not show because bcoin lt 1 %f mid:%d", payCenterWallet.BcoinBalance, mid)
		return
	}

	bpGold := int64(payCenterWallet.BcoinBalance * 1000)
	realRechargeGold := realRechargeGold(req.NeedGold)
	if bpGold < realRechargeGold { // 金瓜子差值
		log.Info("gold not show because bcoin lt gold bcoin:%f gold:%d mid:%d", payCenterWallet.BcoinBalance, realRechargeGold, mid)
		return
	}
	resp.Bp = payCenterWallet.BcoinBalance
	resp.BpCoupon = payCenterWallet.CouponBalance
	resp.Show = 1
	resp.RechargeGold = realRechargeGold
	log.Info("GoldShow mid:%d bp:%f goldPayCnt:%d needGold:%d", mid, payCenterWallet.BcoinBalance, liveWallet.GoldPayCnt, req.NeedGold)
	return
}

// NeedTipRecharge implementation
//
// `midware:"auth"`
func (s *GiftService) NeedTipRecharge(ctx context.Context, req *v1pb.NeedTipRechargeReq) (resp *v1pb.NeedTipRechargeResp, err error) {
	mid := metadata.Value(ctx, metadata.Mid).(int64)
	if req.From == v1pb.From_Gold {
		return s.goldNeedTipRecharge(ctx, mid, req)
	} else if req.From == v1pb.From_Silver {
		return s.silverNeedTipRecharge(ctx, mid, req)
	}
	err = ecode.RequestErr
	resp = &v1pb.NeedTipRechargeResp{}
	return
}

// TipRechargeAction implementation
//
// `midware:"auth"`
// `method:"post"`
func (s *GiftService) TipRechargeAction(ctx context.Context, req *v1pb.TipRechargeActionReq) (resp *v1pb.TipRechargeActionResp, err error) {
	resp = &v1pb.TipRechargeActionResp{}
	mid := metadata.Value(ctx, metadata.Mid).(int64)
	if req.From == v1pb.From_Silver && req.Action == v1pb.UserAction_StopPush {
		err = s.dao.SetUserConf(ctx, mid, model.SilverTarget, StopPush)
	} else {
		err = ecode.RequestErr
	}

	return
}

// GiftConfig implementation
func (s *GiftService) GiftConfig(ctx context.Context, req *v1pb.GiftConfigReq) (resp *v1pb.GiftConfigResp, err error) {
	resp = &v1pb.GiftConfigResp{
		List: make([]*v1pb.GiftConfigResp_Config, 0),
	}
	ret, err := s.client.GiftConfig(ctx, &v1.GiftConfigReq{
		Platform: req.Platform,
		Build:    req.Build,
	})
	if err != nil {
		log.Error("get gift config err,%v", err)
		err = ecode.Error(-1, "系统错误")
		return
	}
	for _, v := range ret.Data {
		countMap := make([]*v1pb.GiftConfigResp_CountMap, 0)
		for _, c := range v.CountMap {
			tmp := &v1pb.GiftConfigResp_CountMap{
				Num:  c.Num,
				Text: c.Text,
			}
			countMap = append(countMap, tmp)
		}
		d := &v1pb.GiftConfigResp_Config{
			Id:                   v.Id,
			Name:                 v.Name,
			Price:                v.Price,
			Type:                 v.Type,
			CoinType:             v.CoinType,
			BagGift:              v.BagGift,
			Effect:               v.Effect,
			CornerMark:           v.CornerMark,
			Broadcast:            v.Broadcast,
			Draw:                 v.Draw,
			StayTime:             v.StayTime,
			AnimationFrameNum:    v.AnimationFrameNum,
			Desc:                 v.Desc,
			Rule:                 v.Rule,
			Rights:               v.Rights,
			PrivilegeRequired:    v.PrivilegeRequired,
			CountMap:             countMap,
			ImgBasic:             v.ImgBasic,
			ImgDynamic:           v.ImgDynamic,
			FrameAnimation:       v.FrameAnimation,
			Gif:                  v.Gif,
			Webp:                 v.Webp,
			FullScWeb:            v.FullScWeb,
			FullScHorizontal:     v.FullScHorizontal,
			FullScVertical:       v.FullScVertical,
			FullScHorizontalSvga: v.FullScHorizontalSvga,
			FullScVerticalSvga:   v.FullScVerticalSvga,
			BulletHead:           v.BulletHead,
			BulletTail:           v.BulletTail,
			LimitInterval:        v.LimitInterval,
		}
		resp.List = append(resp.List, d)
	}

	return
}
