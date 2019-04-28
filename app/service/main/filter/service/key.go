package service

import (
	"context"
	"fmt"
	"regexp"

	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/model/actriearea"
	"go-common/app/service/main/filter/model/actriekey"
	"go-common/app/service/main/filter/model/lrulist"
	"go-common/library/log"
)

// todo 替换*号规则和业务过滤一致
func (s *Service) opFils(msg string, inLevel int8, fils []*model.KeyFilter) (out string, level int8, hitRules []string) {
	out = msg
	level = inLevel
	var (
		msgBytes       = []byte(msg)
		sortedBlackHit = &sortBlackHits{}
		blackHits      []*actriearea.MatchHits
	)
	for _, fil := range fils {
		for _, re := range fil.Regs {
			positions := re.Reg.FindAllIndex(msgBytes, -1)
			if len(positions) > 0 {
				for _, posi := range positions {
					sortedBlackHit.Add(posi)
				}
				s.areaBlackHitProm.Incr("danmu")
				log.Infov(context.Background(), log.KV("log", "filter black hit by key"), log.KV("area", "danmu"), log.KV("hit", re.Reg.String()), log.KV("msg", msg), log.KV("KeyID", re.FkID), log.KV("option", positions))
				hitRules = append(hitRules, re.Reg.String())
				if re.Level > inLevel {
					level = re.Level
				} else {
					level = inLevel
				}
				blackHits = append(blackHits, &actriearea.MatchHits{Fid: re.FkID, Rule: re.Reg.String(), Level: re.Level, Area: re.Area})
			}
		}
		var (
			acLevel int8
			acRules []string
		)
		out, acLevel, acRules = fil.Matcher.Filter(out)
		hitRules = append(hitRules, acRules...)
		if acLevel > level {
			level = acLevel
		}
	}
	out, _ = sortedBlackHit.CoverByStart(msg)
	s.addEvent(func() {
		s.repostHitLog(context.Background(), "danmu", msg, blackHits, "key")
	})
	return
}

func (s *Service) fils(c context.Context, key string, areas []string) (fils []*model.KeyFilter, err error) {
	var (
		lruNode  *lrulist.Node
		missArea []string
		cs       []*model.KeyAreaInfo
	)
	for _, area := range areas {
		s.lruLock.Lock()
		lruNode = s.lruList.Find(s.lruKey(key, area))
		if lruNode != nil {
			fil := lruNode.Value.(*model.KeyFilter)
			fils = append(fils, fil)
			if lruNode == s.lruList.Header() {
				s.lruLock.Unlock()
				continue
			}
			// 热点数据移动到头部
			s.lruList.MoveToHead(s.lruList.Remove(lruNode))
		} else {
			missArea = append(missArea, area)
		}
		s.lruLock.Unlock()
	}
	// 获取此key missArea下规则
	if len(missArea) == 0 {
		return
	}
	if cs, err = s.KeyRules(c, key, missArea); err != nil {
		return
	}
	if len(cs) == 0 {
		return
	}
	var csMap = make(map[string][]*model.KeyAreaInfo, len(cs))
	for _, r := range cs {
		csMap[r.Area] = append(csMap[r.Area], r)
	}
	for area, cs := range csMap {
		var (
			regs    = []*model.KeyRegxp{}
			matcher = actriekey.NewMatcher()
			fil     = &model.KeyFilter{}
		)
		for _, c := range cs {
			switch c.Mode {
			case model.RegMode:
				var reg *regexp.Regexp
				if reg, err = regexp.Compile(c.Filter); err != nil {
					log.Errorv(context.Background(), log.KV("err", err.Error()), log.KV("msg", "load key filter regexp err"), log.KV("regexp", c.Filter), log.KV("area", c.Area))
					continue
				}
				re := &model.KeyRegxp{FkID: c.FKID}
				re.Reg = reg
				re.Level = c.Level
				regs = append(regs, re)
			case model.StrMode:
				matcher.Insert(c.Filter, c.FKID, c.Level)
			}
		}
		matcher.Build()
		fil.Regs = regs
		fil.Matcher = matcher
		fils = append(fils, fil)
		// 加入到lru链表中
		s.lruLock.Lock()
		if s.lruList.Len() >= s.lruMax {
			s.lruList.Remove(s.lruList.Tailer())
		}
		s.lruList.Prepend(s.lruKey(key, area), fil)
		s.lruLock.Unlock()
	}
	return
}

// KeyRules .
func (s *Service) KeyRules(c context.Context, key string, areas []string) (rs []*model.KeyAreaInfo, err error) {
	var (
		missAreas []string
		dbRs      []*model.KeyAreaInfo
		cacheRs   []*model.KeyAreaInfo
	)
	if cacheRs, missAreas, err = s.dao.KeyAreaCache(c, key, areas); err != nil {
		log.Error("d.keyAreaCache(%s,%v) error(%s)", key, areas, err)
		return
	}
	if len(cacheRs) > 0 {
		rs = append(rs, cacheRs...)
	}
	if len(missAreas) > 0 {
		if dbRs, err = s.dao.KeyAreas(c, key, missAreas); err != nil {
			log.Error("d.KeyAreaRule(%s,%v) error(%s)", key, areas, err)
			return
		}
		if len(dbRs) > 0 {
			rs = append(rs, dbRs...)
		} else {
			for _, area := range missAreas {
				// md5 bilibili
				dbRs = append(dbRs, &model.KeyAreaInfo{Filter: "130e29f351572e58c49fd4c910d7beb0", Mode: 0, Area: area, Level: 10})
			}
		}
		// add cache
		s.cacheCh.Save(func() {
			ctx := context.Background()
			var rsMap = make(map[string][]*model.KeyAreaInfo, len(dbRs))
			for _, r := range dbRs {
				rsMap[r.Area] = append(rsMap[r.Area], r)
			}
			for area, rs := range rsMap {
				s.dao.SetKeyAreaCache(ctx, key, area, rs)
			}
		})
	}
	return
}

func (s *Service) lruKey(key, area string) string {
	return fmt.Sprintf("%s_%s", key, area)
}

// 弹幕过滤默认key
func (s *Service) biliKey() string {
	return "bilibili"
}
