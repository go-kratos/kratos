package show

import (
	"context"

	"go-common/library/log"
)

func (s *Service) loadAuditCache() {
	as, err := s.adt.Audits(context.TODO())
	if err != nil {
		log.Error("s.adt.Audits error(%v)", err)
		return
	}
	s.auditCache = as
}

// Audit show tab data list.
func (s *Service) auditTab(mobiApp string, build int, plat int8) (isAudit bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			return true
		}
	}
	return false
}
