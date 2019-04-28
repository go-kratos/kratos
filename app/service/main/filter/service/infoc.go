package service

import (
	"context"
	"strconv"
	"time"

	"go-common/library/log"

	"go-common/app/service/main/filter/model/actriearea"
)

func (s *Service) repostHitLog(c context.Context, resArea, msg string, mh []*actriearea.MatchHits, hitType string) {
	if s.infoc != nil {
		hitTime := strconv.FormatInt(time.Now().Unix(), 10)
		for _, h := range mh {
			s.infoc.Info(resArea, h.Area, hitTime, hitType, h.Fid, h.Rule, h.Level, msg)
			log.Info("%s, %s, %s, %s, %s, %s, %s, %s", resArea, h.Area, hitTime, hitType, h.Fid, h.Rule, h.Level, msg)
		}
	}
}
