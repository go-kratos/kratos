package dao

import (
	"context"

	eleclient "go-common/app/service/main/vip/dao/ele-api-client"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// api name
const (
	_oauthGenerateAccessTokenURI = "/base.openservice/oauth_generate_access_token"
	_unionReceivePrizesURI       = "/member/bilibili/union/receive_prizes"
	_unionUpdateOpenIDURI        = "/member/bilibili/union/update_open_id"
	_bindUnionURI                = "/member/bilibili/bind_union"
	_canPurchaseURI              = "/member/bilibili/union/can_purchase"
	_getUnionMobileURI           = "/member/bilibili/union/get_union_mobile"
	_hongbaosURI                 = "/member/bilibili/union/hongbaos"
	_specailFoodsURI             = "/member/bilibili/union/special_foods"
)

// EleOauthGenerateAccessToken get access_token by auth_code.
func (d *Dao) EleOauthGenerateAccessToken(c context.Context, a *model.ArgEleAccessToken) (data *model.EleAccessTokenResp, err error) {
	args := new(struct {
		Request *model.ArgEleAccessToken `json:"request"`
	})
	args.Request = a
	resp := new(struct {
		Message string                    `json:"message"`
		Name    string                    `json:"name"`
		Data    *model.EleAccessTokenResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _oauthGenerateAccessTokenURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s) a(%+v) resp(%+v) error(%+v)", _oauthGenerateAccessTokenURI, a, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) a(%+v) resp(%+v) error(%+v)", _oauthGenerateAccessTokenURI, a, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	data = resp.Data
	log.Info("ele success api(%s) a(%+v) resp(%+v) data(%+v)", _oauthGenerateAccessTokenURI, a, resp, data)
	return
}

// EleUnionReceivePrizes union receive prizes.
func (d *Dao) EleUnionReceivePrizes(c context.Context, a *model.ArgEleReceivePrizes) (data []*model.EleReceivePrizesResp, err error) {
	args := new(struct {
		Request *model.ArgEleReceivePrizes `json:"request"`
	})
	args.Request = a
	resp := new(struct {
		Message string                        `json:"message"`
		Name    string                        `json:"name"`
		Data    []*model.EleReceivePrizesResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _unionReceivePrizesURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s) a(%+v) resp(%+v) error(%+v)", _unionReceivePrizesURI, a, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) a(%+v) resp(%+v) error(%+v)", _unionReceivePrizesURI, a, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	data = resp.Data
	log.Info("ele success api(%s) a(%+v) resp(%+v) data(%+v)", _unionReceivePrizesURI, a, resp, data)
	return
}

// EleUnionUpdateOpenID update_open_id req.
func (d *Dao) EleUnionUpdateOpenID(c context.Context, a *model.ArgEleUnionUpdateOpenID) (data *model.EleUnionUpdateOpenIDResp, err error) {
	args := new(struct {
		Request *model.ArgEleUnionUpdateOpenID `json:"request"`
	})
	args.Request = a
	resp := new(struct {
		Message string                          `json:"message"`
		Name    string                          `json:"name"`
		Data    *model.EleUnionUpdateOpenIDResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _unionUpdateOpenIDURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s) a(%+v) resp(%+v) error(%+v)", _unionUpdateOpenIDURI, a, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) a(%+v) resp(%+v) error(%+v)", _unionUpdateOpenIDURI, a, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	data = resp.Data
	log.Info("ele success api(%s) a(%+v) resp(%+v) data(%+v)", _unionUpdateOpenIDURI, a, resp, data)
	return
}

// EleBindUnion ele bind union salary vip.
func (d *Dao) EleBindUnion(c context.Context, a *model.ArgEleBindUnion) (data *model.EleBindUnionResp, err error) {
	args := new(struct {
		Request *model.ArgEleBindUnion `json:"request"`
	})
	args.Request = a
	resp := new(struct {
		Message string                  `json:"message"`
		Name    string                  `json:"name"`
		Data    *model.EleBindUnionResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _bindUnionURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s) a(%+v) resp(%+v) error(%+v)", _bindUnionURI, a, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) a(%+v) resp(%+v) error(%+v)", _bindUnionURI, a, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	data = resp.Data
	// 1.发放成功
	if data.Status != 1 && data.Status != 6 {
		log.Error("ele status faild api(%s) a(%+v) resp(%+v) data(%+v) error(%+v)", _bindUnionURI, a, resp, data, err)
		err = ecode.VipOrderEleVipGrantFaildErr
	}
	log.Info("ele success api(%s) a(%+v) resp(%+v) data(%+v)", _bindUnionURI, a, resp, data)
	return
}

// EleCanPurchase ele can purchase.
func (d *Dao) EleCanPurchase(c context.Context, a *model.ArgEleCanPurchase) (data *model.EleCanPurchaseResp, err error) {
	args := new(struct {
		Request *model.ArgEleCanPurchase `json:"request"`
	})
	args.Request = a
	resp := new(struct {
		Message string                    `json:"message"`
		Name    string                    `json:"name"`
		Data    *model.EleCanPurchaseResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _canPurchaseURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s) a(%+v) resp(%+v) error(%+v)", _canPurchaseURI, a, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	// 系统请求是否有误
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) a(%+v) resp(%+v) error(%+v)", _canPurchaseURI, a, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	// 业务逻辑是否有误
	data = resp.Data
	if data.Status != 1 {
		log.Error("ele status faild api(%s) a(%+v) resp(%+v) data(%+v) error(%+v)", _canPurchaseURI, a, resp, data, err)
		err = ecode.VipEleUnionBuyCanPurchaseErr
		return
	}
	log.Info("ele success api(%s) a(%+v) resp(%+v) data(%+v)", _canPurchaseURI, a, resp, data)
	return
}

// EleUnionMobile get ele union mobile.
func (d *Dao) EleUnionMobile(c context.Context, a *model.ArgEleUnionMobile) (data *model.EleUnionMobileResp, err error) {
	args := new(struct {
		Request *model.ArgEleUnionMobile `json:"request"`
	})
	args.Request = a
	resp := new(struct {
		Message string                    `json:"message"`
		Name    string                    `json:"name"`
		Data    *model.EleUnionMobileResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _getUnionMobileURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s) a(%+v) resp(%+v) error(%+v)", _getUnionMobileURI, a, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) a(%+v) resp(%+v) error(%+v)", _getUnionMobileURI, a, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	data = resp.Data
	log.Info("ele success api(%s) a(%+v) resp(%+v) data(%+v)", _getUnionMobileURI, a, resp, data)
	return
}

// EleRedPackages get ele red packages.
func (d *Dao) EleRedPackages(c context.Context) (data []*model.EleRedPackagesResp, err error) {
	args := new(struct{})
	resp := new(struct {
		Message string                      `json:"message"`
		Name    string                      `json:"name"`
		Data    []*model.EleRedPackagesResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _hongbaosURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s)resp(%+v) error(%+v)", _hongbaosURI, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) resp(%+v) error(%+v)", _hongbaosURI, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	data = resp.Data
	log.Info("ele success api(%s) resp(%+v) data(%+v)", _hongbaosURI, resp, data)
	return
}

// EleSpecailFoods get ele specail foods.
func (d *Dao) EleSpecailFoods(c context.Context) (data []*model.EleSpecailFoodsResp, err error) {
	args := new(struct{})
	resp := new(struct {
		Message string                       `json:"message"`
		Name    string                       `json:"name"`
		Data    []*model.EleSpecailFoodsResp `json:"data"`
	})
	err = d.eleclient.Post(c, d.c.Host.Ele, _specailFoodsURI, args, resp)
	if err != nil {
		log.Error("ele faild api(%s)resp(%+v) error(%+v)", _specailFoodsURI, resp, err)
		err = ecode.VipEleUnionReqErr
		return
	}
	if !eleclient.IsSuccess(resp.Message) {
		log.Error("ele message faild api(%s) resp(%+v) error(%+v)", _specailFoodsURI, resp, err)
		err = ecode.VipEleUnionRespCodeErr
		return
	}
	data = resp.Data
	log.Info("ele success api(%s) resp(%+v) data(%+v)", _specailFoodsURI, resp, data)
	return
}
