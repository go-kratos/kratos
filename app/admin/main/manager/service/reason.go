package service

import (
	"context"
	"strconv"
	"strings"

	"go-common/app/admin/main/manager/model"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// CateSecExtList .
func (s *Service) CateSecExtList(c context.Context, e *model.CateSecExt) (res []*model.CateSecExt, err error) {
	if res, err = s.dao.CateSecExtList(c, e); err != nil {
		log.Error("s.CateSecExtList type (%d) error (%v)", e.Type, err)
	}
	return
}

// AssociationList .
func (s *Service) AssociationList(c context.Context, e *model.Association) (res []*model.Association, err error) {
	var (
		aRes       []*model.Association
		rRes       []*model.BusinessRole
		catesecRes []*model.CateSecExt
	)
	// Display all record of Association
	if aRes, err = s.dao.AssociationList(c, model.AllState, e.BusinessID); err != nil {
		log.Error("s.AssociationList error (%v)", err)
		return
	}
	if len(aRes) <= 0 {
		res = []*model.Association{}
		return
	}
	// Get all mapping of role
	rMap := make(map[int64]string)
	bean := &model.BusinessRole{
		BID:   e.BusinessID,
		Type:  model.AllType,
		State: model.AllState,
	}
	if rRes, err = s.dao.RoleListByBID(c, bean); err != nil {
		log.Error("s.RoleListByBID error (%v)", err)
		return
	}
	for _, r := range rRes {
		rMap[r.RID] = r.Name
	}
	// GET all mapping of category and second
	cMap := make(map[int64]string)
	sMap := make(map[int64]*model.CateSecExt)
	if catesecRes, err = s.dao.CateSecList(c, e.BusinessID); err != nil {
		log.Error("s.CateSecList error (%v)", err)
		return
	}
	for _, cs := range catesecRes {
		if cs.Type == model.CategoryCode {
			cMap[cs.ID] = cs.Name
		} else if cs.Type == model.SecondeCode {
			sMap[cs.ID] = cs
		}
	}
	res = []*model.Association{}
	for _, value := range aRes {
		temp := &model.Association{
			ID:           value.ID,
			RoleID:       value.RoleID,
			BusinessID:   value.BusinessID,
			RoleName:     rMap[value.RoleID],
			CategoryID:   value.CategoryID,
			CategoryName: cMap[value.CategoryID],
			SecondIDs:    value.SecondIDs,
			Ctime:        value.Ctime,
			Mtime:        value.Mtime,
			State:        value.State,
		}
		sids := strings.Split(value.SecondIDs, ",")
		temp.Child = []*model.CateSecExt{}
		if sids[0] != "" {
			for _, sid := range sids {
				id, _ := strconv.ParseInt(sid, 10, 64)
				if value, ok := sMap[id]; ok {
					temp.Child = append(temp.Child, value)
				}
			}
		}
		res = append(res, temp)
	}
	return
}

// AddCateSecExt .
func (s *Service) AddCateSecExt(c context.Context, arg *model.CateSecExt) (err error) {
	if err = s.dao.AddCateSecExt(c, arg); err != nil {
		log.Error("s.AddCateSecExt (%s) error (%v)", arg.Name, err)
	}
	return
}

// UpdateCateSecExt .
func (s *Service) UpdateCateSecExt(c context.Context, arg *model.CateSecExt) (err error) {
	if err = s.dao.UpdateCateSecExt(c, arg); err != nil {
		log.Error("s.UpdateCateSecExt (%s) error (%v)", arg.Name, err)
	}
	return
}

// BanCateSecExt .
func (s *Service) BanCateSecExt(c context.Context, arg *model.CateSecExt) (err error) {
	if err = s.dao.BanCateSecExt(c, arg); err != nil {
		log.Error("s.BanCateSecExt (%d) error (%v)", arg.ID, err)
	}
	return
}

// AddAssociation .
func (s *Service) AddAssociation(c context.Context, arg *model.Association) (err error) {
	if err = s.dao.AddAssociation(c, arg); err != nil {
		log.Error("s.AddAssociation error %v", err)
	}
	return
}

// UpdateAssociation .
func (s *Service) UpdateAssociation(c context.Context, arg *model.Association) (err error) {
	if err = s.dao.UpdateAssociation(c, arg); err != nil {
		log.Error("s.UpdateAssociation error %v", err)
	}
	return
}

// BanAssociation .
func (s *Service) BanAssociation(c context.Context, arg *model.Association) (err error) {
	if err = s.dao.BanAssociation(c, arg); err != nil {
		log.Error("s.BanAssociation error %v", err)
	}
	return
}

// AddReason .
func (s *Service) AddReason(c context.Context, arg *model.Reason) (err error) {
	if err = s.dao.AddReason(c, arg); err != nil {
		log.Error("s.AddReason (%v) error (%v)", arg, err)
	}
	return
}

// UpdateReason .
func (s *Service) UpdateReason(c context.Context, arg *model.Reason) (err error) {
	if err = s.dao.UpdateReason(c, arg); err != nil {
		log.Error("s.UpdateReason (%v) error (%v)", arg, err)
	}
	return
}

// ReasonList .
func (s *Service) ReasonList(c context.Context, e *model.SearchReasonParams) (res []*model.Reason, total int64, err error) {
	var (
		rids  []int64
		csids []int64
		eRes  map[int64]*model.BusinessRole
		csRes map[int64]string
	)
	// Search the user_id
	if e.UName != "" {
		if userID, ok := s.userIds[e.UName]; ok {
			e.UID = userID
		}
	}
	if res, err = s.dao.ReasonList(c, e); err != nil {
		log.Error("s.dao.ReasonList error (%v)", err)
		return
	}
	if len(res) <= 0 {
		return
	}
	// Search relation data
	for _, r := range res {
		rids = append(rids, r.RoleID)
		csids = append(csids, r.CategoryID)
		csids = append(csids, r.SecondID)
	}
	if eRes, err = s.dao.RoleByRIDs(c, e.BusinessID, rids); err != nil {
		log.Error("s.dao.ExecutorByIDs error (%v)", err)
		return
	}
	if csRes, err = s.dao.CateSecByIDs(c, csids); err != nil {
		log.Error("s.dao.CateSecByIDs error (%v)", err)
		return
	}
	for _, value := range res {
		if r, ok := eRes[value.RoleID]; ok {
			value.RoleName = r.Name
		}
		if c, ok := csRes[value.CategoryID]; ok {
			value.CategoryName = c
		}
		if s, ok := csRes[value.SecondID]; ok {
			value.SecondName = s
		}
		if u, ok := s.userNames[value.UID]; ok {
			value.UName = u
		}
	}
	total = int64(len(res))
	start := (e.PN - 1) * e.PS
	if start >= total {
		res = []*model.Reason{}
		return
	}
	end := start + e.PS
	if end > total {
		end = total
	}
	res = res[start:end]
	return
}

// BatchUpdateReasonState .
func (s *Service) BatchUpdateReasonState(c context.Context, b *model.BatchUpdateReasonState) (err error) {
	if err = s.dao.BatchUpdateReasonState(c, b); err != nil {
		log.Error("s.dao.BatchUpdateReasonState %v error (%v)", b.IDs, err)
	}
	return
}

// DropDownList .
func (s *Service) DropDownList(c context.Context, e *model.Association) (res []*model.DropList, err error) {
	var (
		aRes  []*model.Association
		rRes  []*model.BusinessRole
		csRes []*model.CateSecExt
	)
	// Only display validate association
	if aRes, err = s.dao.AssociationList(c, model.ValidateState, e.BusinessID); err != nil {
		log.Error("s.AssociationList error (%v)", err)
		return
	}
	// Get all association record
	rMap := make(map[int64]string)
	bean := &model.BusinessRole{
		BID:   e.BusinessID,
		Type:  model.AllType,
		State: model.AllState,
	}
	if rRes, err = s.dao.RoleListByBID(c, bean); err != nil {
		log.Error("s.RoleListByBID error (%v)", err)
		return
	}
	for _, r := range rRes {
		rMap[r.RID] = r.Name
	}
	// Get all mapping of category and second
	csMap := make(map[int64]string)
	if csRes, err = s.dao.CateSecList(c, e.BusinessID); err != nil {
		log.Error("s.CateSecList error (%v)", err)
		return
	}
	for _, cs := range csRes {
		csMap[cs.ID] = cs.Name
	}
	// Mapping data,use map store unique data
	resMap := make(map[int64]map[int64]map[int64]int64)
	for _, v := range aRes {
		if _, ok := resMap[v.RoleID]; !ok {
			//first element
			resMap[v.RoleID] = make(map[int64]map[int64]int64)
		}
		if _, ok := resMap[v.RoleID][v.CategoryID]; !ok {
			resMap[v.RoleID][v.CategoryID] = make(map[int64]int64)
		}
		sids := strings.Split(v.SecondIDs, ",")
		for _, sid := range sids {
			id, _ := strconv.ParseInt(sid, 10, 64)
			resMap[v.RoleID][v.CategoryID][id] = id
		}
	}
	// Output the data
	res = []*model.DropList{}
	for keyR, valueR := range resMap {
		temprRes := &model.DropList{
			ID:   keyR,
			Name: rMap[keyR],
		}
		childCategory := []*model.DropList{}
		for keyC, valueC := range valueR {
			tempcRes := &model.DropList{
				ID:   keyC,
				Name: csMap[keyC],
			}
			childSecond := []*model.DropList{}
			for keyS, valueS := range valueC {
				tempsRes := &model.DropList{
					ID:    keyS,
					Name:  csMap[valueS],
					Child: []*model.DropList{},
				}
				childSecond = append(childSecond, tempsRes)
			}
			tempcRes.Child = childSecond
			childCategory = append(childCategory, tempcRes)
		}
		temprRes.Child = childCategory
		res = append(res, temprRes)
	}
	return
}

// BusinessAttr .
func (s *Service) BusinessAttr(c context.Context, b *model.BusinessAttr) (res map[string]bool, err error) {
	res = make(map[string]bool)
	res["isTag"] = false
	var tempRes []*model.CateSecExt
	if tempRes, err = s.dao.CateSecList(c, b.BID); err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return
	}
	for _, tr := range tempRes {
		if tr.Type == model.ExtensionCode && tr.State == 1 {
			res["isTag"] = true
		}
	}
	return
}
