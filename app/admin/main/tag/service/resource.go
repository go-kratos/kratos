package service

import (
	"context"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
)

// ResourceByOid ResourceByOid.
func (s *Service) ResourceByOid(c context.Context, oid int64, tp int32) (res *model.LimitRes, err error) {
	var arc *model.SearchRes
	if arc, err = s.arcInfo(c, oid); err != nil {
		return
	}
	if arc == nil {
		return nil, ecode.ArchiveNotExist
	}
	if res, err = s.dao.ResLimitByOid(c, oid, tp); err != nil {
		return
	}
	if res == nil {
		res = &model.LimitRes{
			Oid:       arc.ID,
			Type:      int64(tp),
			Operation: model.ResLimitNone,
		}
	}
	res.Title = arc.Title
	if arc.Mid <= 0 {
		return
	}
	if userInfo, _ := s.userInfo(c, arc.Mid); userInfo != nil {
		res.Author = userInfo.Name
	}
	return
}

// ResByOperate ResByOperate.
func (s *Service) ResByOperate(c context.Context, operateState, pn, ps int32) (res []*model.LimitRes, count int64, err error) {
	var (
		oids, authorIDs []int64
		arcs            map[int64]*model.SearchRes
		userInfoMap     map[int64]*model.UserInfo
	)
	start := (pn - 1) * ps
	end := ps
	if count, err = s.dao.ResLimitCount(c, operateState); err != nil {
		return
	}
	if res, oids, err = s.dao.ResLimitByOpState(c, operateState, start, end); err != nil {
		return
	}
	if len(oids) > 0 {
		arcs, authorIDs, _ = s.arcInfos(c, oids)
	}
	if len(authorIDs) > 0 {
		userInfoMap, _ = s.userInfos(c, authorIDs)
	}
	for _, v := range arcs {
		for _, k := range res {
			if v.ID != k.Oid {
				continue
			}
			k.Title = v.Title
			if u, ok := userInfoMap[v.Mid]; ok {
				k.Author = u.Name
			}
		}
	}
	return
}

// UpdateResLimitState UpdateResLimitState.
func (s *Service) UpdateResLimitState(c context.Context, oid int64, tp, operate int32) (err error) {
	var resource = new(model.LimitRes)
	if resource, err = s.dao.ResLimitByOid(c, oid, tp); err != nil {
		return ecode.TagOperateFail
	}
	if resource != nil {
		if _, err = s.dao.UpResLimitState(c, oid, tp, operate); err != nil {
			return ecode.TagOperateFail
		}
		return
	}
	if _, err = s.dao.ResLimitAdd(c, oid, tp, operate); err != nil {
		return ecode.TagOperateFail
	}
	return
}
