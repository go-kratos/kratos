package service

import (
	"context"
	"sync"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/model/actriearea"
	"go-common/library/log"
)

// Article .
func (s *Service) Article(c context.Context, area, msg string, tpID int64) (res []string, err error) {
	if msg == "" {
		return
	}
	return s.hit(c, area, msg, tpID, false)
}

// MultiHit 批量Hit方法
func (s *Service) MultiHit(c context.Context, area string, msgMap map[string]string, tpid int64) (rMap map[string][]string, err error) {
	rMap = make(map[string][]string, len(msgMap))
	if len(msgMap) == 0 {
		return
	}
	a := s.areas.Area(area)
	if a == nil {
		log.Error("s.MultiHit got invalid area [%s]", area)
		return
	}

	for id, msg := range msgMap {
		var (
			hits []string
		)
		if hits, err = s.hit(c, area, msg, tpid, a.IsFilterCommon()); err != nil {
			return
		}
		rMap[id] = hits
	}
	return
}

// Hit .
func (s *Service) Hit(c context.Context, area, msg string, tpID int64) (res []string, err error) {
	if msg == "" {
		return
	}
	a := s.areas.Area(area)
	if a == nil {
		log.Error("s.Hit invalid area [%s]", area)
		return
	}
	return s.hit(c, area, msg, tpID, a.IsFilterCommon())
}

func (s *Service) hit(c context.Context, area, msg string, tpID int64, useComm bool) (res []string, err error) {
	var (
		finds    []string
		allFinds []string
		resMap   map[string]struct{}
	)
	// area filter
	if finds, err = s.hitFilter(c, area, msg, []byte(msg), tpID, useComm); err != nil {
		return
	}
	allFinds = append(allFinds, finds...)

	resMap = make(map[string]struct{}, len(allFinds))
	for _, find := range allFinds {
		if _, ok := resMap[find]; !ok {
			res = append(res, find)
			resMap[find] = struct{}{}
		}
	}
	return
}

func (s *Service) hitFilter(c context.Context, area string, msg string, oriMsg []byte, tpID int64, useComm bool) (finds []string, err error) {
	finds = make([]string, 0, 5)
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
							if re.Level <= conf.Conf.Property.CriticalFilterLevel {
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
		blackHits = append(blackHits, f.Matcher.Filter(msg, tpID, conf.Conf.Property.CriticalFilterLevel)...)
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
				finds = append(finds, string([]byte(msg)[posi[0]:posi[1]]))
			}
		}
	}
	s.addEvent(func() {
		s.repostHitLog(context.Background(), area, msg, blackHits, "area")
	})
	return
}
