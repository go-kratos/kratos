package service

import (
	"context"
	"regexp"
	"time"

	"go-common/app/service/main/filter/model"
	"go-common/app/service/main/filter/model/actriearea"
	"go-common/app/service/main/filter/model/actriekey"
	"go-common/library/ecode"
	"go-common/library/log"
)

// FilterTest 测试敏感词
func (s *Service) FilterTest(c context.Context, areas []string, msg string) (ws, bs []*model.FilterTestResult) {
	ws = make([]*model.FilterTestResult, 0)
	bs = make([]*model.FilterTestResult, 0)
	for _, area := range areas {
		resWs, resBs := s.filterTest(c, area, msg)
		ws = append(ws, resWs...)
		bs = append(bs, resBs...)
	}
	return
}

func (s *Service) filterTest(c context.Context, area, msg string) (resWs, resBs []*model.FilterTestResult) {
	var (
		filters []*model.Filter
	)
	// 1. 白名单匹配
	whiteHits := s.whites.Hits([]byte(msg), area, -1)
	for _, whiteHit := range whiteHits {
		h := &model.FilterTestResult{
			Area: area,
			Rule: &actriearea.Rule{
				Msg:   whiteHit.Rule,
				Level: whiteHit.Level,
				Tpids: whiteHit.TypeIDs,
				Fid:   whiteHit.Fid,
				Mode:  whiteHit.Mode,
			},
		}
		resWs = append(resWs, h)
	}
	if filters = s.filters.GetFiltersByArea(area); len(filters) == 0 {
		return
	}
	// 2. 黑名单匹配
	// 正则匹配
	for _, filter := range filters {
		for _, re := range filter.Regs {
			if re.Reg.MatchString(msg) {
				h := &model.FilterTestResult{
					Area: area,
					Rule: &actriearea.Rule{
						Msg:   re.Reg.String(),
						Level: re.Level,
						Tpids: re.TypeIDs,
						Fid:   re.Fid,
						Mode:  0,
					},
				}
				resBs = append(resBs, h)
			}
		}
	}
	// 字符串匹配
	for _, f := range filters {
		_, rs := f.Matcher.Test(msg)
		for _, r := range rs {
			tr := &model.FilterTestResult{
				Area: area,
				Rule: r,
			}
			resBs = append(resBs, tr)
		}
	}
	return
}

// KeyTest 测试key维度的敏感词
func (s *Service) KeyTest(c context.Context, key, msg string, areas []string) (rs []*model.KeyTestResult, err error) {
	// 构建reg规则和actrie filter
	var (
		regs    = []*model.KeyRegxp{}
		matcher = actriekey.NewMatcher()
		filter  = &model.KeyFilter{}
	)
	if !s.areas.CheckArea(areas) {
		err = ecode.FilterIllegalArea
		return
	}
	// 获取此key当前area下规则
	var ks []*model.KeyAreaInfo
	if ks, err = s.KeyRules(c, key, areas); err != nil {
		return
	}
	rsMap := make(map[int64]*model.KeyTestResult, len(ks))

	for _, k := range ks {
		switch k.Mode {
		case model.RegMode:
			re := &model.KeyRegxp{FkID: k.FKID}
			reg, reErr := regexp.Compile(k.Filter)
			if reErr != nil {
				log.Errorv(context.Background(), log.KV("err", err.Error()), log.KV("msg", "key test regexp err"), log.KV("regexp", k.Filter), log.KV("area", k.Area))
				continue
			}
			re.Reg = reg
			regs = append(regs, re)
		case model.StrMode:
			matcher.Insert(k.Filter, k.FKID, k.Level)
		}

		r := &model.KeyTestResult{
			KeyAreaInfo: k,
		}
		r.Shelve = s.shelve(k.STime.Time(), k.CTime.Time())
		r.TpIDs = []int64{}
		rsMap[r.FKID] = r
	}
	matcher.Build()
	filter.Regs = regs
	filter.Matcher = matcher
	// reg check
	for _, reg := range filter.Regs {
		if reg.Reg.MatchString(msg) {
			if r, ok := rsMap[reg.FkID]; ok {
				rs = append(rs, r)
			}
		}
	}
	// actire check
	fkIDs := filter.Matcher.Test(msg)
	for _, fkID := range fkIDs {
		if r, ok := rsMap[fkID]; ok {
			rs = append(rs, r)
		}
	}
	if len(rs) == 0 {
		rs = []*model.KeyTestResult{}
	}
	return
}

func (s *Service) shelve(stime, etime time.Time) bool {
	if time.Now().Unix() >= stime.Unix() && time.Now().Unix() <= etime.Unix() {
		return true
	}
	return false
}
