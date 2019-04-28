package service

import (
	"context"
	"strconv"

	"go-common/app/admin/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// SettingList get all setting
func (s *Service) SettingList(c context.Context) (list []*model.Setting, err error) {
	if list, err = s.spyDao.SettingList(c); err != nil {
		log.Error("s.spyDao.SettingList() error(%v)", err)
		return
	}
	return
}

// UpdateSetting update setting
func (s *Service) UpdateSetting(c context.Context, name string, property string, val string) (err error) {
	if err = s.checkSettingVal(property, val); err != nil {
		return
	}
	var effected int64
	if effected, err = s.spyDao.UpdateSetting(c, property, val); err != nil {
		log.Error("s.spyDao.UpdateSetting(%s,%d) error(%v)", property, val, err)
		return
	}
	if effected > 0 {
		updatedSetting := &model.Setting{Property: property, Val: val}
		if err := s.AddLog(c, name, model.UpdateSetting, updatedSetting); err != nil {
			log.Error("s.AddLog(%s,%d,%+v) error(%v)", name, model.UpdateSetting, updatedSetting, err)
		}
	}
	return
}

func (s *Service) checkSettingVal(prop string, val string) (err error) {
	switch prop {
	case model.AutoBlock:
		var block int64
		if block, err = strconv.ParseInt(val, 10, 64); err != nil {
			err = ecode.SpySettingValTypeError
			return
		}
		if block != 1 && block != 0 {
			err = ecode.SpySettingValueOutOfRange
			return
		}
	case model.LimitBlockCount:
		var count int64
		if count, err = strconv.ParseInt(val, 10, 64); err != nil {
			err = ecode.SpySettingValTypeError
			return
		}
		if count < 0 {
			err = ecode.SpySettingValueOutOfRange
			return
		}
	case model.LessBlockScore:
		var score int64
		if score, err = strconv.ParseInt(val, 10, 64); err != nil {
			err = ecode.SpySettingValTypeError
			return
		}
		if score < 0 || score > 30 {
			err = ecode.SpySettingValueOutOfRange
			return
		}
	default:
		err = ecode.SpySettingUnknown
	}
	return err
}
