package pgc

import (
	"context"
	"fmt"
	"strings"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
	"go-common/library/xstr"
)

// batchFilter picks a batch of seasonCMS data, to define their newest ep and update the struct
func (s *Service) batchFilter(ctx context.Context, snCMS []*model.SeasonCMS) {
	if len(snCMS) == 0 {
		return
	}
	for _, v := range snCMS {
		if newest, err := s.newestNB(v.SeasonID); err != nil || newest == 0 {
			continue
		} else {
			v.NewestNb = newest
		}
	}
}

// newestNB picks all the eps of the season and do the title fitler to calculate the newest episode
func (s *Service) newestNB(sid int) (newest int, err error) {
	var (
		keywords = s.c.Cfg.TitleFilter
		strategy = s.c.Cfg.LessStrategy
	)
	eps, err := s.dao.AllEP(ctx, sid, strategy)
	if err != nil {
		log.Warn("AllEP newestNB SeasonID %d, Err %v", sid, err)
		return
	}
	for _, v := range eps {
		if titleCheck(keywords, v.Title) {
			continue
		}
		newest++
	}
	if newest == 0 {
		log.Warn("AllEP newestNB SeasonID %d, After Filter it's empty", sid)
	}
	return
}

// titleCheck checks whether the title matches some forbidden keywords
func titleCheck(keywords []string, title string) bool {
	for _, v := range keywords {
		if strings.Contains(title, v) {
			return true
		}
	}
	return false
}

func (s *Service) cmsShelve() {
	var (
		ctx           = context.Background()
		cfg           = s.c.Cfg.Merak
		validMap      map[int64]int
		onIDs, offIDs []int64
		err           error
	)
	if validMap, err = s.cmsDao.ValidSns(ctx, cfg.Onlyfree); err != nil {
		log.Error("cmsShelve ValidSns Err %v", err)
		return
	}
	if onIDs, offIDs, err = s.cmsDao.ShelveOp(ctx, validMap); err != nil {
		log.Error("cmsShelve ShelveOp err %v", err)
		return
	}
	if len(onIDs) > 0 {
		if err = s.cmsDao.ActOps(ctx, onIDs, true); err != nil {
			log.Error("cmsShelve ActOps OnIDs %v, Err %v", onIDs, err)
		}
	}
	if len(offIDs) > 0 {
		if err = s.cmsDao.ActOps(ctx, offIDs, false); err != nil {
			log.Error("cmsShelve ActOps OffIDs %v, Err %v", offIDs, err)
		}
	}
	log.Info("cmsShelve OnIDs %v, OffIDs %v", onIDs, offIDs)
	content := fmt.Sprintf(cfg.Template, xstr.JoinInts(onIDs), xstr.JoinInts(offIDs))
	if err = s.cmsDao.MerakNotify(ctx, cfg.Title, content); err != nil {
		log.Error("Merak Content %s, Err %v", content, err)
	}
}
