package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"sync"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/model/actriearea"
	rpcmodel "go-common/app/service/main/filter/model/rpc"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Filter 单个过滤
func (s *Service) Filter(c context.Context, area, msg string, tpid int64, id int64, oid int64, mid int64, keys []string, replyType int8) (fmsg string, level int8, tpIDs []int64, hitRules []string, limit int, aiScore *model.AiScore, err error) {
	if msg == "" {
		return
	}
	if mid > 0 {
		if _, ok := s.whiteMids.White(mid); ok {
			fmsg = msg
			log.Info("Filte white mid hit s.conf.Ai.Whites(%d %s)", mid, msg)
			return
		}
	}
	// 1. get result from cache
	var (
		cacheRes     *model.FilterCacheRes
		cacheFlag    = false
		cacheContent string
	)
	aiScore = &model.AiScore{}
	if len(msg) <= conf.Conf.Property.FilterCacheShortMaxSize {
		cacheFlag = true
		cacheContent = msg
		if cacheRes, err = s.dao.FilterCache(c, area, tpid, keys, msg); err != nil {
			log.Error("%+v", err)
			err = nil
			cacheFlag = false
		}
	} else if len(msg) >= conf.Conf.Property.FilterCacheLongMinSize {
		cacheFlag = true
		hash := hashContent(msg)
		cacheContent = hash
		if cacheRes, err = s.dao.FilterCache(c, area, tpid, keys, hash); err != nil {
			log.Error("%+v", err)
			err = nil
			cacheFlag = false
		}
	}
	if cacheRes != nil {
		fmsg, level, tpIDs, hitRules, limit, aiScore = cacheRes.Fmsg, cacheRes.Level, cacheRes.TpIDs, cacheRes.HitRules, cacheRes.Limit, cacheRes.AI
		return
	}

	// 2. check area validation
	fmsg = msg
	a := s.areas.Area(area)
	if a == nil {
		log.Error("s.Filter invalid area [%s]", area)
		return
	}

	var (
		wg sync.WaitGroup
	)
	// 3. 敏感词过滤
	if a.IsFilter() {
		wg.Add(1)
		go func() {
			defer func() {
				if x := recover(); x != nil {
					log.Error("filter proc panic (%+v)", err)
				}
			}()
			defer wg.Done()
			if fmsg, level, tpIDs, hitRules, err = s.filter(context.Background(), area, msg, tpid, a.IsFilterCommon(), a.IsFullLevel()); err != nil {
				return
			}
			// key 维度过滤
			if a.IsFilterKey(keys) {
				if level == 0 || len(hitRules) == 0 {
					level = 0
					for _, key := range append(keys, s.biliKey()) {
						var fils []*model.KeyFilter
						if fils, err = s.fils(c, key, []string{"common", area}); err != nil {
							return
						}
						if fmsg, level, hitRules = s.opFils(fmsg, level, fils); level != 0 {
							break
						}
					}
				}
			}
		}()
	}

	// 4. 反垃圾过滤
	if a.IsFilterRubbish(oid) {
		wg.Add(1)
		go func() {
			defer func() {
				if x := recover(); x != nil {
					log.Error("antispam proc panic (%+v)", err)
				}
			}()
			defer wg.Done()
			var (
				limitType   string
				antispamErr error
			)
			if limitType, _, antispamErr = s.rubbishFilter(context.Background(), a, msg, oid, id, mid); antispamErr != nil {
				log.Error("%+v", antispamErr)
				return
			}
			switch limitType {
			case model.LimitTypeBlack:
				limit = ecode.FilterHitLimitBlack.Code()
				return
			case model.LimitTypeRestrict:
				limit = ecode.FilterHitStrictLimit.Code()
				return
			case model.LimitTypeExceed:
				limit = ecode.FilterHitRubLimit.Code()
				return
			case model.LimitTypeOK:
			default:
				log.Errorv(c,
					log.KV("log", "antispam limittype unknown"),
					log.KV("limittype", limitType),
					log.KV("area", area),
					log.KV("msg", msg),
				)
				return
			}
		}()
	}
	// AI过滤
	if a.IsAIFilter() {
		wg.Add(1)
		go func() {
			defer func() {
				if x := recover(); x != nil {
					log.Error("ai filter proc panic (%+v)", err)
				}
			}()
			defer wg.Done()
			var aiFilterErr error

			switch a.Name {
			case "reply":
				if aiFilterErr = s.FilterAiScore(context.Background(), msg, mid, model.ADMINID, oid, id, replyType); aiFilterErr != nil {
					log.Error("%+v", aiFilterErr)
					return
				}
			default:
				if aiScore, aiFilterErr = s.dao.AiScore(context.Background(), msg, a.Name); aiFilterErr != nil {
					log.Error("%s: %+v", a.Name, aiFilterErr)
					return
				}
			}

		}()
	}
	wg.Wait()
	// 5. 存储hbase
	if conf.Conf.HBase != nil {
		s.replyToHbase(c, level, id, area, msg)
	}
	// 6. save to cache
	if cacheFlag {
		s.cacheCh.Save(func() {
			cacheRes = &model.FilterCacheRes{
				Fmsg:     fmsg,
				Level:    level,
				TpIDs:    tpIDs,
				HitRules: hitRules,
				Limit:    limit,
				AI:       aiScore,
			}
			if theErr := s.dao.SetFilterCache(context.Background(), area, tpid, keys, cacheContent, cacheRes); theErr != nil {
				log.Error("%+v", theErr)
			}
		})
	}
	return
}

