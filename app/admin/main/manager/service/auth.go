package service

import (
	"go-common/app/admin/main/manager/model"
	"go-common/library/log"
)

func (s *Service) authItems(typ int) (items []*model.AuthItem, err error) {
	if err = s.dao.DB().Where("type = ?", typ).Find(&items).Error; err != nil {
		log.Error("s.authItem(%d) error(%v)", typ, err)
	}
	return
}

func (s *Service) assigns(id int) (assigns []*model.AuthAssign, err error) {
	if err = s.dao.DB().Where("item_id = ?", id).Find(&assigns).Error; err != nil {
		log.Error("s.assigns(%d) error(%v)", id, err)
	}
	return
}

// admins
func (s *Service) adms() (as map[int64]bool, err error) {
	assigns, err := s.assigns(model.Admin)
	if err != nil {
		return
	}
	as = make(map[int64]bool)
	for _, v := range assigns {
		as[v.UserID] = true
	}
	return
}

// points
func (s *Service) ptrs() (ps []*model.AuthItem, mps map[int64]*model.AuthItem, err error) {
	ps, err = s.authItems(model.TypePointer)
	if err != nil {
		return
	}
	mps = make(map[int64]*model.AuthItem)
	for _, v := range ps {
		mps[v.ID] = v
	}
	return
}

func (s *Service) roleAuths() (ra map[int64][]int64, err error) {
	var aic []*model.AuthItemChild
	if err = s.dao.DB().Joins("left join auth_item on auth_item_child.parent=auth_item.id").Where("auth_item.type=?", model.TypeRole).Find(&aic).Error; err != nil {
		log.Error("s.roleAuths() error(%v)", err)
		return
	}
	ra = make(map[int64][]int64)
	for _, v := range aic {
		ra[v.Parent] = append(ra[v.Parent], v.Child)
	}
	return
}

func (s *Service) groupAuths() (ga map[int64][]int64, err error) {
	var aic []*model.AuthItemChild
	if err = s.dao.DB().Joins("left join auth_item on auth_item_child.parent=auth_item.id").Where("auth_item.type=?", model.TypeGroup).Find(&aic).Error; err != nil {
		log.Error("s.groupAuths() error(%v)", err)
		return
	}
	ga = make(map[int64][]int64)
	for _, v := range aic {
		ga[v.Parent] = append(ga[v.Parent], v.Child)
	}
	return
}

func (s *Service) orgAuths() (gsInfo map[int64]*model.AuthOrg, err error) {
	var aic []*model.AuthOrg
	if err = s.dao.DB().Where("auth_item.type IN (?,?)", model.TypeGroup, model.TypeRole).Find(&aic).Error; err != nil {
		log.Error("s.groups() error(%v)", err)
		return
	}
	gsInfo = make(map[int64]*model.AuthOrg)
	for _, v := range aic {
		gsInfo[v.ID] = v
	}
	return
}
