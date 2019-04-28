package vip

import (
	"bytes"
	"context"

	"go-common/app/interface/main/account/model"
	v1 "go-common/app/service/main/vip/api"
	"go-common/library/log"
)

// ElemeOAuthURI get ele oauth uri.
func (s *Service) ElemeOAuthURI(c context.Context, csrf string) string {
	var buf bytes.Buffer
	buf.WriteString(s.c.Vipproperty.EleOAuthURI)
	buf.WriteString("?consumer_key=")
	buf.WriteString(s.c.Vipproperty.EleConsumerKey)
	buf.WriteString("&scope=user_basic_info&redirect_uri=")
	buf.WriteString(s.c.Vipproperty.EleOAuthCallBackURI)
	buf.WriteString("&state=")
	buf.WriteString(csrf)
	return buf.String()
}

// OpenIDByAuthCode third get open_id by auth code.
func (s *Service) OpenIDByAuthCode(c context.Context, a *model.ArgAuthCode) (res *model.OpenIDResp, err error) {
	var (
		data *model.OAuth2InfoResp
		r    *v1.RegisterOpenIDReply
	)
	//oauth2 token.
	if data, err = s.vipDao.OAuth2ByCode(c, a); err != nil {
		return
	}
	res = new(model.OpenIDResp)
	if r, err = s.vipgRPC.RegisterOpenID(c, &v1.RegisterOpenIDReq{Mid: data.Mid, AppId: a.APPID}); err != nil {
		return
	}
	res.OpenID = r.OpenId
	return
}

// OpenAuthCallBack open auth callback[third->bilibili].
func (s *Service) OpenAuthCallBack(c context.Context, a *model.ArgOpenAuthCallBack) (uri string) {
	var (
		ret = "0"
		err error
	)
	if _, err = s.vipgRPC.OpenAuthCallBack(c, &v1.OpenAuthCallBackReq{
		Mid:       a.Mid,
		ThirdCode: a.ThirdCode,
		AppId:     a.AppID,
	}); err != nil {
		log.Error("vipSvc.OpenAuthCallBack(%+v) err(%+v)", a, err)
		ret = "1"
	}
	uri = s.c.Vipproperty.ActivityURI + "?bind_ret=" + ret
	return
}

// BilibiliPrizeGrant vip prize grant.
func (s *Service) BilibiliPrizeGrant(c context.Context, a *model.ArgBilibiliPrizeGrant) (res *v1.BilibiliPrizeGrantReply, err error) {
	return s.vipgRPC.BilibiliPrizeGrant(c, &v1.BilibiliPrizeGrantReq{
		PrizeKey: a.PrizeKey,
		UniqueNo: a.UniqueNo,
		OpenId:   a.OpenID,
		AppId:    a.AppID,
	})
}

// BilibiliVipGrant vip grant.
func (s *Service) BilibiliVipGrant(c context.Context, a *model.ArgBilibiliVipGrant) (err error) {
	_, err = s.vipgRPC.BilibiliVipGrant(c, &v1.BilibiliVipGrantReq{
		OpenId:     a.OpenID,
		AppId:      a.AppID,
		OutOpenId:  a.OutOpenID,
		OutOrderNo: a.OutOrderNO,
		Duration:   a.Duration,
	})
	return
}

// OpenBindByOutOpenID associate user bind by out_open_id [third -> bilibili].
func (s *Service) OpenBindByOutOpenID(c context.Context, a *model.ArgBind) (err error) {
	_, err = s.vipgRPC.OpenBindByOutOpenID(c, &v1.OpenBindByOutOpenIDReq{
		OpenId:    a.OpenID,
		OutOpenId: a.OutOpenID,
		AppId:     a.AppID,
	})
	return
}

// UserInfoByOpenID get userinfo by open_id.
func (s *Service) UserInfoByOpenID(c context.Context, a *model.ArgUserInfoByOpenID) (res *v1.UserInfoByOpenIDReply, err error) {
	return s.vipgRPC.UserInfoByOpenID(c, &v1.UserInfoByOpenIDReq{
		OpenId: a.OpenID,
		Ip:     a.IP,
		AppId:  a.AppID,
	})
}
