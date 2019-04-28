package like

import (
	"context"
	"sync/atomic"
	"time"

	"go-common/app/interface/main/activity/model/like"
	"go-common/library/log"
)

// RedDot get hot dot.
func (s *Service) RedDot(c context.Context, mid int64) (redDot *like.RedDot, err error) {
	var lastTs int64
	redDot = new(like.RedDot)
	if mid <= 0 {
		return
	}
	if lastTime, e := s.dao.CacheRedDotTs(c, mid); e != nil {
		log.Error("s.dao.CacheRedDotTs mid(%d) error(%+v)", mid, e)
	} else {
		lastTs = lastTime
	}
	if s.newestSubTs > lastTs {
		redDot.RedDot = true
	}
	return
}

// ClearRetDot clear red dot.
func (s *Service) ClearRetDot(c context.Context, mid int64) (err error) {
	if err = s.dao.AddCacheRedDotTs(c, mid, time.Now().Unix()); err != nil {
		log.Error("s.dao.AddCacheRedDotTs mid(%d) error(%+v)", mid, err)
	}
	return
}

func (s *Service) newestSubTsproc() {
	for {
		if like, err := s.dao.NewestSubject(context.Background(), like.VIDEOALL); err != nil || like == nil {
			log.Error("actNewTsproc s.dao.NewestSubject error(%+v)", err)
			time.Sleep(5 * time.Second)
		} else {
			newTs := like.Ctime.Time().Unix()
			if newTs > s.newestSubTs {
				atomic.StoreInt64(&s.newestSubTs, newTs)
			}
		}
		time.Sleep(time.Duration(s.c.Interval.NewestSubTsInterval))
	}
}
