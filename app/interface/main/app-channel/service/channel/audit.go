package channel

import (
	"context"

	"go-common/app/interface/main/app-channel/model"
	"go-common/library/log"
)

var (
	_auditRids = map[int8]map[int]struct{}{
		model.PlatIPad: map[int]struct{}{
			1:     struct{}{},
			160:   struct{}{},
			119:   struct{}{},
			155:   struct{}{},
			165:   struct{}{},
			5:     struct{}{},
			181:   struct{}{},
			65552: struct{}{},
			65556: struct{}{},
		},
		model.PlatIPhone: map[int]struct{}{
			1:     struct{}{},
			160:   struct{}{},
			119:   struct{}{},
			155:   struct{}{},
			165:   struct{}{},
			5:     struct{}{},
			181:   struct{}{},
			65552: struct{}{},
			65556: struct{}{},
		},
	}
)

// auditRegion region data list.
func (s *Service) auditRegion(mobiApp string, plat int8, build, rid int) (isAudit bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			if rids, ok := _auditRids[plat]; ok {
				if _, ok = rids[rid]; ok {
					return true
				}
			}
		}
	}
	return false
}

func (s *Service) auditList(mobiApp string, plat int8, build int) (isAudit bool) {
	if plats, ok := s.auditCache[mobiApp]; ok {
		if _, ok = plats[build]; ok {
			return true
		}
	}
	return false
}

func (s *Service) loadAuditCache() {
	as, err := s.adt.Audits(context.TODO())
	if err != nil {
		log.Error("s.adt.Audits error(%v)", err)
		return
	}
	s.auditCache = as
}
