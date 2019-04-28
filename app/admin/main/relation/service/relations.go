package service

import (
	"context"

	"go-common/app/admin/main/relation/model"
	relationPB "go-common/app/service/main/relation/api"
)

// Followers is.
func (s *Service) Followers(ctx context.Context, param *model.FollowersParam) (*model.FollowersListPage, error) {
	list, err := s.dao.Followers(ctx, param.Fid, param.Mid)
	if err != nil {
		return nil, err
	}
	list.OrderByMTime(param.Desc())
	from, err := model.ParseTime(param.MTimeFrom)
	if err == nil {
		list = list.FilterMTimeFrom(from)
	}
	to, err := model.ParseTime(param.MTimeTo)
	if err == nil {
		list = list.FilterMTimeTo(to)
	}
	plist := list.Paginate(param.PS*(param.PN-1), param.PS)
	flist := plist.FollowersList()

	uids := make([]int64, 0, len(flist)*2)
	for _, r := range flist {
		uids = append(uids, r.Mid)
		uids = append(uids, r.Fid)
	}
	uinfos, err := s.dao.RPCInfos(ctx, uids)
	if err != nil {
		return nil, err
	}

	for _, r := range flist {
		if mi, ok := uinfos[r.Mid]; ok {
			r.MemberName = mi.Name
		}
	}
	for _, r := range flist {
		if fi, ok := uinfos[r.Fid]; ok {
			r.FollowerName = fi.Name
		}
	}

	page := &model.FollowersListPage{}
	page.Sort = param.Sort
	page.Order = param.Order
	page.List = flist
	page.TotalCount = len(list)
	page.PN = param.PN
	page.PS = param.PS
	return page, nil
}

// Followings is.
func (s *Service) Followings(ctx context.Context, param *model.FollowingsParam) (*model.FollowingsListPage, error) {
	list, err := s.dao.Followings(ctx, param.Mid, param.Fid)
	if err != nil {
		return nil, err
	}
	list.OrderByMTime(param.Desc())
	from, err := model.ParseTime(param.MTimeFrom)
	if err == nil {
		list = list.FilterMTimeFrom(from)
	}
	to, err := model.ParseTime(param.MTimeTo)
	if err == nil {
		list = list.FilterMTimeTo(to)
	}
	plist := list.Paginate(param.PS*(param.PN-1), param.PS)
	flist := plist.FollowingsList()

	uids := make([]int64, 0, len(flist)*2)
	for _, r := range flist {
		uids = append(uids, r.Mid)
		uids = append(uids, r.Fid)
	}
	minfos, err := s.dao.RPCInfos(ctx, uids)
	if err != nil {
		return nil, err
	}
	for _, r := range flist {
		if mi, ok := minfos[r.Mid]; ok {
			r.MemberName = mi.Name
		}
	}
	for _, r := range flist {
		if fi, ok := minfos[r.Fid]; ok {
			r.FollowingName = fi.Name
		}
	}

	page := &model.FollowingsListPage{}
	page.Sort = param.Sort
	page.Order = param.Order
	page.List = flist
	page.TotalCount = len(list)
	page.PN = param.PN
	page.PS = param.PS
	return page, nil
}

// Stat is
func (s *Service) Stat(ctx context.Context, param *model.ArgMid) (*relationPB.StatReply, error) {
	return s.dao.Stat(ctx, param.Mid)
}

// Stats is
func (s *Service) Stats(ctx context.Context, param *model.ArgMids) (map[int64]*relationPB.StatReply, error) {
	return s.dao.Stats(ctx, param.Mids)
}
