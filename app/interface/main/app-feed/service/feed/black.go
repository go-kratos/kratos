package feed

import (
	"context"
	"time"

	"go-common/library/log"
)

func (s *Service) loadBlackCache() {
	bs, err := s.blk.Black(context.Background())
	if err != nil {
		log.Error("s.blk.Black error(%v)", err)
		return
	}
	s.blackCache = bs
	log.Info("reBlackList success")
}

// blackproc load blacklist cache.
func (s *Service) blackproc() {
	for {
		time.Sleep(s.tick)
		s.loadBlackCache()
	}
}

func (s *Service) BlackList(c context.Context, mid int64) (aidm map[int64]struct{}, err error) {
	if mid == 0 {
		return
	}
	return s.blk.BlackList(c, mid)
}
