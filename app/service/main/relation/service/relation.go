package service

import (
	"context"

	"go-common/app/service/main/relation/model"
)

// Relation get relation of mid -> fid.
// black(128) friend(4) following(2) whisper(1), null for no relation.
func (s *Service) Relation(c context.Context, mid int64, fid int64) (f *model.Following, err error) {
	res, err := s.Relations(c, mid, []int64{fid})
	return res[fid], err
}

// Relations get relations of mid -> fids.
// black(128) friend(4) following(2) whisper(1), key absent for no relation.
func (s *Service) Relations(c context.Context, mid int64, fids []int64) (rm map[int64]*model.Following, err error) {
	var (
		ok  bool
		fid int64
		f   *model.Following
		fs  []*model.Following
		arm map[int64]*model.Following
	)
	if mid <= 0 {
		return
	}
	for _, v := range fids {
		if v <= 0 {
			return
		}
	}
	if rm, err = s.dao.RelationsCache(c, mid, fids); err != nil {
		err = nil
	} else if len(rm) > 0 {
		delete(rm, 0)
		return
	}
	if fs, err = s.followings(c, mid); err != nil {
		return
	} else if len(fs) == 0 {
		rm = _emptyFollowingMap
		return
	}
	rm = make(map[int64]*model.Following, len(fids))
	arm = make(map[int64]*model.Following, len(fs))
	for _, f = range fs {
		arm[f.Mid] = f
	}
	for _, fid = range fids {
		if f, ok = arm[fid]; ok {
			rm[f.Mid] = f
		}
	}
	return
}
