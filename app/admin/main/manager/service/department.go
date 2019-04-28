package service

import (
	"go-common/app/admin/main/manager/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Departments .
func (s *Service) Departments(c *bm.Context) (res []*model.DepartCustom, err error) {
	if err = s.dao.DB().Table("user_department").Where("status = ?", model.ValidateState).Find(&res).Error; err != nil {
		log.Error("s.Departments err (%v)", err)
	}
	return
}

// Roles .
func (s *Service) Roles(c *bm.Context) (res []*model.RoleCustom, err error) {
	if err = s.dao.DB().Table("auth_item").Where("type = ?", model.TypeRole).Find(&res).Error; err != nil {
		log.Error("s.Roles error (%v)", err)
	}
	return
}

// UsersByDepartment .
func (s *Service) UsersByDepartment(c *bm.Context, ID int64) (res []*model.UserCustom, err error) {
	if err = s.dao.DB().Table("user").Where("state =? and department_id = ?", model.ValidateState, ID).Find(&res).Error; err != nil {
		log.Error("s.UsersByDepartment ID (%d) error (%v)", ID, err)
	}
	return
}

// UsersByRole .
func (s *Service) UsersByRole(c *bm.Context, ID int64) (res []*model.UserCustom, err error) {
	resAuthAssign := []*model.AuthAssign{}
	if err = s.dao.DB().Table("auth_assignment").Where("item_id = ?", ID).Find(&resAuthAssign).Error; err != nil {
		log.Error("s.UsersByRole ID (%d) error (%v)", ID, err)
		return
	}
	for _, v := range resAuthAssign {
		temp := &model.UserCustom{
			ID:       v.UserID,
			Username: s.userNames[v.UserID],
			Nickname: s.userNicknames[v.UserID],
		}
		res = append(res, temp)
	}
	return
}
