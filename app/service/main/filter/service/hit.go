package service

import (
	"context"
	"sync"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/model/actriearea"
	"go-common/library/log"
)

// HitV3 returns hit lists that contains message and risk level.
func (s *Service) HitV3(c context.Context, area, msg string, tpID int64, level int8) (res []*model.HitRes, err error) {
	if msg == "" {
		return
	}
	a := s.areas.Area(area)
	if a == nil {
		log.Error("s.HitV3 invalid area [%s]", area)
		return
	}
	return s.hitV3(c, area, msg, []byte(msg), tpID, a.IsFilterCommon(), level)
}

func (s *Service) hitV3(c context.Context, area string, msg string, oriMsg []byte, tpID int64, useComm bool, level int8) (hits []*model.HitRes, err error) {
	hits = make([]*model.HitRes, 0, 5)
	var (
		filters []*model.Filter
		allRegs []*model.Regexp
	)
	// 第一步： 白名单命中规则
	whiteHits := s.whites.Hits([]byte(msg), area, tpID)
	// 第二步： 黑名单命中规则
	if filters = s.filters.GetFilters(area, useComm); len(filters) == 0 {
		return
	}
	for _, f := range filters {
		allRegs = append(allRegs, f.Regs...)
	}
	var (
		blackHits []*actriearea.MatchHits
		mutex     sync.Mutex
		paraSize  = conf.Conf.Property.ParallelSize
		regLen    = len(allRegs)
		wg        = new(sync.WaitGroup)
	)
	// 正则命中
	for i := 0; i < regLen; i += paraSize {
		var regs []*model.Regexp
		if i+paraSize > regLen {
			regs = allRegs[i:]
		} else {
			regs = allRegs[i : i+paraSize]
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if x := recover(); x != nil {
					log.Error("service.artFilter(%s,%d,%t) error(%v)", area, tpID, useComm, x)
				}
			}()
			for _, re := range regs {
				for _, id := range re.TypeIDs {
					if id == 0 || id == tpID {
						if re.Reg.MatchString(msg) {
							if re.Level <= level {
								hit := &actriearea.MatchHits{Level: re.Level}
								hit.TypeIDs = append(hit.TypeIDs, id)
								hit.Rule = re.Reg.String()
								mutex.Lock()
								blackHits = append(blackHits, hit)
								mutex.Unlock()
							} else {
								positions := re.Reg.FindAllStringIndex(msg, -1)
								if len(positions) > 0 {
									hit := &actriearea.MatchHits{Positions: positions, Level: re.Level}
									hit.TypeIDs = append(hit.TypeIDs, id)
									hit.Rule = re.Reg.String()
									mutex.Lock()
									blackHits = append(blackHits, hit)
									mutex.Unlock()
								}
							}
						}
						break
					}
				}
			}
		}()
	}
	wg.Wait()
	// actrie hits
	for _, f := range filters {
		blackHits = append(blackHits, f.Matcher.Filter(msg, tpID, level)...)
	}
	// 终极过滤
	for _, blackhit := range blackHits {
		log.Info("article black hit(%+v)", blackhit)
		for _, posi := range blackhit.Positions {
			var whiteCover = false
			if len(whiteHits) > 0 {
			OUTER:
				for _, whhit := range whiteHits {
					for _, whPosi := range whhit.Positions {
						if whPosi[0] <= posi[0] && whPosi[1] >= posi[1] {
							whiteCover = true
							break OUTER
						}
					}
				}
			}
			if !whiteCover {
				hit := &model.HitRes{
					Msg:   string([]byte(msg)[posi[0]:posi[1]]),
					Level: blackhit.Level,
				}
				hits = append(hits, hit)
			}
		}
	}
	s.addEvent(func() {
		s.repostHitLog(context.Background(), area, msg, blackHits, "area")
	})
	return
}
