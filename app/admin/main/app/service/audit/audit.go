package audit

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/app/conf"
	auditdao "go-common/app/admin/main/app/dao/audit"
	"go-common/app/admin/main/app/model/audit"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_initAuditKey = "audit_key_%v_%v"
)

// Service audit service.
type Service struct {
	dao *auditdao.Dao
}

// New new a audit service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: auditdao.New(c),
	}
	return
}

// Audits select All
func (s *Service) Audits(c context.Context) (res []*audit.Audit, err error) {
	if res, err = s.dao.Audits(c); err != nil {
		log.Error("s.dao.Audits error(%v)", err)
		return
	}
	return
}

// AuditByID by id
func (s *Service) AuditByID(c context.Context, id int64) (res *audit.Audit, err error) {
	if res, err = s.dao.AuditByID(c, id); err != nil {
		log.Error("s.dao.AuditByID error(%v)", err)
		return
	}
	return
}

// AddAudit insert audit
func (s *Service) AddAudit(c context.Context, a *audit.Param, now time.Time) (err error) {
	var (
		row map[string]*audit.Audit
	)
	key := fmt.Sprintf(_initAuditKey, a.MobiApp, a.Build)
	if row, err = s.dao.AuditExist(c, a); err != nil {
		log.Error("s.dao.AuditExist(%v)", err)
		return
	}
	if _, ok := row[key]; ok {
		err = ecode.NotModified
		return
	}
	if err = s.dao.Insert(c, a, now); err != nil {
		log.Error("s.dao.Insert(%v)", err)
		return
	}
	return
}

// UpdateAudit update audit
func (s *Service) UpdateAudit(c context.Context, a *audit.Param, now time.Time) (err error) {
	var (
		row map[string]*audit.Audit
	)
	key := fmt.Sprintf(_initAuditKey, a.MobiApp, a.Build)
	if row, err = s.dao.AuditExist(c, a); err != nil {
		log.Error("s.dao.AuditExist(%v)", err)
		return
	}
	if _, ok := row[key]; ok {
		err = ecode.NotModified
		return
	}
	if err = s.dao.Update(c, a, now); err != nil {
		log.Error("s.dao.Insert(%v)", err)
		return
	}
	return
}

// DelAudit del audit by id
func (s *Service) DelAudit(c context.Context, id int64) (err error) {
	if err = s.dao.Del(c, id); err != nil {
		log.Error("s.dao.Del(%v)", err)
		return
	}
	return
}
