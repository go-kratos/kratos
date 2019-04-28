package service

import (
	"go-common/app/job/main/search/model"
)

// stat get stat
func (s *Service) stat(appid string) (st *model.Stat) {
	s.mutex.RLock()
	st = s.stats[appid]
	s.mutex.RUnlock()
	return
}

// updateStat update stat
func (s *Service) updateStat(appid string, st *model.Stat) {
	s.mutex.Lock()
	s.stats[appid] = st
	s.mutex.Unlock()
}
