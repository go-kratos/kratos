package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_mailCreateCouponCodeURI = "/mall-marketing/coupon_code/create"
)

// MailCouponCodeCreate mail coupon code create.
func (d *Dao) MailCouponCodeCreate(c context.Context, a *model.ArgMailCouponCodeCreate) (err error) {
	resp := new(struct {
		Code    int                             `json:"code"`
		Message string                          `json:"message"`
		Data    *model.MailCouponCodeCreateResp `json:"data"`
	})
	err = d.mailclient.Post(c, d.c.Host.Mail+_mailCreateCouponCodeURI, a, resp)
	if err != nil {
		log.Error("mail faild api(%s) args(%+v) res(%+v) error(%+v)", _mailCreateCouponCodeURI, a, resp, err)
		err = ecode.VipMailReqErr
		return
	}
	if resp.Code != ecode.OK.Code() {
		log.Error("mail code faild api(%s) args(%+v) res(%+v) error(%+v)", _mailCreateCouponCodeURI, a, resp, err)
		err = ecode.VipMailRespCodeErr
		return
	}
	log.Info("mail success api(%s) args(%+v) res(%+v) data(%+v)", _mailCreateCouponCodeURI, a, resp, resp.Data)
	return
}
