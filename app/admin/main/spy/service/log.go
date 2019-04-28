package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/main/spy/model"
	"go-common/library/log"
)

// AddLog add log.
func (s *Service) AddLog(c context.Context, name string, m int8, v interface{}) (err error) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Error("AddLog json.Marshal(%v) error(%v)", v, err)
		return
	}
	l := &model.Log{Name: name, Module: m, Context: string(b)}
	_, err = s.spyDao.AddLog(c, l)
	if err != nil {
		log.Error("userDao spyDao.AddLog(%v) error(%v)", l, err)
		return
	}
	return
}

// AddLog2 add log.
func (s *Service) AddLog2(c context.Context, l *model.Log) (err error) {
	l.Ctime = time.Now()
	_, err = s.spyDao.AddLog(c, l)
	return
}

// LogList add log.
func (s *Service) LogList(c context.Context, refID int64, module int8) (list []*model.Log, err error) {
	list, err = s.spyDao.LogList(c, refID, module)
	return
}

// UpdateStateLog get log.
func (s *Service) UpdateStateLog(refID int64, state string) (log string) {
	log = fmt.Sprintf("%d : 状态修改为[%s]", refID, state)
	return
}

// AddRemarkLog add reamrk log.
func (s *Service) AddRemarkLog(refID int64, remark string) (log string) {
	log = fmt.Sprintf("%d : [备注]%s", refID, remark)
	return
}

// DeleteStatLog delete log .
func (s *Service) DeleteStatLog(refID int64) (log string) {
	log = fmt.Sprintf("%d 删除记录", refID)
	return
}

// UpdateStatCountLog update stat count log.
func (s *Service) UpdateStatCountLog(refID int64, old int64, new int64) (log string) {
	log = fmt.Sprintf("%d [修正数值] %d -> %d", refID, old, new)
	return
}
