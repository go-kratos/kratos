package service

import (
	"context"
	"strconv"

	"go-common/app/service/main/passport-sns/api"
	"go-common/app/service/main/passport-sns/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var emptyCache = &model.SnsProto{}

// GetAuthorizeURL get sns authorize url
func (s *Service) GetAuthorizeURL(c context.Context, req *api.GetAuthorizeURLReq) (resp *api.GetAuthorizeURLReply, err error) {
	platform := parsePlatform(req.Platform)
	if platform == 0 || !s.isAppID(platform, req.AppId) {
		return nil, ecode.RequestErr
	}
	switch platform {
	case model.PlatformQQ:
		return &api.GetAuthorizeURLReply{Url: s.d.QQAuthorize(c, req.AppId, req.RedirectUrl, req.Display)}, nil
	case model.PlatformWEIBO:
		return &api.GetAuthorizeURLReply{Url: s.d.WeiboAuthorize(c, req.AppId, req.RedirectUrl, req.Display)}, nil
	}
	return
}

// Bind bind sns user
func (s *Service) Bind(c context.Context, req *api.BindReq) (resp *api.EmptyReply, err error) {
	platform := parsePlatform(req.Platform)
	if platform == 0 || !s.isAppID(platform, req.AppId) {
		return nil, ecode.RequestErr
	}
	var u *model.SnsUser
	u, err = s.d.SnsUserByMid(c, req.Mid, platform)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return nil, platformToMidBindErr(platform)
	}
	var info *model.Oauth2Info
	switch platform {
	case model.PlatformQQ:
		if info, err = s.d.QQOauth2Info(c, req.Code, req.RedirectUrl, s.AppMap[platform][req.AppId]); err != nil {
			return nil, err
		}
	case model.PlatformWEIBO:
		if info, err = s.d.WeiboOauth2Info(c, req.Code, req.RedirectUrl, s.AppMap[platform][req.AppId]); err != nil {
			return nil, err
		}
	}
	return &api.EmptyReply{}, s.bindSns(c, req.Mid, req.AppId, platform, info)
}

// Unbind unbind sns user
func (s *Service) Unbind(c context.Context, req *api.UnbindReq) (resp *api.EmptyReply, err error) {
	if req.Platform == model.PlatformAll {
		err = s.unbindAll(c, req.Mid)
		return &api.EmptyReply{}, err
	}
	platform := parsePlatform(req.Platform)
	if platform == 0 || !s.isAppID(platform, req.AppId) {
		return nil, ecode.RequestErr
	}
	return &api.EmptyReply{}, s.unbindSns(c, req.Mid, req.AppId, platform)
}

// GetInfo get info by mid
func (s *Service) GetInfo(c context.Context, req *api.GetInfoReq) (resp *api.GetInfoReply, err error) {
	addCache := true
	snsMap := make(map[string]*model.SnsProto)
	for _, platformStr := range s.PlatformStrList {
		var sns *model.SnsProto
		sns, err = s.d.SnsCache(c, req.Mid, platformStr)
		if err != nil {
			addCache = false
			err = nil
		}
		snsMap[platformStr] = sns
	}

	infos := make([]*api.Info, 0)
	hitAll := true
	for _, platformStr := range s.PlatformStrList {
		if snsMap[platformStr] == nil {
			hitAll = false
			break
		}
		if snsMap[platformStr].Mid != 0 {
			infos = append(infos, snsMap[platformStr].ConvertToInfo())
		}
	}
	if hitAll {
		return &api.GetInfoReply{Infos: infos}, nil
	}

	users, err := s.d.SnsUsers(c, req.Mid)
	if err != nil {
		return
	}
	for _, u := range users {
		snsMap[parsePlatformStr(u.Platform)] = u.ConvertToProto()
	}

	infos = make([]*api.Info, 0)
	for _, platformStr := range s.PlatformStrList {
		if snsMap[platformStr] != nil && snsMap[platformStr].Mid != 0 {
			infos = append(infos, snsMap[platformStr].ConvertToInfo())
		} else {
			snsMap[platformStr] = emptyCache
		}
	}

	if addCache {
		s.cache.Do(c, func(c context.Context) {
			for _, platformStr := range s.PlatformStrList {
				s.d.SetSnsCache(context.Background(), req.Mid, platformStr, snsMap[platformStr])
			}
		})
	}
	return &api.GetInfoReply{Infos: infos}, nil
}

