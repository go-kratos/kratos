package service

import (
	"context"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
)

var (
	_emtpyLimitUser = make([]*model.LimitUser, 0)
)

// LimitUsers limit users.
func (s *Service) LimitUsers(c context.Context, pn, ps int32) (res []*model.LimitUser, total int64, err error) {
	start := (pn - 1) * ps
	end := ps
	if total, err = s.dao.LimitUserCount(c); err != nil {
		return
	}
	if total == 0 {
		return _emtpyLimitUser, 0, nil
	}
	if res, err = s.dao.LimitUsers(c, start, end); err != nil {
		return
	}
	if len(res) == 0 {
		res = _emtpyLimitUser
	}
	return
}

// LimitUserAdd add limit user.
func (s *Service) LimitUserAdd(c context.Context, mid int64, cname string) (err error) {
	var (
		userInfo  *model.UserInfo
		userLimit *model.LimitUser
	)
	if userLimit, err = s.dao.LimitUser(c, mid); err != nil {
		return
	}
	if userLimit != nil {
		return ecode.TagLimitUserExist
	}
	if userInfo, err = s.userInfo(c, mid); err != nil || userInfo == nil {
		return
	}
	id, err := s.dao.InsertLimitUser(c, mid, userInfo.Name, cname)
	if err != nil || id <= 0 {
		err = ecode.TagUpdateLimitUserFail
	}
	return
}

func (s *Service) userInfo(c context.Context, mid int64) (user *model.UserInfo, err error) {
	userInfoReply, err := s.dao.UserInfo(c, mid)
	if err != nil {
		return
	}
	if userInfoReply == nil || userInfoReply.Info == nil || userInfoReply.Info.Mid != mid {
		err = ecode.UserNoMember
		return
	}
	user = &model.UserInfo{
		Mid:  userInfoReply.Info.Mid,
		Name: userInfoReply.Info.Name,
	}
	return
}

func (s *Service) userInfos(c context.Context, mids []int64) (unameMap map[int64]*model.UserInfo, err error) {
	unameMap = make(map[int64]*model.UserInfo, len(mids))
	userInfosReply, err := s.dao.UserInfos(c, mids)
	if err != nil || userInfosReply == nil {
		return
	}
	for _, u := range userInfosReply.Infos {
		unameMap[u.Mid] = &model.UserInfo{
			Mid:  u.Mid,
			Name: u.Name,
		}
	}
	return
}

// LimitUserDel LimitUserDel.
func (s *Service) LimitUserDel(c context.Context, mid int64) (err error) {
	if _, err := s.dao.DelLimitUser(c, mid); err != nil {
		err = ecode.TagUpdateLimitUserFail
	}
	return
}
