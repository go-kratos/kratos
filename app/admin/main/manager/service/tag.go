package service

import (
	"context"

	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddType .
func (s *Service) AddType(c context.Context, tt *model.TagType) (err error) {
	if err = s.dao.AddType(c, tt); err != nil {
		log.Error("s.dao.AddType error(%v)", err)
	}
	return
}

// UpdateType .
func (s *Service) UpdateType(c context.Context, tt *model.TagType) (err error) {
	if err = s.dao.UpdateTypeName(c, tt); err != nil {
		log.Error("s.dao.UpdateTypeName error(%v)", err)
		return
	}
	if err = s.dao.DeleteNonRole(c, tt); err != nil {
		log.Error("s.dao.DeleteNonRole error(%v)", err)
		return
	}
	if err = s.dao.UpdateType(c, tt); err != nil {
		log.Error("s.dao.UpdateType error(%v)", err)
	}
	return
}

// DeleteType .
func (s *Service) DeleteType(c context.Context, td *model.TagTypeDel) (err error) {
	// judge before deleting
	var typeRes []*model.Tag
	if typeRes, err = s.dao.TagByType(c, td.ID); err != nil {
		log.Error("s.dao.TagByType error(%v)", err)
		return
	}
	if len(typeRes) != 0 {
		err = ecode.ManagerTagTypeDelErr
		return
	}
	if err = s.dao.DeleteType(c, td); err != nil {
		log.Error("s.dao.DeleteType error(%v)", err)
	}
	return
}

// AddTag .
func (s *Service) AddTag(c context.Context, t *model.Tag) (err error) {
	if err = s.dao.AddTag(c, t); err != nil {
		log.Error("s.dao.AddTag error(%v)", err)
	}
	return
}

// UpdateTag .
func (s *Service) UpdateTag(c context.Context, t *model.Tag) (err error) {
	if err = s.dao.UpdateTag(c, t); err != nil {
		log.Error("s.dao.UpdateTag error(%v)", err)
	}
	return
}

// AddControl .
func (s *Service) AddControl(c context.Context, tc *model.TagControl) (err error) {
	if err = s.dao.AddControl(c, tc); err != nil {
		log.Error("s.dao.AddControl error(%v)", err)
	}
	return
}

// UpdateControl .
func (s *Service) UpdateControl(c context.Context, tc *model.TagControl) (err error) {
	if err = s.dao.UpdateControl(c, tc); err != nil {
		log.Error("s.dao.UpdateControl error(%v)", err)
	}
	return
}

// BatchUpdateState .
func (s *Service) BatchUpdateState(c context.Context, b *model.BatchUpdateState) (err error) {
	if err = s.dao.BatchUpdateState(c, b); err != nil {
		log.Error("s.dao.BatchUpdateState error(%v)", err)
	}
	return
}

// TagList .
func (s *Service) TagList(c context.Context, t *model.SearchTagParams) (res []*model.Tag, total int64, err error) {
	var (
		tRes map[int64]*model.TagType
		rRes map[int64]*model.BusinessRole
	)
	if t.UName != "" {
		if uid, ok := s.userIds[t.UName]; ok {
			t.UID = uid
		}
	}
	if res, err = s.dao.TagList(c, t); err != nil {
		log.Error("s.dao.TagList error(%v)", err)
		return
	}
	tids := []int64{}
	rids := []int64{}
	for _, r := range res {
		tids = append(tids, r.Tid)
		rids = append(rids, r.Rid)
	}
	if tRes, err = s.dao.TypeByIDs(c, tids); err != nil {
		log.Error("s.dao.TypeByIDs error(%v)", err)
		return
	}
	if rRes, err = s.dao.RoleByRIDs(c, t.Bid, rids); err != nil {
		log.Error("s.dao.RoleByIDs error(%v)", err)
		return
	}
	for _, r := range res {
		if tr, ok := tRes[r.Tid]; ok {
			r.TName = tr.Name
		}
		if rr, ok := rRes[r.Rid]; ok {
			r.RName = rr.Name
		}
		if u, ok := s.userNames[r.UID]; ok {
			r.UName = u
		}
	}
	total = int64(len(res))
	start := (t.PN - 1) * t.PS
	if start >= total {
		res = []*model.Tag{}
		return
	}
	end := start + t.PS
	if end > total {
		end = total
	}
	res = res[start:end]
	return
}

// TypeList .
func (s *Service) TypeList(c context.Context, tt *model.TagTypeList) (res []*model.TagType, err error) {
	var (
		tids  []int64
		rids  []int64
		tRole []*model.TagTypeRole
		rRes  map[int64]*model.BusinessRole
	)
	if res, err = s.dao.TagTypeByBID(c, tt.BID); err != nil {
		log.Error("s.dao.TagTypeByBID error(%v)", err)
		return
	}
	for _, r := range res {
		tids = append(tids, r.ID)
	}
	if tRole, err = s.dao.TagTypeRoleByTids(c, tids); err != nil {
		log.Error("s.dao.TagTypeRoleByTids error(%v)", err)
		return
	}
	for _, rt := range tRole {
		rids = append(rids, rt.Rid)
	}
	if rRes, err = s.dao.RoleByRIDs(c, tt.BID, rids); err != nil {
		log.Error("s.dao.RoleByRIDs error(%v)", err)
		return
	}
	for _, r := range res {
		for _, t := range tRole {
			if r.ID == t.Tid {
				if role, ok := rRes[t.Rid]; ok {
					r.Roles = append(r.Roles, role)
				}
			}
		}
	}
	return
}

// AttrList .
func (s *Service) AttrList(c context.Context, tba *model.TagBusinessAttr) (res *model.TagBusinessAttr, err error) {
	if res, err = s.dao.AttrList(c, tba.Bid); err != nil {
		log.Error("s.dao.AttrList error(%v)", err)
		return
	}
	if res.ID == 0 {
		t := &model.TagBusinessAttr{
			Bid:    tba.Bid,
			Button: model.DefaultButton,
		}
		if err = s.dao.InsertAttr(c, t); err != nil {
			log.Error("s.dao.InsertAttr error(%v)", err)
			return
		}
		if res, err = s.dao.AttrList(c, tba.Bid); err != nil {
			log.Error("s.dao.AttrList error(%v)", err)
			return
		}
	}
	return
}

// AttrUpdate .
func (s *Service) AttrUpdate(c context.Context, tba *model.TagBusinessAttr) (err error) {
	if err = s.dao.AttrUpdate(c, tba); err != nil {
		log.Error("s.dao.AttrUpdate error(%v)", err)
	}
	return
}

// TagControl .
func (s *Service) TagControl(c context.Context, tc *model.TagControlParam) (res []*model.TagControl, err error) {
	if res, err = s.dao.TagControl(c, tc); err != nil {
		log.Error("s.dao.TagControl error(%v)", err)
	}
	return
}
