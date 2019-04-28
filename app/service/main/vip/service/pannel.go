package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

//PannelInfoNew pannelInfo new.
func (s *Service) PannelInfoNew(c context.Context, mid int64, arg *model.ArgPannel) (pi *model.PannelInfo, err error) {
	var (
		vipInfo         *model.VipInfoResp
		discountHistory *model.VipUserDiscountHistory
	)
	if vipInfo, err = s.ByMid(c, mid); err != nil {
		err = errors.WithStack(err)
		return
	}
	if discountHistory, err = s.dao.DiscountSQL(c, mid, model.FirstDiscountBuyVip); err != nil {
		err = errors.WithStack(err)
		return
	}
	if pi, err = s.pannelInfoWithCache(c, mid, arg, metadata.String(c, metadata.RemoteIP)); err != nil {
		err = errors.WithStack(err)
		return
	}
	if vipInfo.PayType == 1 {
		vipMonths := make([]*model.VipMonthsPriceBo, 0)
		for _, v := range pi.VipMonths {
			if v.OrderType == 1 {
				continue
			}
			vipMonths = append(vipMonths, v)
		}
		pi.VipMonths = vipMonths
	} else if discountHistory == nil {
		for _, v := range pi.VipMonths {
			if v.OrderType == 1 {
				v.Price = v.FirstDiscountMoney
				v.DiscountRate = s.calcDiscountRate(v.OriginalPrice, v.Price) + "折"
			}
		}
	}
	return
}

func (s *Service) pannelInfoWithCache(c context.Context, mid int64, arg *model.ArgPannel, realIP string) (pi *model.PannelInfo, err error) {
	var (
		payTypes       []*model.PayTypeBo
		vipMonthPannel []*model.VipMonthsPriceBo
	)
	pi = new(model.PannelInfo)
	payTypes = s.payType(model.PlatformByName[arg.Platform], realIP)
	pi.PayTypes = payTypes
	if vipMonthPannel, err = s.findVipMonthPanel(c, arg, model.PlatformByName[arg.Platform]); err != nil {
		log.Error("%+v", err)
		return
	}
	pi.VipMonths = vipMonthPannel
	pi.BcoinTips = strconv.Itoa(s.c.Property.AnnualVipBcoinCouponMoney)
	return
}

func (s *Service) calcDiscountRate(originPrice, price float64) string {
	divide := price / originPrice
	mul := divide * 10
	return fmt.Sprintf("%.2f", mul)
}

func (s *Service) payType(plat int, realIP string) (payType []*model.PayTypeBo) {

	switch plat {
	case model.DevicePC:
		payType = s.payChannel(realIP, model.ALIPAY, model.WECHAT, model.BCION, model.BANK)
	case model.DeviceANDROID:
		payType = s.payChannelFromConfig("android")
	case model.DeviceIOS:
		payType = s.payChannelFromConfig("ios")
	case model.DeviceIPAD:
		payType = s.payChannelFromConfig("ios")
	default:
		payType = s.payChannel(realIP, model.ALIPAY, model.WECHAT, model.BCION, model.BANK)
	}
	return
}

func (s *Service) payChannelFromConfig(key string) (payTypeBos []*model.PayTypeBo) {
	val := s.c.Property.PayType[key]
	payTypes := strings.Split(val, ",")
	for _, v := range payTypes {
		payTypeNum := model.PayTypeName[v]
		r := new(model.PayTypeBo)
		r.Name = model.PayWayName[payTypeNum]
		r.Code = v
		payTypeBos = append(payTypeBos, r)
	}
	return
}

func (s *Service) payChannel(clientIP string, payTypes ...int8) (payTypeBo []*model.PayTypeBo) {
	for _, v := range payTypes {
		r := new(model.PayTypeBo)
		r.Name = model.PayWayName[v]
		r.Code = model.PayType[v]
		if model.BANK == v {
			banks := make([]*model.PayBankResp, 0)
			if err := s.dao.PayBanks(context.TODO(), clientIP, banks); err != nil {
				return
			}
			payBanks := make([]*model.PayBankBo, 0)
			for _, v := range banks {
				bank := new(model.PayBankBo)
				bank.Code = v.BanckCode
				bank.Name = v.Name
				bank.Image = v.Res
				payBanks = append(payBanks, bank)
			}
			r.Banks = payBanks
		}
		payTypeBo = append(payTypeBo, r)
	}
	return
}

func (s *Service) findVipMonthPanel(c context.Context, arg *model.ArgPannel, plat int) (vipMonthPanels []*model.VipMonthsPriceBo, err error) {
	var (
		vipMonths  []*model.Month
		monthPrice *model.PriceMapping
	)

	if vipMonths, err = s.dao.AllMonthByOrder(c, "desc"); err != nil {
		log.Error("%+v", err)
		return
	}
	//vipMonthsTemp := make([]*model.Month, 0)
	//for _, v := range vipMonths {
	//	vipMonthsTemp = append(vipMonthsTemp, v)
	//}
	for _, v := range vipMonths {
		if monthPrice, err = s.dao.PriceMapping(c, v.ID, model.PayPlatform[plat]); err != nil {
			log.Error("%+v", err)
			return
		}
		if monthPrice == nil {
			continue
		}
		bo := new(model.VipMonthsPriceBo)
		bo.ID = monthPrice.ID
		bo.OriginalPrice = monthPrice.Money
		bo.Remark = monthPrice.Remark
		bo.OrderType = v.MonthType
		bo.MonthType = monthPrice.MonthType
		bo.Selected = monthPrice.Selected
		useDiscountMoney := s.checkIsDiscount(monthPrice.StartTime.Time(), monthPrice.EndTime.Time(), monthPrice.DiscountMoney)
		if useDiscountMoney {
			bo.Price = monthPrice.DiscountMoney
		} else {
			bo.Price = monthPrice.Money
		}

		if v.Month == 1 {
			bo.MonthStr = "月度大会员"
		} else if v.Month == 3 {
			bo.MonthStr = "季度大会员"
		} else if v.Month == 12 {
			bo.MonthStr = "年度大会员"
		}

		if v.MonthType == 1 {
			bo.FirstDiscountMoney = monthPrice.FirstDiscountMoney
			if v.Month == 1 {
				bo.MonthStr = "连续包月大会员"
			}
		}
		bo.DiscountRate = s.calcDiscountRate(bo.OriginalPrice, bo.Price) + "折"
		vipMonthPanels = append(vipMonthPanels, bo)
	}

	return
}

func (s *Service) checkIsDiscount(start, end time.Time, discountPrice float64) bool {
	if discountPrice == 0 {
		return false
	}
	now := time.Now()
	if start.Before(now) && end.After(now) {
		return true
	}
	return false
}
