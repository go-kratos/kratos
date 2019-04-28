package service

import (
	"context"
	"time"

	"go-common/app/job/main/credit-timer/model"
	"go-common/library/log"
)

func (s *Service) caseProc(ctx context.Context) {
	affect, err := s.dao.UpdateCaseEndTime(ctx, time.Now())
	log.Info("UpdateCaseEndTime affect %d err %v", affect, err)
	if s.c.Judge.CaseEndVoteTotal > 0 {
		affect, err = s.dao.UpdateCaseEndVote(ctx, s.c.Judge.CaseEndVoteTotal, time.Now().Add(time.Duration(s.c.Judge.ReservedTime)))
		log.Info("UpdateCaseEndVote affect %d err %v", affect, err)
	}
}

func (s *Service) loadConf(ctx context.Context) {
	vTotal, err := s.dao.LoadConf(ctx)
	if err != nil {
		log.Error("loadConf error(%v)", err)
		return
	}
	s.c.Judge.CaseEndVoteTotal = vTotal
}

func (s *Service) juryProc(c context.Context) {
	affect, err := s.dao.UpdateJury(c, time.Now())
	log.Info("update jury affect %d err %v", affect, err)
}

func (s *Service) voteProc(c context.Context) {
	affect, err := s.dao.UpdateVote(c, time.Now())
	log.Info("update vote affect %d err %v", affect, err)
}

// ComputePoint compute KPI point.
func (s *Service) ComputePoint(c context.Context, mid int64) (r model.KpiPoint, err error) {
	var (
		voteTotal, voteRight, blockedTotal, activeDays, opinionNums, opinionQuality int64
		voteRightViolate, voteRightLegal, likes, hates                              int64
		//vr:投准率 vf:投准率系数 af:活跃系数 bf:违规系数 of:观点数量系数 oqf:观点质量系数
		vr, vf, af, bf, of, oqf float64
		point                   float64
		begin, end              string
	)
	begin = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	end = time.Now().Format("2006-01-02")
	if blockedTotal, err = s.dao.CountBlocked(c, mid, begin, end); err != nil {
		return
	}
	if blockedTotal == 0 {
		bf = float64(1)
	}
	if voteTotal, err = s.dao.CountVoteTotal(c, mid, begin, end); err != nil {
		return
	}
	if voteRightViolate, err = s.dao.CountVoteRightViolate(c, mid, begin, end); err != nil {
		return
	}
	if voteRightLegal, err = s.dao.CountVoteRightLegal(c, mid, begin, end); err != nil {
		return
	}
	voteRight = voteRightViolate + voteRightLegal
	if voteTotal > 0 {
		vr = float64(voteRight) / float64(voteTotal)
		vf = s.voteRightRatio(vr)
	}
	if activeDays, err = s.dao.CountVoteActive(c, mid, begin, end); err != nil {
		return
	}
	af = s.activeDaysRatio(activeDays)
	if opinionNums, err = s.dao.CountOpinion(c, mid, begin, end); err != nil {
		return
	}
	of = s.opinionNumsRatio(opinionNums)
	if likes, hates, err = s.dao.OpinionQuality(c, mid, begin, end); err != nil {
		return
	}
	opinionQuality = likes - hates
	oqf = s.opinionQualityRatio(opinionQuality)
	log.Info("mid:%d voteTotal:%v vr:%v vf:%v af:%v bf:%v of:%v oqf:%v", mid, voteTotal, vr, vf, af, bf, of, oqf)
	point = float64(voteTotal) * vr * vf * af * bf * of * oqf * 100
	r.Point = int64(point)
	r.ActiveDays = activeDays
	r.Day = time.Now()
	r.Mid = mid
	r.VoteTotal = voteTotal
	r.VoteRadio = int64(vr * 100)
	r.BlockedTotal = blockedTotal
	r.OpinionNum = opinionNums
	r.OpinionLikes = likes
	r.OpinionHates = hates
	return
}

func (s *Service) kpiPointProc(c context.Context) (err error) {
	var (
		mids []int64
		mid  int64
		r    model.KpiPoint
	)
	if mids, err = s.dao.JuryList(c); err != nil {
		log.Error("kpiPoint err(%v)", err)
		return
	}

	for _, mid = range mids {
		if r, err = s.ComputePoint(c, mid); err != nil {
			log.Error("computePoint err(%v)", err)
			continue
		}
		s.dao.UpdateKPIPoint(c, &r)
	}
	return
}

