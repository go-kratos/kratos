package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// RegisterOpenID register open_id.
func (s *Service) RegisterOpenID(c context.Context, a *model.ArgRegisterOpenID) (res *model.RegisterOpenIDResp, err error) {
	openID := createOpenID(a.Mid, a.AppID)
	if err = s.dao.AddOpenInfo(c, &model.OpenInfo{
		Mid:    a.Mid,
		OpenID: openID,
		AppID:  a.AppID,
	}); err != nil {
		return
	}
	res = &model.RegisterOpenIDResp{
		OpenID: openID,
	}
	s.dao.DelOpenInfoCache(c, openID, a.AppID)
	return
}

func createOpenID(mid int64, appID int64) string {
	mh := md5.Sum([]byte(fmt.Sprintf("%d_%d", mid, appID)))
	return hex.EncodeToString(mh[:])
}

// OpenAuthCallBack third open call back.
func (s *Service) OpenAuthCallBack(c context.Context, a *model.ArgOpenAuthCallBack) (err error) {
	var data *model.EleAccessTokenResp
	if data, err = s.dao.EleOauthGenerateAccessToken(c, &model.ArgEleAccessToken{AuthCode: a.ThirdCode}); err != nil {
		return
	}
	// add mid to open info.
	if _, err = s.RegisterOpenID(c, &model.ArgRegisterOpenID{
		Mid:   a.Mid,
		AppID: a.AppID,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = s.OpenBindByMid(c, &model.ArgOpenBindByMid{
		Mid:       a.Mid,
		AppID:     a.AppID,
		OutOpenID: data.OpenID,
	}); err != nil {
		return
	}
	return
}

// UserInfoByOpenID get userinfo by open_id.
func (s *Service) UserInfoByOpenID(c context.Context, a *model.ArgUserInfoByOpenID) (
	res *model.UserInfoByOpenIDResp, err error) {
	var (
		oi *model.OpenInfo
		bi *memmdl.BaseInfo
	)
	// check open_id
	if oi, err = s.dao.OpenInfoByOpenID(c, a.OpenID, a.AppID); err != nil {
		return
	}
	if oi == nil {
		return nil, ecode.VipAssociateOpenIDNotExsitErr
	}
	res = new(model.UserInfoByOpenIDResp)
	if bi, err = s.memRPC.Base(c, &memmdl.ArgMemberMid{
		Mid:      oi.Mid,
		RemoteIP: a.IP,
	}); err != nil {
		log.Error("s.memRPC.Base(%d) err(%+v)", oi.Mid, err)
	}
	res.Name = strconv.FormatInt(oi.Mid, 10)
	if bi != nil && bi.Name != "" {
		res.Name = bi.Name
	}
	// if ob, err = s.dao.BindInfoByMid(c, oi.Mid, a.AppID); err != nil {
	// 	return
	// }
	// if ob != nil {
	// 	res.BindState = 1 //had bind
	// 	res.OutOpenID = ob.OutOpenID
	// }
	return
}
