package service

import (
	"context"
	"math"
	"time"

	"go-common/app/admin/main/reply/model"
	"go-common/library/log"
)

// AddReplyConfig create a new administrator configuration for reply business
func (s *Service) AddReplyConfig(c context.Context, m *model.Config) (id int64, err error) {
	sub, err := s.subject(c, m.Oid, m.Type)
	if err != nil {
		return
	}
	now := time.Now()
	if _, err = s.dao.AddConfig(c, m.Type, m.Category, m.Oid, m.AdminID, m.Operator, m.Config, now); err != nil {
		return
	}
	if m.ShowEntry == 1 && m.ShowAdmin == 1 {
		sub.AttrSet(model.AttrNo, model.SubAttrConfig)
	} else {
		sub.AttrSet(model.AttrYes, model.SubAttrConfig)
	}
	if _, err = s.dao.UpSubjectAttr(c, m.Oid, m.Type, sub.Attr, now); err != nil {
		log.Error("s.dao.UpSubjectAttr(%d,%d,%d,%d) error(%v)", m.Type, m.Oid, model.SubAttrConfig, m.ShowEntry, err)
		return
	}
	if err = s.dao.DelSubjectCache(c, m.Oid, m.Type); err != nil {
		log.Error("ReplyConfig del subject cache error(%v)", err)
	}
	if err = s.dao.DelConfigCache(c, m.Oid, m.Type, m.Category); err != nil {
		log.Error("ReplyConfig del config cache error(%v)", err)
	}
	return
}

// LoadReplyConfig load a configuration record of reply business.
func (s *Service) LoadReplyConfig(c context.Context, typ, category int32, oid int64) (m *model.Config, err error) {
	m, err = s.dao.LoadConfig(c, typ, category, oid)
	return
}

//PaginateReplyConfig paginate configuration list of records indexing from start to end, and a total count of records
func (s *Service) PaginateReplyConfig(c context.Context, typ, category int32, oid int64, operator string, offset, count int) (configs []*model.Config, totalCount, pages int64, err error) {
	configs, _ = s.dao.PaginateConfig(c, typ, category, oid, operator, offset, count)
	totalCount, _ = s.dao.PaginateConfigCount(c, typ, category, oid, operator)
	pages = int64(math.Ceil(float64(totalCount) / float64(count)))
	return
}

//RenewReplyConfig reset reply configuration by default, with deleting the detail configurations from db
func (s *Service) RenewReplyConfig(c context.Context, id int64) (result bool, err error) {
	now := time.Now()
	config, err := s.dao.LoadConfigByID(c, id)
	if err != nil {
		log.Error("s.dao.LoadConfigByID(%d) error(%v)", id, err)
	}
	if config == nil {
		return false, nil
	}
	sub, err := s.dao.Subject(c, config.Oid, config.Type)
	if err != nil {
		return
	}
	sub.AttrSet(model.AttrNo, model.SubAttrConfig)
	_, err = s.dao.UpSubjectAttr(c, config.Oid, config.Type, sub.Attr, now)
	if err != nil {
		log.Error("s.dao.UpSubjectAttr(%d,%d,%d,%d) error(%v)", config.Type, config.Oid, model.SubAttrConfig, config.ShowEntry, err)
		return
	}
	if _, err = s.dao.DeleteConfig(c, id); err != nil {
		log.Error("s.dao.DeleteConfig(%d) error(%v)", id, err)
		return
	}
	if err = s.dao.DelSubjectCache(c, config.Oid, config.Type); err != nil {
		log.Error("ReplyConfig del subject cache error(%v)", err)
	}
	if err = s.dao.DelConfigCache(c, config.Oid, config.Type, config.Category); err != nil {
		log.Error("ReplyConfig del config cache error(%v)", err)
	}
	result = true
	return
}