// GetInfoByCode get info by code
func (s *Service) GetInfoByCode(c context.Context, req *api.GetInfoByCodeReq) (resp *api.GetInfoByCodeReply, err error) {
	platform := parsePlatform(req.Platform)
	if platform == 0 || !s.isAppID(platform, req.AppId) {
		return nil, ecode.RequestErr
	}

	var (
		info *model.Oauth2Info
		mid  int64
	)
	switch platform {
	case model.PlatformQQ:
		if info, err = s.d.QQOauth2Info(c, req.Code, req.RedirectUrl, s.AppMap[platform][req.AppId]); err != nil {
			return nil, err
		}
	case model.PlatformWEIBO:
		if info, err = s.d.WeiboOauth2Info(c, req.Code, req.RedirectUrl, s.AppMap[platform][req.AppId]); err != nil {
			return nil, err
		}
	}
	snsUser, err := s.d.SnsUserByUnionID(c, info.UnionID, platform)
	if err != nil {
		return nil, err
	}
	if snsUser != nil {
		mid = snsUser.Mid
	}

	proto := &model.Oauth2Proto{
		Mid:      mid,
		Platform: int32(platform),
		UnionID:  info.UnionID,
		OpenID:   info.OpenID,
		Token:    info.Token,
		Expires:  info.Expires,
		AppID:    req.AppId,
	}
	if err = s.d.SetOauth2Cache(c, info.OpenID, req.Platform, proto); err != nil {
		return nil, err
	}

	log.Info("GetInfoByCode request(%+v) response (%+v)", req, &api.GetInfoByCodeReply{
		Mid:     mid,
		UnionId: info.UnionID,
		OpenId:  info.OpenID,
		Expires: info.Expires,
		Token:   info.Token,
	})
	return &api.GetInfoByCodeReply{
		Mid:     mid,
		UnionId: info.UnionID,
		OpenId:  info.OpenID,
		Expires: info.Expires,
		Token:   info.Token,
	}, nil
}

// UpdateInfo update info
func (s *Service) UpdateInfo(c context.Context, req *api.UpdateInfoReq) (resp *api.EmptyReply, err error) {
	platform := parsePlatform(req.Platform)
	if platform == 0 || !s.isAppID(platform, req.AppId) {
		return nil, ecode.RequestErr
	}

	info, err := s.d.Oauth2Cache(c, req.OpenId, req.Platform)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, ecode.RequestErr
	}
	s.cache.Do(c, func(c context.Context) {
		s.d.DelOauth2Cache(context.Background(), req.OpenId, req.Platform)
	})
	if info.Mid == 0 {
		var u *model.SnsUser
		u, err = s.d.SnsUserByMid(c, req.Mid, platform)
		if err != nil {
			return nil, err
		}
		if u != nil {
			return nil, platformToMidBindErr(platform)
		}
		oauth2Info := &model.Oauth2Info{
			UnionID: info.UnionID,
			OpenID:  info.OpenID,
			Token:   info.Token,
			Expires: info.Expires,
		}
		return nil, s.bindSns(c, req.Mid, info.AppID, platform, oauth2Info)
	}

	tx, err := s.d.BeginTran(context.Background())
	if err != nil {
		log.Error("s.d.BeginTran error(%+v)", err)
		return
	}

	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
			return
		}
		snsLog := &model.SnsLog{
			Mid:      info.Mid,
			OpenID:   info.OpenID,
			UnionID:  info.UnionID,
			AppID:    info.AppID,
			Platform: platform,
			Operator: "",
			Operate:  model.OperateLogin,
		}
		s.job.Do(c, func(c context.Context) {
			s.sendSnsLog(context.Background(), snsLog)
		})
		s.cache.Do(c, func(c context.Context) {
			s.d.DelSnsCache(context.Background(), info.Mid, parsePlatformStr(platform))
		})
	}()

	snsUser := &model.SnsUser{
		Mid:      info.Mid,
		Platform: platform,
		Expires:  info.Expires,
	}
	if _, err = s.d.TxUpdateSnsUser(tx, snsUser); err != nil {
		return
	}
	snsToken := &model.SnsToken{
		Mid:      info.Mid,
		Platform: platform,
		Token:    info.Token,
		Expires:  info.Expires,
	}
	if _, err = s.d.TxUpdateSnsToken(tx, snsToken); err != nil {
		return
	}
	return
}