// HTTPMultiAreaFilter .
func (s *Service) HTTPMultiAreaFilter(c context.Context, area string, msgs []string, tpid int64) (list []*model.HTTPAreaFilterRes, err error) {
	for _, msg := range msgs {
		data := &model.HTTPAreaFilterRes{}
		data.MSG, data.Level, data.TypeID, err = s.FilterArea(c, area, msg, tpid)
		if err != nil {
			return
		}
		list = append(list, data)
	}
	return
}

// HTTPMultiFilter .
func (s *Service) HTTPMultiFilter(c context.Context, area string, msgs []string, tpid int64, id int64, oid int64, mid int64, keys []string, replyType int8) (list []*model.HTTPFilterRes, err error) {
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	list = make([]*model.HTTPFilterRes, len(msgs))

	for index, msg := range msgs {
		wg.Add(1)
		go func(i int, m string) {
			defer func() {
				if x := recover(); x != nil {
					log.Error("HTTPMultiFilter filter panic (%+v)", x)
				}
			}()
			defer wg.Done()
			var (
				res    = &model.HTTPFilterRes{}
				theErr error
			)
			res.MSG, res.Level, res.TypeID, res.Hit, res.Limit, res.AI, theErr = s.Filter(c, area, m, tpid, id, oid, mid, keys, replyType)
			if theErr != nil {
				mu.Lock()
				err = theErr
				mu.Unlock()
				log.Error("%+v", err)
				return
			}
			mu.Lock()
			list[i] = res
			mu.Unlock()
		}(index, msg)
	}
	wg.Wait()
	return
}

// RPCMultiFilter RPC批量过滤
func (s *Service) RPCMultiFilter(c context.Context, area string, msgMap map[string]string, tpid int64) (rMap map[string]*rpcmodel.FilterRes, err error) {
	rMap = make(map[string]*rpcmodel.FilterRes, len(msgMap))
	if len(msgMap) == 0 {
		return
	}
	a := s.areas.Area(area)
	if a == nil {
		log.Error("s.RPCMFilter invalid area [%s]", area)
		return
	}

	for id, msg := range msgMap {
		var (
			level int8
			fmsg  string
		)
		// check common
		if fmsg, level, _, _, err = s.filter(c, area, msg, tpid, a.IsFilterCommon(), a.IsFullLevel()); err != nil {
			return
		}
		result := &rpcmodel.FilterRes{Result: fmsg, Level: level}
		rMap[id] = result
	}
	return
}

// FilterArea 不带common业务的过滤
func (s *Service) FilterArea(c context.Context, area, msg string, tpid int64) (fmsg string, level int8, tpIDs []int64, err error) {
	if msg == "" {
		return
	}
	a := s.areas.Area(area)
	if a == nil {
		log.Error("s.FilterArea invalid area [%s]", area)
		return
	}
	if fmsg, level, tpIDs, _, err = s.filter(c, area, msg, tpid, false, a.IsFullLevel()); err != nil {
		return
	}
	if tpIDs == nil {
		tpIDs = make([]int64, 0)
	}
	return
}

