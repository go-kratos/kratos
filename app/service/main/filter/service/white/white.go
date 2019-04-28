package white

import (
	"context"
	"strings"
	"sync"

	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/model/actriearea"
	"go-common/app/service/main/filter/service/area"
	"go-common/app/service/main/filter/service/regexp"
	"go-common/library/log"
)

type WhiteAreas func(c context.Context, area string) (fs []*model.WhiteAreaInfo, err error)

type White struct {
	whiteMap map[string]*model.Filter
	sync.RWMutex
}

func New() (white *White) {
	white = &White{
		whiteMap: make(map[string]*model.Filter),
	}
	return
}

func (w *White) Load(ctx context.Context, loader WhiteAreas, area *area.Area) (err error) {
	var (
		whiteMap  = make(map[string]*model.Filter)
		areaNames = area.AreaNames()
	)
	for _, areaName := range areaNames {
		var (
			regs    []*model.Regexp
			matcher = actriearea.NewMatcher()
			rules   []*model.WhiteAreaInfo
		)
		if rules, err = loader(ctx, areaName); err != nil {
			return
		}
		for _, rule := range rules {
			switch rule.Mode {
			case model.RegMode:
				reg, err := regexp.Compile(rule.Content)
				if err != nil {
					log.Errorv(context.TODO(), log.KV("err", err.Error()), log.KV("msg", "load white regexp err"), log.KV("regexp", rule.Content), log.KV("area", rule.Area))
					continue
				} else {
					re := &model.Regexp{Reg: reg, TypeIDs: rule.TpIDs, Fid: rule.ID}
					regs = append(regs, re)
				}
			case model.StrMode:
				matcher.Insert(strings.ToLower(rule.Content), 0, rule.TpIDs, rule.ID)
			}
		}
		whiteMap[areaName] = &model.Filter{}
		whiteMap[areaName].Regs = regs
		matcher.Build()
		whiteMap[areaName].Matcher = matcher
	}
	w.Lock()
	w.whiteMap = whiteMap
	w.Unlock()
	return
}

// Hits 获得白名单 tpID < 0 == all , tpID = 0 全站分区 , tpID > 0 子分区
func (w *White) Hits(byteMsg []byte, area string, tpID int64) (hits []*actriearea.MatchHits) {
	var (
		filter *model.Filter
		ok     bool
	)
	w.RLock()
	if filter, ok = w.whiteMap[area]; !ok {
		w.RUnlock()
		return
	}
	w.RUnlock()
	// reg hit
	var reHits []*actriearea.MatchHits
	for _, re := range filter.Regs {
		if tpID < 0 {
			positions := re.Reg.FindAllIndex(byteMsg)
			if len(positions) > 0 {
				hit := &actriearea.MatchHits{Positions: positions, Level: re.Level}
				hit.TypeIDs = re.TypeIDs
				hit.Rule = re.Reg.String()
				hit.Fid = re.Fid
				hit.Mode = 0
				reHits = append(reHits, hit)
			}
		} else {
			for _, tid := range re.TypeIDs {
				if tid == 0 || tid == tpID {
					positions := re.Reg.FindAllIndex(byteMsg)
					if len(positions) > 0 {
						hit := &actriearea.MatchHits{Positions: positions, Level: re.Level}
						hit.TypeIDs = re.TypeIDs
						hit.Rule = re.Reg.String()
						hit.Fid = re.Fid
						hit.Mode = 0
						reHits = append(reHits, hit)
					}
				}
			}
		}
	}
	hits = append(hits, reHits...)

	// actire hit
	acHits := filter.Matcher.GetWhiteHits(string(byteMsg), tpID)
	hits = append(hits, acHits...)
	return
}
