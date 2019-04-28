package service

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"time"

	"go-common/app/service/live/xroom-feed/internal/model"
	"go-common/library/log"
)

// 保底数据源标识
const _defaultSourceName = "__default_source__"

func (s *Service) getPoolConf() (tidyRs map[string][]*model.RecPoolConf) {
	data := s.dao.GetCacheData()
	tidyRs = make(map[string][]*model.RecPoolConf)
	rs := make([]*model.RecPoolConf, 0)
	if err := json.Unmarshal(data, &rs); err != nil {
		log.Error("[getPoolConf] loadDataTOcache json unmarshal err: %+v; rs: %+v;", err, rs)
		return
	}
	if rs == nil {
		log.Error("[getPoolConf] rs err: %+v; ", rs)
		return
	}

	for _, recConf := range rs {
		recPoolDataMap := s.dao.GetRecPoolByID(context.Background(), []int64{recConf.ID})
		if _, ok := recPoolDataMap[recConf.ID]; !ok {
			continue
		}
		if len(recPoolDataMap[recConf.ID]) == 0 {
			continue
		}
		key := strconv.FormatInt(recConf.ModuleType, 10) + "_" +
			strconv.FormatInt(recConf.Position, 10)
		tidyRs[key] = append(tidyRs[key], recConf)
	}

	for k, recConfSlice := range tidyRs {
		sort.Sort(model.RecPoolSlice(recConfSlice))
		//加入保底数据源
		tidyRs[k] = append(tidyRs[k], &model.RecPoolConf{
			ID:          0,
			Name:        _defaultSourceName,
			Type:        0,
			Rule:        "",
			Priority:    1,
			Percent:     0,
			TruePercent: 0,
			ModuleType:  0,
			Position:    0,
		})

		var rest float64
		rest = 100
		for i, recConf := range tidyRs[k] {
			if rest < 0 {
				log.Error("[getPoolConf] true percent cal err: %+v", tidyRs[k])
				recConf.TruePercent = recConf.Percent
				continue
			}
			if i == 0 {
				recConf.TruePercent = recConf.Percent
				rest = rest - recConf.Percent
				continue
			}
			if recConf.Name == _defaultSourceName {
				//保底逻辑
				recConf.TruePercent = rest
				break
			}
			truePercent := rest * recConf.Percent / 100
			recConf.TruePercent = truePercent
			rest = rest - truePercent
		}
	}

	return
}

// GetPoolConfFromMem get pool conf from mem...
func (s *Service) GetPoolConfFromMem(key string) (res []*model.RecPoolConf) {
	rcl := s.recCache.Load()

	rc, ok := rcl.(map[string][]*model.RecPoolConf)
	var raw map[string][]*model.RecPoolConf
	if !ok {
		raw = s.getPoolConf()
		log.Error("[GetPoolConfFromMem] mem cache miss, data assert err, rcl:%+v, raw: %+v", rcl, raw)
	} else {
		raw = rc
	}
	res, ok = raw[key]
	if !ok {
		log.Error("[GetPoolConfFromMem] key not exist, key: %+v", key)
		return
	}
	return
}

func (s *Service) poolConfProc() {
	for {
		time.Sleep(time.Second * 5)
		s.loadPoolConf()
	}
}
func (s *Service) loadPoolConf() {
	recPoolCache := s.getPoolConf()
	if len(recPoolCache) == 0 {
		log.Info("[loadPoolConf] getPoolConf empty")
	}
	s.recCache.Store(recPoolCache)
}
