package goblin

import (
	"time"

	"go-common/library/log"
)

// reload hotword data from MC
func (s *Service) loadHotword() {
	var err error
	if s.Hotword, err = s.dao.Hotword(ctx); err != nil {
		log.Error("loadHotword Error %v, List %v", err, s.Hotword)
		return
	}
}

// load hotword data regularly
func (s *Service) loadHotwordproc() {
	for {
		time.Sleep(time.Duration(s.conf.Search.HotwordFre))
		s.loadHotword()
	}
}