// KPIProc KPI process.
func (s *Service) KPIProc(c context.Context) (err error) {
	var (
		ps  []model.Kpi
		res []model.KpiPoint
		k   model.KpiPoint
		kd  model.KpiData
		at  int64
		day string
		m   int64
	)
	day = time.Now().Format("2006-01-02")
	if res, err = s.dao.KPIPointDay(c, day); err != nil {
		log.Error("kpiPoint(%s) err(%v)", day, err)
	}
	for _, k = range res {
		a := model.Kpi{}
		a.Point = k.Point
		if len(ps) == 0 {
			a.PreCount = 0
			m = 1
			ps = append(ps, a)
			continue
		}
		if ps[len(ps)-1].Point == k.Point {
			m = m + 1
			continue
		}
		a.PreCount = m
		ps = append(ps, a)
		m = m + 1
	}
	at = int64(len(res))
	d := time.Now()
	for _, k = range res {
		if k.Expired.Format("2006-01-02") != day {
			log.Info("Expired(%s)!=day(%s)", k.Expired.Format("2006-01-02"), day)
			continue
		}
		for i, r := range ps {
			log.Info("RankPer r(%+v) k(%+v)", r, k)
			if r.Point == k.Point {
				b := model.Kpi{}
				b.Point = k.Point
				b.Rank = int64(i + 1)
				b.RankPer = (r.PreCount + 1) * 100 / at
				b.RankTotal = at
				p := b.RankPer
				if p == 0 {
					b.Rate = 1
					b.RankPer = 1
				} else if p > 0 && p <= 10 {
					b.Rate = 1
				} else if p > 10 && p <= 25 {
					b.Rate = 2
				} else if p > 25 && p <= 40 {
					b.Rate = 3
				} else if p > 40 && p <= 60 {
					b.Rate = 4
				} else if p > 60 && p <= 100 {
					b.Rate = 5
				}
				if r.Point == 0 {
					b.Rate = 5
				}
				b.Day = d
				b.Mid = k.Mid
				if err = s.dao.UpdateKPI(c, &b); err != nil {
					log.Error("dao.UpdateKPI(%+v) err(%v)", b, err)
				}
				kd.KpiPoint = k
				end := k.Day
				begin := k.Day.AddDate(0, 0, -30)
				count, _ := s.dao.CountVoteByTime(c, k.Mid, begin, end)
				kd.VoteRealTotal = count
				if err = s.dao.UpdateKPIData(c, &kd); err != nil {
					log.Error("dao.UpdateKPIData(%+v) err(%v)", kd, err)
				}
				break
			}
		}
	}
	return
}

// FixKPI fix kpi.
func (s *Service) FixKPI(c context.Context, year, month, dd int, mid int64) (res []model.KpiPoint, err error) {
	var (
		ps  []model.Kpi
		k   model.KpiPoint
		kd  model.KpiData
		at  int64
		day string
		m   int64
	)
	t := time.Date(year, time.Month(month), dd, 0, 0, 0, 0, time.UTC)
	day = t.Format("2006-01-02")
	if res, err = s.dao.KPIPointDay(c, day); err != nil {
		log.Error("kpiPoint(%s) err(%v)", day, err)
	}
	for _, k = range res {
		a := model.Kpi{}
		a.Point = k.Point
		if len(ps) == 0 {
			a.PreCount = 0
			m = 1
			ps = append(ps, a)
			continue
		}
		if ps[len(ps)-1].Point == k.Point {
			m = m + 1
			continue
		}
		a.PreCount = m
		ps = append(ps, a)
		m = m + 1
	}
	at = int64(len(res))
	for _, k = range res {
		for i, r := range ps {
			if r.Point == k.Point {
				b := model.Kpi{}
				b.Point = k.Point
				b.Rank = int64(i + 1)
				b.RankPer = (r.PreCount + 1) * 100 / at
				b.RankTotal = at
				p := b.RankPer
				if p == 0 {
					b.Rate = 1
					b.RankPer = 1
				} else if p > 0 && p <= 10 {
					b.Rate = 1
				} else if p > 10 && p <= 25 {
					b.Rate = 2
				} else if p > 25 && p <= 40 {
					b.Rate = 3
				} else if p > 40 && p <= 60 {
					b.Rate = 4
				} else if p > 60 && p <= 100 {
					b.Rate = 5
				}
				if r.Point == 0 {
					b.Rate = 5
				}
				b.Day = t
				b.Mid = k.Mid
				if b.Mid == mid {
					log.Info("fix kpi mid %d kpi %v", mid, b)
					if err = s.dao.UpdateKPI(c, &b); err != nil {
						log.Error("dao.UpdateKPI(%+v) err(%v)", b, err)
					}
					kd.KpiPoint = k
					end := k.Day
					begin := k.Day.AddDate(0, 0, -30)
					count, _ := s.dao.CountVoteByTime(c, k.Mid, begin, end)
					kd.VoteRealTotal = count
					if err = s.dao.UpdateKPIData(c, &kd); err != nil {
						log.Error("dao.UpdateKPIData(%+v) err(%v)", kd, err)
					}
				}
				break
			}
		}
	}
	return
}