// MFilterArea 批量不带common业务的过滤
func (s *Service) MFilterArea(c context.Context, area string, msgMap map[string]string, tpid int64) (rMap map[string]*rpcmodel.FilterRes, err error) {
	rMap = make(map[string]*rpcmodel.FilterRes, len(msgMap))
	a := s.areas.Area(area)
	if a == nil {
		log.Error("s.MFilterArea invalid area [%s]", area)
		return
	}
	for id, msg := range msgMap {
		var (
			alevel, level int8
			fmsg          string
		)
		// check area
		if fmsg, alevel, _, _, err = s.filter(c, area, msg, tpid, false, a.IsFullLevel()); err != nil {
			return
		}
		if alevel > level {
			level = alevel
		}
		result := &rpcmodel.FilterRes{Result: fmsg, Level: level}
		rMap[id] = result
	}
	return
}

// 并行过滤处理
func (s *Service) filter(c context.Context, area, msg string, tpID int64, commonHit bool, isFullLevel bool) (fmsg string, level int8, allTpIDs []int64, hitRules []string, err error) {
	hitRules = make([]string, 0, 5)
	allTpIDs = make([]int64, 0, 5)
	var (
		filters []*model.Filter
		allRegs = make([]*model.Regexp, 0, 5000)
	)
	// 第一步：获取白名单命中规则
	whiteHits := s.whites.Hits([]byte(msg), area, tpID)

	// 第二步：获取黑名单命中规则
	if filters = s.filters.GetFilters(area, commonHit); len(filters) == 0 {
		fmsg = msg
		return
	}
	// combine filters' regs
	for _, f := range filters {
		allRegs = append(allRegs, f.Regs...)
	}
	var (
		blackHits []*actriearea.MatchHits

		mutex    sync.Mutex
		paraSize = conf.Conf.Property.ParallelSize
		regLen   = len(allRegs)
		wg       = new(sync.WaitGroup)
		msgBytes = []byte(msg)
	)
	// regex hits
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
					log.Error("service.filter(%s,%s,%d,%b) panic(%v)", area, msg, tpID, commonHit, x)
				}
			}()
			for _, re := range regs {
				for _, id := range re.TypeIDs {
					if id == 0 || id == tpID {
						if re.Level <= conf.Conf.Property.CriticalFilterLevel && !isFullLevel {
							if re.Reg.MatchString(msg) {
								hit := &actriearea.MatchHits{Level: re.Level, Fid: re.Fid, Area: re.Area, Rule: re.Reg.String()}
								hit.TypeIDs = append(hit.TypeIDs, id)
								hit.Rule = re.Reg.String()
								mutex.Lock()
								blackHits = append(blackHits, hit)
								mutex.Unlock()
							}
						} else {
							positions := re.Reg.FindAllIndex(msgBytes)
							if len(positions) > 0 {
								hit := &actriearea.MatchHits{Positions: positions, Level: re.Level, Fid: re.Fid, Area: re.Area, Rule: re.Reg.String()}
								hit.TypeIDs = append(hit.TypeIDs, id)
								hit.Rule = re.Reg.String()
								mutex.Lock()
								blackHits = append(blackHits, hit)
								mutex.Unlock()
							}
						}
						break
					}
				}
			}
		}()
	}
	// actrie hits
	for _, f := range filters {
		if isFullLevel {
			mutex.Lock()
			blackHits = append(blackHits, f.Matcher.Filter(msg, tpID, 0)...)
			mutex.Unlock()
		} else {
			mutex.Lock()
			blackHits = append(blackHits, f.Matcher.Filter(msg, tpID, conf.Conf.Property.CriticalFilterLevel)...)
			mutex.Unlock()
		}
	}
	wg.Wait()
	// 终极过滤  重点是，如果白名单完全覆盖黑名单命中位置，则黑名单不生效
	var (
		sortedBlackHit = &sortBlackHits{}
	)
	for _, blackhit := range blackHits {
		s.areaBlackHitProm.Incr(area)
		log.Infov(c,
			log.KV("log", "filter black hit"),
			log.KV("area", area),
			log.KV("hit", blackhit.Rule),
			log.KV("msg", msg),
		)
		for _, posi := range blackhit.Positions {
			var whiteCover = false
			if len(whiteHits) > 0 {
			OUTER:
				for _, whhit := range whiteHits {
					log.Info("filter white hit (%+v) : msg(%s) area(%s) ", whhit, msg, area)
					for _, whPosi := range whhit.Positions {
						// 如果白名单完全覆盖了黑名单则忽略，否则算黑
						if whPosi[0] <= posi[0] && whPosi[1] >= posi[1] {
							log.Info("filter white cover (%+v) : msg(%s) area(%s) ", whhit, msg, area)
							whiteCover = true
							break OUTER
						}
					}
				}
			}
			if !whiteCover {
				sortedBlackHit.Add(posi)
				if blackhit.Level > level {
					level = blackhit.Level
				}
				allTpIDs = append(allTpIDs, blackhit.TypeIDs...)
				hitRules = append(hitRules, blackhit.Rule)
			}
		}
	}
	//replace '*'
	fmsg, err = sortedBlackHit.CoverByStart(msg)
	s.addEvent(func() {
		s.repostHitLog(context.Background(), area, msg, blackHits, "area")
	})
	return
}

