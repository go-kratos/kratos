package http

import (
	"strings"

	"go-common/app/interface/main/account/model"
	col "go-common/app/service/main/coupon/model"
	vipmol "go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_headerBuvid = "Buvid"
)

func codeVerify(c *bm.Context) {
	c.JSON(vipSvc.CodeVerify(c))
}

func codeOpen(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(struct {
		Token  string `form:"token" validate:"required"`
		Code   string `form:"code" validate:"required"`
		Verify string `form:"verify" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.Code = strings.Trim(arg.Code, " ")
	c.JSON(vipSvc.CodeOpen(c, mid.(int64), arg.Code, arg.Token, arg.Verify))
}

// tips info.
func tips(c *bm.Context) {
	var (
		res *vipmol.TipsResp
		arg = new(model.TipsReq)
		err error
	)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.Tips(c, arg); err != nil {
		log.Error("vipSvc.Tips(%+v) err(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

// tips info.
func tipsv2(c *bm.Context) {
	var (
		res []*vipmol.TipsResp
		arg = new(model.TipsReq)
		err error
	)
	if err = c.Bind(arg); err != nil {
		log.Error("c.Bind err(%+v)", err)
		return
	}
	if res, err = vipSvc.TipsV2(c, arg); err != nil {
		log.Error("vipSvc.Tips(%+v) err(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func vipPanel(c *bm.Context) {
	var (
		err error
		res *vipmol.VipPirceResp5
	)
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(model.VipPanelRes)
	if err = c.Bind(arg); err != nil {
		return
	}
	if res, err = vipSvc.VipPanel5(c, mid.(int64), arg); err != nil {
		log.Error("vipSvc.VipPanel(%+v) err(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func couponUsable(c *bm.Context) {
	var (
		err error
		res *col.CouponAllowancePanelInfo
	)
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(model.ArgVipCoupon)
	if err = c.Bind(arg); err != nil {
		return
	}
	if res, err = vipSvc.CouponBySuitID(c, mid.(int64), arg.ID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func couponList(c *bm.Context) {
	var (
		err error
		res *col.CouponAllowancePanelResp
	)
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(model.ArgVipCoupon)
	if err = c.Bind(arg); err != nil {
		return
	}
	if res, err = vipSvc.CouponsForPanelV2(c, mid.(int64), arg.ID); err != nil {
		log.Error("vipSvc.CouponsForPanelV2(%+v) err(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func couponUnlock(c *bm.Context) {
	var err error
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(model.ArgVipCancelPay)
	if err = c.Bind(arg); err != nil {
		return
	}
	if err = vipSvc.CancelUseCoupon(c, &vipmol.ArgCancelUseCoupon{
		Mid:         mid.(int64),
		CouponToken: arg.CouponToken,
	}); err != nil {
		log.Error("vipSvc.CancelUseCoupon(%+v) err(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(model.CouponCancelExplain, nil)
}

func vipPanelV2(c *bm.Context) {
	var err error
	arg := new(model.ArgVipPanel)
	if err = c.Bind(arg); err != nil {
		return
	}
	mid, exists := c.Get("mid")
	if exists {
		arg.Mid = mid.(int64)
	}
	arg.IP = metadata.String(c, metadata.RemoteIP)
	c.JSON(vipSvc.VipPanelV2(c, arg))
}

func vipPanelV8(c *bm.Context) {
	var err error
	arg := new(model.ArgVipPanel)
	if err = c.Bind(arg); err != nil {
		return
	}
	mid, exists := c.Get("mid")
	if exists {
		arg.Mid = mid.(int64)
	}
	arg.IP = metadata.String(c, metadata.RemoteIP)
	c.JSON(vipSvc.VipPanelV8(c, arg))
}

func privilegeBySid(c *bm.Context) {
	var err error
	arg := new(vipmol.ArgPrivilegeBySid)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.PrivilegebySid(c, arg))
}

func privilegeByType(c *bm.Context) {
	var err error
	arg := new(vipmol.ArgPrivilegeDetail)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.PrivilegebyType(c, arg))
}

func vipManagerInfo(c *bm.Context) {
	c.JSON(vipSvc.ManagerInfo(c))
}

func codeOpeneds(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.CodeInfoReq)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.CodeOpeneds(c, arg, metadata.String(c, metadata.RemoteIP)))
}

func unfrozen(c *bm.Context) {

	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}

	c.JSON(nil, vipSvc.Unfrozen(c, mid.(int64)))
}

func frozenTime(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.AccountNotLogin)
		return
	}
	c.JSON(vipSvc.FrozenTime(c, mid.(int64)))
}

func publicPriceList(c *bm.Context) {
	var (
		err error
		res *vipmol.VipPirceResp
		mid int64
	)
	midStr, exists := c.Get("mid")
	if exists {
		mid = midStr.(int64)
	}
	arg := new(model.VipPanelRes)
	if err = c.Bind(arg); err != nil {
		return
	}
	if res, err = vipSvc.VipPanel(c, mid, arg); err != nil {
		log.Error("vipSvc.VipPanel(%+v) err(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}

func useBatch(c *bm.Context) {
	var err error
	arg := new(vipmol.ArgUseBatch)
	if err = c.Bind(arg); err != nil {
		log.Error("use batch bind err(%+v) arg(%+v)", err, arg)
		return
	}
	c.JSON(nil, vipSvc.UseBatch(c, arg))
}

func orderStatus(c *bm.Context) {
	var (
		err error
	)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(vipmol.ArgDialog)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.Mid = midI.(int64)
	c.JSON(vipSvc.OrderStatus(c, arg))
}

func resourceBanner(c *bm.Context) {
	var (
		err error
	)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(model.ArgResource)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.MID = midI.(int64)
	arg.Buvid = c.Request.Header.Get(_headerBuvid)
	arg.Plat = model.Plat(arg.MobiApp, arg.Device)
	c.JSON(vipSvc.ResourceBanner(c, arg))
}

func resourceBuy(c *bm.Context) {
	var (
		err error
	)
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(model.ArgResource)
	if err = c.Bind(arg); err != nil {
		return
	}
	arg.MID = midI.(int64)
	arg.Buvid = c.Request.Header.Get(_headerBuvid)
	arg.Plat = model.Plat(arg.MobiApp, arg.Device)
	c.JSON(vipSvc.ResourceBuy(c, arg))
}
func couponBySuitIDV2(c *bm.Context) {
	var err error
	arg := new(model.ArgCouponBySuitID)
	if err = c.Bind(arg); err != nil {
		return
	}
	mid, exists := c.Get("mid")
	if exists {
		arg.Mid = mid.(int64)
	}
	c.JSON(vipSvc.CouponBySuitIDV2(c, arg))
}

func vipPanelV9(c *bm.Context) {
	var err error
	arg := new(model.ArgVipPanel)
	if err = c.Bind(arg); err != nil {
		return
	}
	mid, exists := c.Get("mid")
	if exists {
		arg.Mid = mid.(int64)
	}
	arg.IP = metadata.String(c, metadata.RemoteIP)
	c.JSON(vipSvc.VipPanelV9(c, arg))
}
