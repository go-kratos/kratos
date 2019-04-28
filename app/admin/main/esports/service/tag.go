package service

import (
	"context"

	"go-common/app/admin/main/esports/model"
	"go-common/library/log"
)

var _emptyTagList = make([]*model.Tag, 0)

// TagInfo .
func (s *Service) TagInfo(c context.Context, id int64) (data *model.Tag, err error) {
	data = new(model.Tag)
	if err = s.dao.DB.Model(&model.Tag{}).Where("id=?", id).First(&data).Error; err != nil {
		log.Error("TagInfo Error (%v)", err)
	}
	return
}

// TagList .
func (s *Service) TagList(c context.Context, pn, ps int64) (list []*model.Tag, count int64, err error) {
	s.dao.DB.Model(&model.Tag{}).Count(&count)
	if err = s.dao.DB.Model(&model.Tag{}).Offset((pn - 1) * ps).Limit(ps).Find(&list).Error; err != nil {
		log.Error("TagList Error (%v)", err)
	}
	return
}

// AddTag .
func (s *Service) AddTag(c context.Context, param *model.Tag) (err error) {
	if err = s.dao.DB.Model(&model.Tag{}).Create(param).Error; err != nil {
		log.Error("AddTag s.dao.DB.Model Create(%+v) error(%v)", param, err)
	}
	return
}

// EditTag .
func (s *Service) EditTag(c context.Context, param *model.Tag) (err error) {
	preData := new(model.Tag)
	if err = s.dao.DB.Where("id=?", param.ID).First(&preData).Error; err != nil {
		log.Error("EditTag s.dao.DB.Where id(%d) error(%d)", param.ID, err)
		return
	}
	if err = s.dao.DB.Model(&model.Tag{}).Update(param).Error; err != nil {
		log.Error("EditTag s.dao.DB.Model Update(%+v) error(%v)", param, err)
	}
	return
}

// ForbidTag .
func (s *Service) ForbidTag(c context.Context, id int64, state int) (err error) {
	preTag := new(model.Tag)
	if err = s.dao.DB.Where("id=?", id).First(&preTag).Error; err != nil {
		log.Error("TagForbid s.dao.DB.Where id(%d) error(%d)", id, err)
		return
	}
	if err = s.dao.DB.Model(&model.Tag{}).Where("id=?", id).Update(map[string]int{"status": state}).Error; err != nil {
		log.Error("TagForbid s.dao.DB.Model error(%v)", err)
	}
	return
}
