package service

import (
	"context"
	"time"

	"go-common/app/service/main/ugcpay/model"
	"go-common/library/log"
)

// IncomeUserAssetOverview .
func (s *Service) IncomeUserAssetOverview(ctx context.Context, mid int64, currency string) (user *model.AggrIncomeUser, monthReady int64, newDailyBill *model.Bill, err error) {
	if user, err = s.dao.RawAggrIncomeUser(ctx, mid, currency); err != nil {
		return
	}
	var (
		dayVer   = s.dayVer(time.Now().Add(-time.Hour * 48))
		monthVer = s.monthVer(time.Now())
		billList []*model.Bill
		biz      = "asset"
	)
	log.Info("IncomeUserAssetOverview mid: %d, biz: %s, currency: %s, month_ver: %d", mid, biz, currency, monthVer)

	if billList, err = s.dao.RawBillUserDailyByMonthVer(ctx, mid, biz, currency, monthVer); err != nil {
		return
	}
	for _, b := range billList {
		monthReady += b.In - b.Out
	}
	if newDailyBill, err = s.dao.RawBillUserDaily(ctx, mid, biz, currency, dayVer); err != nil {
		return
	}
	return
}

func (s *Service) dayVer(t time.Time) (ver int64) {
	return int64(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

func (s *Service) monthVer(t time.Time) (ver int64) {
	return int64(t.Year()*100 + int(t.Month()))
}

// IncomeUserAssetList .
func (s *Service) IncomeUserAssetList(ctx context.Context, mid int64, currency string, ver int64, pn, ps int64) (list *model.AggrIncomeUserAssetList, page *model.Page, err error) {
	var (
		assets []*model.AggrIncomeUserAsset
		limit  = 1000
	)
	if assets, err = s.dao.RawAggrIncomeUserAssetList(ctx, mid, currency, ver, limit); err != nil {
		return
	}
	l, page := s.pageIncomeUseAsset(assets, pn, ps)
	list = &model.AggrIncomeUserAssetList{
		MID:    mid,
		Ver:    ver,
		Assets: l,
		Page:   page,
	}
	return
}

// IncomeUserAsset .
func (s *Service) IncomeUserAsset(ctx context.Context, mid int64, oid int64, otype string, currency string, ver int64) (asset *model.AggrIncomeUserAsset, err error) {
	return s.dao.RawAggrIncomeUserAsset(ctx, mid, currency, oid, otype, ver)
}

func (s *Service) pageIncomeUseAsset(list []*model.AggrIncomeUserAsset, pn, ps int64) (l []*model.AggrIncomeUserAsset, page *model.Page) {
	page = &model.Page{
		Num:   pn,
		Size:  ps,
		Total: int64(len(list)),
	}
	l = make([]*model.AggrIncomeUserAsset, 0)
	from := (pn - 1) * ps
	to := pn * ps
	if page.Total < from {
		return
	}
	if page.Total < to {
		l = append(l, list[from:]...)
	} else {
		l = append(l, list[from:to]...)
	}
	return
}
