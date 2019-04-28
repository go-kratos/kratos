package service

import (
	"math"
	"time"

	"go-common/app/job/main/up-rating/model"
	xtime "go-common/library/time"
)

// Copy copy data from month to month
func (s *Service) Copy(rch chan []*model.Rating, wch chan []*model.Rating, past map[int64]*model.Past, params *model.RatingParameter) {
	defer close(wch)
	for rs := range rch {
		for _, r := range rs {
			if _, ok := past[r.MID]; !ok {
				continue
			}
			r.MetaCreativityScore = 0
			csr := past[r.MID].MetaCreativityScore
			if csr == 0 {
				r.CreativityScore = 0
			} else {
				r.CreativityScore = int64(math.Min(float64(params.WCS)*math.Log(float64(csr)), float64(params.WCSR)))
			}

			r.MetaInfluenceScore = 0
			isr := past[r.MID].MetaInfluenceScore
			if isr == 0 {
				r.InfluenceScore = 0
			} else {
				r.InfluenceScore = int64(math.Min(float64(params.WIS)*math.Log(float64(isr)), float64(params.WISR)))
			}
			r.Date = xtime.Time(time.Date(r.Date.Time().Year(), r.Date.Time().Month()+1, 1, 0, 0, 0, 0, time.Local).Unix())
			r.MagneticScore = r.CreativityScore + r.InfluenceScore + r.CreditScore
		}
		wch <- rs
	}
}

// CalScore cal rating score
func (s *Service) CalScore(rch chan []*model.BaseInfo,
	wch chan []*model.Rating,
	params *model.RatingParameter,
	past map[int64]*model.Past, date time.Time) {
	defer close(wch)
	for bs := range rch {
		m := make([]*model.Rating, 0)
		for _, b := range bs {
			if !b.Date.Time().Equal(date) {
				continue
			}
			r := &model.Rating{
				MID:                 b.MID,
				TagID:               b.TagID,
				MetaCreativityScore: calCreativetyMetaScore(b, params),
				CreativityScore:     calCreativityScore(b, params, past),
				MetaInfluenceScore:  calInfluenceMetaScore(b, params),
				InfluenceScore:      calInfluenceScore(b, params, past),
				CreditScore:         calCreditScore(b, params, past),
				Date:                b.Date,
			}
			r.MagneticScore = r.CreativityScore + r.InfluenceScore + r.CreditScore
			m = append(m, r)
		}
		wch <- m
	}
}

func calCreativetyMetaScore(b *model.BaseInfo, params *model.RatingParameter) int64 {
	// ps: 当月播放分
	ps := params.WDP*b.PlayIncr + params.WDC*b.CoinIncr
	// ubs: 当月投稿低保分
	ubs := params.WDV * int64(math.Min(float64(b.Avs), float64(params.WMDV)))
	// csm: 当月创作力得分
	csm := ps + ubs
	return csm
}

func calCreativityScore(b *model.BaseInfo, params *model.RatingParameter, past map[int64]*model.Past) int64 {
	csm := calCreativetyMetaScore(b, params)
	// csr: csm + past 创作力原始分
	var csr int64
	if _, ok := past[b.MID]; ok {
		csr = csm + past[b.MID].MetaCreativityScore
	} else {
		csr = csm
	}
	if csr < 1 {
		return 0
	}
	// cs: 创作力总分
	cs := math.Min(float64(params.WCS)*math.Log(float64(csr)), float64(params.WCSR))
	return int64(cs)
}

func calInfluenceMetaScore(b *model.BaseInfo, params *model.RatingParameter) int64 {
	// mfans: 当月活跃粉丝数
	mfans := params.WMAAFans*(b.MAAFans+b.MAHFans) + params.WMAHFans*b.MAHFans
	return mfans
}

func calInfluenceScore(b *model.BaseInfo, params *model.RatingParameter, past map[int64]*model.Past) int64 {
	mfans := calInfluenceMetaScore(b, params)
	// isr: 影响力原始分
	var isr int64
	if _, ok := past[b.MID]; ok {
		isr = mfans + past[b.MID].MetaInfluenceScore
	} else {
		isr = mfans
	}
	if isr < 1 {
		return 0
	}
	// is: up主影响力分
	is := math.Min(float64(params.WIS)*math.Log(float64(isr)), float64(params.WISR))
	return int64(is)
}

func calCreditScore(b *model.BaseInfo, params *model.RatingParameter, past map[int64]*model.Past) int64 {
	addScore := min(b.OpenAvs*params.HV, params.HVM)
	minusScore := min(b.LockedAvs*params.HL, params.HLM)
	var cs int64
	if _, ok := past[b.MID]; ok {
		cs = past[b.MID].CreditScore + addScore - minusScore
	} else {
		cs = params.HBASE + addScore - minusScore
	}
	if cs < 0 {
		cs = 0
	}
	return min(cs, params.HR)
}

func min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}
