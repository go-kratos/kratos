package service

import (
	"context"
	"time"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Oauth oauth.
func (s *Service) Oauth(c context.Context, app *model.App, accessToken, from string) (token *model.Token, err error) {
	if accessToken == "" {
		err = ecode.AccessKeyErr
		return
	}
	region, ok := region(accessToken)
	if !ok {
		err = ecode.AccessKeyErr
		return
	}
	if region == s.currentRegion {
		return s.currentOauth(c, accessToken)
	}
	if region == _origin {
		token, err = s.currentOauth(c, accessToken)
	} else {
		token, err = s.otherOauth(c, accessToken)
	}
	if err == nil && token != nil && token.Mid > 0 {
		return
	}
	if from != "" {
		s.dispatcherErrStats.Incr("dispatcher_error")
		err = ErrDispatcherError
		return
	}
	return s.originOauth(c, s.oauth[region], accessToken, s.currentRegion)
}

func (s *Service) originOauth(c context.Context, uri, accessKey, from string) (*model.Token, error) {
	cache := true
	token, err := s.d.OriginTokenCache(c, accessKey)
	if token != nil {
		if token.Mid > 0 {
			return token, nil
		}
		return nil, ecode.NoLogin
	}
	if err != nil {
		cache = false
		log.Error("Faield to get origin token cache with access key: %s: %+v", accessKey, err)
	}
	token, err = s.d.Oauth(c, uri, accessKey, from)
	if token != nil {
		// s.addCache(func() { s.d.SetOriginTokenCache(context.Background(), token) })
		return token, nil
	}

	log.Error("Failed to oauth token by origin: %s, %s, %s: %+v", uri, accessKey, from, err)
	if !cache {
		return nil, err
	}
	if ec := ecode.Cause(err); ec.Equal(ecode.AccessKeyErr) || ec.Equal(ecode.NoLogin) {
		token = &model.Token{
			Mid:         -1,
			AccessToken: accessKey,
		}
		s.addCache(func() { s.d.SetOriginTokenCache(context.Background(), token) })
	}
	return nil, err
}

func (s *Service) currentOauth(c context.Context, accessToken string) (token *model.Token, err error) {
	var t *model.Perm
	if t, err = s.tokenInfo(c, accessToken); err != nil {
		err = nil
	}
	return s.parseToken(c, t)
}

func (s *Service) otherOauth(c context.Context, accessToken string) (token *model.Token, err error) {
	var t *model.Perm
	if t, err = s.d.TokenFromOtherRegion(c, accessToken); err != nil {
		err = nil
	}
	return s.parseToken(c, t)
}

func (s *Service) parseToken(c context.Context, t *model.Perm) (token *model.Token, err error) {
	if t == nil {
		err = ecode.AccessKeyErr
		return
	}
	duration := time.Now().Unix() - t.Expires
	if duration > _gameAdditionalExpireSeconds {
		err = ecode.AccessKeyErr
		return
	}
	if duration > 0 && duration <= _gameAdditionalExpireSeconds {
		err = ecode.AccessTokenExpires
		return
	}
	accInfo := s.Info(c, t.Mid)
	token = &model.Token{
		Mid:         t.Mid,
		AppID:       t.AppID,
		AccessToken: t.AccessToken,
		CreateAt:    t.CreateAt,
		UserID:      accInfo.UserID,
		Uname:       accInfo.Uname,
		Expires:     t.Expires,
		Permission:  "ALL",
	}
	return
}

// thinOauth oauth and return mid.
func (s *Service) thinOauth(c context.Context, app *model.App, accessToken string) (res *model.Info, err error) {
	if accessToken == "" {
		err = ecode.AccessKeyErr
		return
	}
	region, ok := region(accessToken)
	if !ok {
		err = ecode.AccessKeyErr
		return
	}
	var mid int64
	if region == s.currentRegion {
		mid, err = s.thinCurrentOauth(c, app.AppKey, accessToken)
		if err != nil {
			return
		}
		res = &model.Info{
			Mid: mid,
		}
		return
	}
	if region == _origin {
		mid, err = s.thinCurrentOauth(c, app.AppKey, accessToken)
		if err == nil {
			res = &model.Info{
				Mid: mid,
			}
			return
		}
	}
	t, err := s.d.Oauth(c, s.oauth[region], accessToken, s.currentRegion)
	if err != nil {
		return
	}
	res = &model.Info{
		Mid:    t.Mid,
		UserID: t.UserID,
		Uname:  t.Uname,
	}
	return
}

func (s *Service) thinCurrentOauth(c context.Context, appKey, accessToken string) (res int64, err error) {
	var t *model.Perm
	if t, err = s.tokenInfo(c, accessToken); err != nil {
		err = nil
	}
	if t == nil {
		err = ecode.AccessKeyErr
		return
	}
	duration := time.Now().Unix() - t.Expires
	if duration > 0 && duration <= _gameAdditionalExpireSeconds {
		err = ecode.AccessTokenExpires
		return
	}
	res = t.Mid
	return
}
