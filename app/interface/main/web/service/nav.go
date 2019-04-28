package service

import (
	"context"

	"go-common/app/interface/main/web/model"
	accmdl "go-common/app/service/main/account/api"
	coupon "go-common/app/service/main/coupon/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// Nav api service
func (s *Service) Nav(c context.Context, mid int64, cookie string) (resp *model.NavResp, err error) {
	var (
		wallet    *model.Wallet
		hasShop   bool
		shopURL   string
		allowance int
	)
	profile := new(accmdl.ProfileStatReply)
	eg, egCtx := errgroup.WithContext(c)
	eg.Go(func() error {
		var e error
		if profile, e = s.accClient.ProfileWithStat3(egCtx, &accmdl.MidReq{Mid: mid}); e != nil {
			log.Error("s.accClient.ProfileWithStat3(%d) error %v", mid, e)
			profile = model.DefaultProfile
			profile.Profile.Mid = mid
		}
		return nil
	})
	eg.Go(func() error {
		var shop *model.ShopInfo
		var e error
		if shop, e = s.dao.ShopInfo(egCtx, mid); e == nil && shop != nil {
			hasShop = true
			shopURL = shop.JumpURL
		} else {
			log.Warn("s.dao.ShopInfo(%v) error(%+v)", mid, e)
		}
		return nil
	})
	eg.Go(func() error {
		var e error
		if wallet, e = s.dao.Wallet(egCtx, mid); e != nil || wallet == nil {
			log.Error("s.dao.Wallet(%d) error(%v)", mid, e)
			if wallet, e = s.dao.OldWallet(egCtx, mid); e != nil || wallet == nil {
				log.Error("s.dao.OldWallet(%d) error(%v)", mid, e)
			}
		} else {
			log.Info("account wallet mid(%d)", mid)
		}
		return nil
	})
	eg.Go(func() error {
		var e error
		if allowance, e = s.coupon.AllowanceCount(egCtx, &coupon.ArgAllowanceMid{Mid: mid}); e != nil {
			log.Error("s.coupon.AllowanceCount(%d) error(%v)", mid, e)
		}
		return nil
	})
	eg.Wait()
	resp = &model.NavResp{
		IsLogin:        true,
		EmailVerified:  int(profile.Profile.EmailStatus),
		Face:           profile.Profile.Face,
		Mid:            profile.Profile.Mid,
		MobileVerified: int(profile.Profile.TelStatus),
		Coins:          profile.Coins,
		Moral:          float32(profile.Profile.Moral),
		Pendant:        profile.Profile.Pendant,
		Uname:          profile.Profile.Name,
		VipDueDate:     profile.Profile.Vip.DueDate,
		VipStatus:      int(profile.Profile.Vip.Status),
		VipType:        int(profile.Profile.Vip.Type),
		VipPayType:     profile.Profile.Vip.VipPayType,
		Wallet:         wallet,
		HasShop:        hasShop,
		ShopURL:        shopURL,
		AllowanceCount: allowance,
	}
	if profile.Profile.Official.Role == 0 {
		resp.OfficialVerify.Type = -1
	} else {
		if profile.Profile.Official.Role <= 2 {
			resp.OfficialVerify.Type = 0
		} else {
			resp.OfficialVerify.Type = 1
		}
		resp.OfficialVerify.Desc = profile.Profile.Official.Title
	}
	resp.LevelInfo.Cur = int(profile.LevelInfo.Cur)
	resp.LevelInfo.Min = int(profile.LevelInfo.Min)
	resp.LevelInfo.NowExp = int(profile.LevelInfo.NowExp)
	resp.LevelInfo.NextExp = profile.LevelInfo.NextExp
	if profile.LevelInfo.NextExp == -1 {
		resp.LevelInfo.NextExp = "--"
	}
	return
}
