package service

import (
	"context"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/sync/errgroup"
)

// MyInfo get user's info by token.
func (s *Service) MyInfo(c context.Context, app *model.App, accessKey string) (*model.Info, error) {
	ti, err := s.thinOauth(c, app, accessKey)
	if err != nil {
		return nil, err
	}
	return s.Info(c, ti.Mid), nil
}

// Info get user's info by mid.
func (s *Service) Info(c context.Context, mid int64) (res *model.Info) {
	var err error
	cache := true
	if res, err = s.d.InfoCache(c, mid); err != nil {
		err = nil
		cache = false
	} else if res != nil {
		return
	}
	var userid, uname, face string
	var email, tel *string
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() (err error) {
		accInfo, err := s.d.AccountInfo(errCtx, mid)
		if err != nil || accInfo == nil {
			userid = model.DefaultUserID(mid)
			uname = model.DefaultUname(mid)
			return
		}
		userid = accInfo.UserID
		uname = accInfo.Uname
		email = accInfo.Email
		tel = accInfo.Tel
		return
	})
	eg.Go(func() (err error) {
		memInfo, err := s.d.MemberInfo(errCtx, mid)
		if err != nil || memInfo == nil {
			face = model.EmptyFace
			return
		}
		face = memInfo.FullFace()
		return
	})
	eg.Wait()
	info := &model.Info{
		Mid:    mid,
		UserID: userid,
		Uname:  uname,
		Face:   face,
	}
	if email != nil {
		info.HasEmail = true
	}
	if tel != nil {
		info.HasTel = true
	}
	if cache {
		s.addCache(func() {
			s.d.SetInfoCache(context.Background(), info)
		})
	}
	res = info
	return
}
