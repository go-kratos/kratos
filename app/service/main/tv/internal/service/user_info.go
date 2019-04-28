package service

import (
	"context"
	"time"

	"go-common/app/service/main/tv/internal/model"
	"go-common/app/service/main/tv/internal/service/validator"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) UserInfo(c context.Context, mid int64) (ui *model.UserInfo, err error) {
	ui, err = s.dao.UserInfoByMid(c, int64(mid))
	if err != nil {
		return nil, err
	}
	if ui.IsEmpty() {
		return nil, ecode.NothingFound
	}
	if ui.IsExpired() {
		ui.MarkExpired()
		s.ExpireUserAsync(c, ui.Mid)
	}
	return ui, nil
}

func (s *Service) YstUserInfo(c context.Context, req *model.YstUserInfoReq) (ui *model.UserInfo, err error) {
	v := &validator.SignerValidator{
		Sign:   req.Sign,
		Signer: s.dao.Signer(),
		Val:    req,
	}
	if err = v.Validate(); err != nil {
		log.Error("signValidator.Validate(%+v) err(%+v)", req, err)
		return
	}
	ui, err = s.dao.UserInfoByMid(c, int64(req.Mid))
	if err != nil {
		return nil, err
	}
	if ui.IsEmpty() {
		return nil, ecode.NothingFound
	}
	if ui.IsExpired() {
		ui.MarkExpired()
		s.ExpireUserAsync(c, ui.Mid)
	}
	return ui, nil
}

func (s *Service) ExpireUserAsync(c context.Context, mid int64) {
	s.mission(func() {
		if err := s.ExpireUser(context.TODO(), mid); err != nil {
			log.Error("s.ExpireUser(%d) err(%+v)", mid, err)
		}
		log.Info("s.ExpireUserAsync(%d) msg(success)", mid)
	})
}

func (s *Service) ExpireUser(c context.Context, mid int64) (err error) {
	var (
		ui *model.UserInfo
		tx *xsql.Tx
	)
	// check overdue time
	if ui, err = s.dao.RawUserInfoByMid(c, mid); err != nil {
		return
	}
	if !ui.IsExpired() {
		return nil
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	defer func() {
		err = s.dao.EndTran(tx, err)
	}()
	ui.MarkExpired()
	if err = s.dao.TxUpdateUserInfo(c, tx, ui); err != nil {
		return
	}
	s.flushUserInfoAsync(context.TODO(), mid)
	return nil
}

func (s *Service) incrVipDuration(c context.Context, tx *xsql.Tx, ui *model.UserInfo, d time.Duration) (err error) {
	// update user status
	ui.Status = model.VipStatusActive
	// update user overdue time
	now := time.Now()
	nowUnix := time.Now().Unix()
	ui.RecentPayTime = xtime.Time(nowUnix)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Add(time.Hour * 24)
	if int64(ui.OverdueTime) <= nowUnix {
		ui.OverdueTime = xtime.Time(startDate.Unix() + int64(d.Seconds()))
	} else {
		ui.OverdueTime = xtime.Time(int64(ui.OverdueTime) + int64(d.Seconds()))
	}
	if err = s.dao.TxUpdateUserInfo(c, tx, ui); err != nil {
		return
	}
	log.Info("s.incrVipDuration(%+v, %+v)", ui, d)
	return nil
}

func (s *Service) RenewPanel(c context.Context, mid int64) (p *model.PanelPriceConfig, err error) {
	uc, err := s.dao.UserContractByMid(c, mid)
	if err != nil {
		return
	}
	po, err := s.dao.PayOrderByOrderNo(c, uc.OrderNo)
	if err != nil {
		return
	}
	p, err = s.PanelPriceConfigByProductId(c, po.ProductId)
	if err != nil {
		return
	}
	if p == nil {
		return nil, ecode.TVIPPanelNotFound
	}
	return
}

func (s *Service) RenewVip(c context.Context, mid int64) (err error) {
	ui, err := s.dao.RawUserInfoByMid(c, mid)
	if err != nil {
		return
	}
	// validate
	validator := &validator.RenewVipValidator{
		UserInfo:     ui,
		FromDuration: s.c.PAY.RenewFromDuration,
		ToDuration:   s.c.PAY.RenewToDuration,
	}
	if err = validator.Validate(); err != nil {
		log.Error("renew.Validate() err(%v)", err)
		return
	}
	// get renew panel
	rp, err := s.RenewPanel(c, mid)
	if err != nil {
		return
	}
	// add pay order
	payOrder := &model.PayOrder{
		Platform:    model.PlatformSystem,
		Mid:         mid,
		Status:      1,
		PaymentType: ui.PayChannelId,
		Ver:         1,
		OrderType:   model.PayOrderTypeSub,
	}
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	payOrder.CopyFromPanel(rp)
	if _, err = s.dao.TxInsertPayOrder(c, tx, payOrder); err != nil {
		return
	}
	s.mission(func() {
		// 允许云视听失败，job查询兜底
		s.dao.RenewYstOrder(c, nil)
	})
	return
}
