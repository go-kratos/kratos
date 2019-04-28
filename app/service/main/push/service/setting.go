package service

import (
	"context"

	"go-common/app/service/main/push/dao"
	"go-common/app/service/main/push/model"
	"go-common/library/log"
)

// Setting gets user notify setting.
func (s *Service) Setting(c context.Context, mid int64) (st map[int]int, err error) {
	st, err = s.dao.Setting(c, mid)
	if err != nil {
		log.Error("s.dao.Setting(%d) error(%v)", mid, err)
		return
	}
	if st == nil {
		st = make(map[int]int, len(model.Settings))
		for k, v := range model.Settings {
			st[k] = v
		}
	}
	return
}

// SetSetting saves setting.
func (s *Service) SetSetting(c context.Context, mid int64, typ, val int) (err error) {
	st, err := s.dao.Setting(c, mid)
	if err != nil {
		log.Error("s.dao.Setting(%d) error(%v)", mid, err)
		return
	}
	if st == nil {
		st = make(map[int]int, len(model.Settings))
		for k, v := range model.Settings {
			st[k] = v
		}
	}
	st[typ] = val
	if err = s.dao.SetSetting(c, mid, st); err != nil {
		log.Error("s.dao.AddSetting(%d,%v) error(%v)", mid, st, err)
		dao.PromError("setting:保存设置")
	}
	return
}