type sortBlackHits struct {
	posi [][]int
}

func (s *sortBlackHits) Len() int {
	return len(s.posi)
}

func (s *sortBlackHits) Less(i, j int) bool {
	return s.posi[i][0] <= s.posi[j][0]
}

func (s *sortBlackHits) Swap(i, j int) {
	s.posi[i], s.posi[j] = s.posi[j], s.posi[i]
}

func (s *sortBlackHits) Add(posi []int) {
	s.posi = append(s.posi, posi)
}

// CoverByStart cover block word by start.
func (s *sortBlackHits) CoverByStart(msg string) (res string, err error) {
	sort.Sort(s)
	var (
		curPos       = 0
		msgBytes     = []byte(msg)
		filterdBytes = &bytes.Buffer{}
	)
	for _, posi := range s.posi {
		if posi[1] < curPos {
			continue
		}
		if posi[0] < curPos {
			posi[0] = curPos
		}
		var maskStr string
		for range msg[posi[0]:posi[1]] {
			maskStr += "*"
		}
		if _, err = filterdBytes.Write(msgBytes[curPos:posi[0]]); err != nil {
			log.Error("byteBuff.Write(%s) error(%v)", string(msgBytes[curPos:posi[0]]), err)
			return
		}
		if _, err = filterdBytes.Write([]byte(maskStr)); err != nil {
			log.Error("byteBuff.Write(%s) error(%v)", maskStr, err)
			return
		}
		curPos = posi[1]
	}
	if _, err = filterdBytes.Write(msgBytes[curPos:]); err != nil {
		log.Error("byteBuff.Write(%s) error(%v)", string(msgBytes[curPos:]), err)
		return
	}
	res = filterdBytes.String()
	return
}

// Rubbish 仅反垃圾过滤,目前仅有私信接入 @2017.12.14
func (s *Service) Rubbish(c context.Context, area, msg string, oid int64) (fmsg string, hits []string, err error) {
	var (
		limitType string
		senderID  int64
		a         *model.Area
	)
	a = s.areas.Area(area)
	if a == nil {
		log.Error("s.Filter invalid area [%s]", area)
		return
	}
	fmsg = msg
	if limitType, hits, err = s.rubbishFilter(c, a, msg, oid, 0, senderID); err != nil {
		return
	}
	switch limitType {
	case model.LimitTypeBlack:
		err = ecode.FilterHitLimitBlack
	case model.LimitTypeRestrict:
		err = ecode.FilterHitStrictLimit
	case model.LimitTypeExceed:
		err = ecode.FilterHitRubLimit
	case model.LimitTypeOK:
	default:
		log.Errorv(c,
			log.KV("log", "antispam limittype unknown"),
			log.KV("limittype", limitType),
			log.KV("area", area),
			log.KV("msg", msg),
		)
	}
	return
}

func hashContent(msg string) (hs string) {
	bytes := md5.Sum([]byte(msg))
	return fmt.Sprintf("%x", bytes)
}