func (s *Service) bindSns(c context.Context, mid int64, appID string, platform int, info *model.Oauth2Info) (err error) {
	snsUser, err := s.d.SnsUserByUnionID(c, info.UnionID, platform)
	if err != nil {
		return err
	}
	if snsUser != nil {
		return platformToSnsBindErr(platform)
	}

	tx, err := s.d.BeginTran(c)
	if err != nil {
		log.Error("s.d.BeginTran error(%+v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
	}()

	snsUser = &model.SnsUser{
		Mid:      mid,
		UnionID:  info.UnionID,
		Platform: platform,
		Expires:  info.Expires,
	}
	affected, err := s.d.TxAddSnsUser(tx, snsUser)
	if err != nil {
		return
	}
	if affected > 0 {
		snsOpenID := &model.SnsOpenID{
			Mid:      mid,
			OpenID:   info.OpenID,
			UnionID:  info.UnionID,
			AppID:    appID,
			Platform: platform,
		}
		if _, err = s.d.TxAddSnsOpenID(tx, snsOpenID); err != nil {
			return
		}

		snsToken := &model.SnsToken{
			Mid:      mid,
			OpenID:   info.OpenID,
			UnionID:  info.UnionID,
			Platform: platform,
			Token:    info.Token,
			Expires:  info.Expires,
			AppID:    appID,
		}
		if _, err = s.d.TxAddSnsToken(tx, snsToken); err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.d.DelSnsCache(context.Background(), mid, parsePlatformStr(platform))
	})
	if affected > 0 {
		s.job.Do(c, func(c context.Context) {
			snsLog := &model.SnsLog{
				Mid:      mid,
				OpenID:   info.OpenID,
				UnionID:  info.UnionID,
				AppID:    appID,
				Platform: platform,
				Operator: "",
				Operate:  model.OperateBind,
			}
			s.sendSnsLog(context.Background(), snsLog)
		})
	}
	return
}

func (s *Service) unbindSns(c context.Context, mid int64, appID string, platform int) (err error) {
	u, err := s.d.SnsUserByMid(c, mid, platform)
	if err != nil {
		return
	}
	affected, err := s.d.DelSnsUser(c, mid, platform)
	if err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.d.DelSnsCache(context.Background(), mid, parsePlatformStr(platform))
	})
	if affected > 0 {
		snsLog := &model.SnsLog{
			Mid: mid,
			//OpenID:   qq.OpenID,
			UnionID:  u.UnionID,
			AppID:    appID,
			Platform: platform,
			Operator: "",
			Operate:  model.OperateUnbind,
		}
		s.job.Do(c, func(c context.Context) {
			s.sendSnsLog(context.Background(), snsLog)
		})
	}
	return
}

func (s *Service) unbindAll(c context.Context, mid int64) (err error) {
	affected, err := s.d.DelSnsUsers(c, mid)
	if err != nil {
		return
	}
	s.cache.Do(c, func(ctx context.Context) {
		for _, platformStr := range s.PlatformStrList {
			s.d.DelSnsCache(context.Background(), mid, platformStr)
		}
	})

	if affected > 0 {
		snsLog := &model.SnsLog{
			Mid:         mid,
			Operator:    "admin",
			Operate:     model.OperateDelete,
			Description: "删除用户",
		}
		s.job.Do(c, func(c context.Context) {
			s.sendSnsLog(context.Background(), snsLog)
		})
	}
	return
}

func (s *Service) sendSnsLog(c context.Context, snsLog *model.SnsLog) {
	for i := 0; i < 3; i++ {
		if err := s.snsLogPub.Send(c, strconv.FormatInt(snsLog.Mid, 10), snsLog); err != nil {
			log.Error("fail to send snsLog(%+v) error(%+v)", snsLog, err)
			continue
		}
		break
	}
}
