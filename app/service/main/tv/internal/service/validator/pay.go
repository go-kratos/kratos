package validator

import (
	"go-common/app/service/main/tv/internal/model"
	"go-common/app/service/main/tv/internal/pkg"
	"go-common/library/ecode"
	"go-common/library/log"
	"time"
)

type PayParamValidator struct {
	PayParam *model.PayParam
}

func (pv *PayParamValidator) Validate() error {
	if pv.PayParam == nil {
		return ecode.TVIPTokenErr
	}
	if pv.PayParam.IsExpired() {
		return ecode.TVIPTokenExpire
	}
	return nil
}

type PanelValidator struct {
	Panel  *model.PanelPriceConfig
	MVip   *model.MainVip
	BuyNum int32
	Ui     *model.UserInfo
}

func (pv *PanelValidator) isSuitable() bool {
	switch pv.Panel.SuitType {
	case model.SuitTypeAll:
		return true
	case model.SuitTypeMvip:
		// NOTE: tv 会员无法通过主站会员升级方式购买
		if pv.Ui != nil && pv.Ui.IsVip() {
			return false
		}
		return pv.MVip.IsVip() && pv.MVip.Months() >= 1
	default:
		log.Error("pv.isSuitable err(%s)", "SuiteTypeUnknown")
		return false
	}
}

func (pv *PanelValidator) isBuyNumExceeded() bool {
	if pv.Panel.SubType != model.SuitTypeMvip {
		return false
	}
	return pv.BuyNum > pv.MVip.Months()
}

func (pv *PanelValidator) Validate() error {
	if pv.Panel == nil {
		return ecode.TVIPPanelNotFound
	}
	if !pv.isSuitable() {
		return ecode.TVIPPanelNotSuitalbe
	}
	if pv.isBuyNumExceeded() {
		return ecode.TVIPBuyNumExceeded
	}
	return nil
}

type CreateOrderValidator struct {
	PayOrder          *model.PayOrder
	PaymentType       string
	PayParam          *model.PayParam
	Panel             *model.PanelPriceConfig
	MVip              *model.MainVip
	Ui                *model.UserInfo
	PayingOrders      []*model.PayOrder
	MVipPanels        []*model.PanelPriceConfig
	MVipRateMaxNumber int
}

func (cv *CreateOrderValidator) isAliContract() bool {
	return cv.Panel.SubType == model.SubTypeContract && cv.PaymentType == model.PaymentTypeAliPay
}

func (cv *CreateOrderValidator) isDupOrderNo() bool {
	return cv.PayOrder != nil && cv.PayOrder.OrderNo == cv.PayParam.OrderNo
}

func (cv *CreateOrderValidator) isRateExceed() bool {
	return len(cv.PayingOrders) >= cv.MVipRateMaxNumber
}

func (cv *CreateOrderValidator) isMVipExceed() bool {
	if cv.Panel.SuitType != model.SuitTypeMvip {
		return false
	}
	for _, order := range cv.PayingOrders {
		productId := order.ProductId
		for _, panel := range cv.MVipPanels {
			if panel.ProductId == productId {
				return true
			}
		}
	}
	return false
}

func (cv *CreateOrderValidator) Validate() error {
	if cv.isRateExceed() {
		return ecode.TVIPBuyRateExceeded
	}
	if cv.isMVipExceed() {
		return ecode.TVIPMVipRateExceeded
	}
	if cv.isAliContract() {
		return ecode.TVIPPanelNotSuitalbe
	}
	if cv.isDupOrderNo() {
		return ecode.TVIPDupOrderNo
	}
	paramValidator := &PayParamValidator{PayParam: cv.PayParam}
	if err := paramValidator.Validate(); err != nil {
		return err
	}
	panelValidator := &PanelValidator{Panel: cv.Panel, MVip: cv.MVip, Ui: cv.Ui}
	return panelValidator.Validate()
}

type RenewVipValidator struct {
	UserInfo     *model.UserInfo
	FromDuration string
	ToDuration   string
}

func (rv *RenewVipValidator) isTooEarly() bool {
	d, err := time.ParseDuration(rv.FromDuration)
	if err != nil {
		log.Error("rv.ParseDuration(%s) err(%v)", rv.FromDuration, err)
		d = time.Hour * 24
	}
	return int64(rv.UserInfo.OverdueTime)-time.Now().Unix() > int64(d.Seconds())
}

func (rv *RenewVipValidator) isTooLate() bool {
	d, err := time.ParseDuration(rv.ToDuration)
	if err != nil {
		log.Error("rv.ParseDuration(%s) err(%v)", rv.ToDuration, err)
		d = time.Hour * 24
	}
	return time.Now().Unix()-int64(rv.UserInfo.OverdueTime) > int64(d.Seconds())
}

func (rv *RenewVipValidator) isContracted() bool {
	return rv.UserInfo.IsContracted()
}

func (rv *RenewVipValidator) Validate() error {
	if rv.isTooEarly() {
		return ecode.TVIPRenewTooEarly
	}
	if rv.isTooLate() {
		return ecode.TVIPRenewTooLate
	}
	if !rv.isContracted() {
		return ecode.TVIPNotContracted
	}
	return nil
}

type PayCallbackValidator struct {
	Signer         *pkg.Signer
	CallbackReq    *model.YstPayCallbackReq
	PayOrder       *model.PayOrder
	ExpireDuration string
}

func (pv *PayCallbackValidator) isSignValid() bool {
	sign, err := pv.Signer.Sign(pv.CallbackReq)
	if err != nil {
		return false
	}
	return sign == pv.CallbackReq.Sign
}

func (pv *PayCallbackValidator) hasPayingOrder() bool {
	if pv.PayOrder == nil {
		return false
	}
	if pv.PayOrder.Status != model.PayOrderStatusPaying {
		return false
	}
	return true
}

func (pv *PayCallbackValidator) isOrderExpired() bool {
	return false
	//d, err := time.ParseDuration(pv.ExpireDuration)
	//if err != nil {
	//	d = time.Hour * 2
	//}
	//return time.Now().Unix()-int64(pv.PayOrder.Ctime) > int64(d.Seconds())
}

func (pv *PayCallbackValidator) Validate() error {
	if !pv.isSignValid() {
		return ecode.TVIPSignErr
	}
	if !pv.hasPayingOrder() {
		return ecode.TVIPPayOrderNotFound
	}
	if pv.isOrderExpired() {
		return ecode.TVIPPayOrderExpired
	}
	return nil
}
