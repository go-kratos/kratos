package resource

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

// checkDiff check diff between ads and ads_active
func (s *Service) checkDiff() {
	now := time.Now()
	s.activeVideos(context.Background(), now)
}

// activeVideos check if VideoAds active need to update
func (s *Service) activeVideos(c context.Context, now time.Time) {
	amtime, err := s.resdao.ActVideoMTimeCount(c)
	if err != nil {
		log.Error("resdao.ActVideoMTimeCount error(%v)", err)
		return
	}
	dmtime, err := s.resdao.AdVideoMTimeCount(c, now)
	if err != nil {
		log.Error("resdao.AdVideoMTimeCount error(%v)", err)
		return
	}
	if amtime == dmtime {
		log.Info("all video active ad are same")
		return
	}
	log.Info("video active avg mtime(%d), ads avg mtime(%d)", amtime, dmtime)
	if err = s.resdao.DelAllVideo(c); err != nil {
		log.Error("sdDao.DelAllVideo(), err (%v)", err)
		return
	}
	ads, err := s.resdao.AllVideoActive(c, now)
	if err != nil {
		log.Error("resdao.AllVideoActive(%v), err (%v)", now, err)
		return
	}
	tx, err := s.resdao.BeginTran(c)
	if err != nil {
		log.Error("BeginTran(), err (%v)", err)
		return
	}
	for _, ad := range ads {
		aids := strings.Split(ad.AidS, ",")
		for _, aid := range aids {
			i, e := strconv.ParseInt(aid, 10, 64)
			if e != nil {
				log.Error("strconv.ParseInt() error(%v)", e)
				continue
			}
			ad.Aid = i
			ad.MTime = dmtime
			if err = s.resdao.TxInsertVideo(tx, ad); err != nil {
				if err = tx.Rollback(); err != nil {
					log.Error("tx.Rollback(), err (%v)", err)
				}
				log.Error("resdao.TxInsertVideo(tx, %v), err(%v)", ad, err)
				return
			}
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(), error(%v)", err)
	}
}
