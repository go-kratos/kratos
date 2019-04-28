package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/up-rating/model"

	"go-common/library/log"

	"golang.org/x/sync/errgroup"
)

// CalTrend cal trend
func (s *Service) CalTrend(c context.Context, date time.Time) (err error) {
	return s.calTrend(c, date)
}

func (s *Service) calTrend(c context.Context, date time.Time) (err error) {
	ds, err := s.getDiffs(c, date)
	if err != nil {
		return
	}

	var (
		magneticSec  = sections(TotalType)
		creativeSec  = sections(CreativeType)
		influenceSec = sections(InfluenceType)
		creditSec    = sections(CreditType)
		g            errgroup.Group

		// for concurrent write
		_ascMagneticCh  = make(chan []*model.Diff, 5)
		_ascCreativeCh  = make(chan []*model.Diff, 5)
		_ascInfluenceCh = make(chan []*model.Diff, 5)
		_ascCreditCh    = make(chan []*model.Diff, 5)

		_descMagneticCh  = make(chan []*model.Diff, 5)
		_descCreativeCh  = make(chan []*model.Diff, 5)
		_descInfluenceCh = make(chan []*model.Diff, 5)
		_descCreditCh    = make(chan []*model.Diff, 5)
	)

	// magnetic ascend
	g.Go(func() (err error) {
		defer close(_ascMagneticCh)
		mr := classify(ds, TotalType, "asc", magneticSec)
		push(_ascMagneticCh, mr)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "asc", _ascMagneticCh)
		return
	})

	// creativity ascend
	g.Go(func() (err error) {
		defer close(_ascCreativeCh)
		cr := classify(ds, CreativeType, "asc", creativeSec)
		push(_ascCreativeCh, cr)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "asc", _ascCreativeCh)
		return
	})

	// influence ascend
	g.Go(func() (err error) {
		defer close(_ascInfluenceCh)
		ir := classify(ds, InfluenceType, "asc", influenceSec)
		push(_ascInfluenceCh, ir)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "asc", _ascInfluenceCh)
		return
	})

	// credit ascend
	g.Go(func() (err error) {
		defer close(_ascCreditCh)
		cr := classify(ds, CreditType, "asc", creditSec)
		push(_ascCreditCh, cr)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "asc", _ascCreditCh)
		return
	})

	// magnetic descend
	g.Go(func() (err error) {
		defer close(_descMagneticCh)
		mr := classify(ds, TotalType, "desc", magneticSec)
		push(_descMagneticCh, mr)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "desc", _descMagneticCh)
		return
	})

	// creativity descend
	g.Go(func() (err error) {
		defer close(_descCreativeCh)
		cr := classify(ds, CreativeType, "desc", creativeSec)
		push(_descCreativeCh, cr)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "desc", _descCreativeCh)
		return
	})

	// influence descend
	g.Go(func() (err error) {
		defer close(_descInfluenceCh)
		ir := classify(ds, InfluenceType, "desc", influenceSec)
		push(_descInfluenceCh, ir)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "desc", _descInfluenceCh)
		return
	})

	// credit descend
	g.Go(func() (err error) {
		defer close(_descCreditCh)
		cr := classify(ds, CreditType, "desc", creditSec)
		push(_descCreditCh, cr)
		return
	})

	g.Go(func() (err error) {
		err = s.batchInsertDiffs(c, "desc", _descCreditCh)
		return
	})

	if err = g.Wait(); err != nil {
		log.Error("g.Wait error(%v)", err)
	}
	return
}

func push(ch chan []*model.Diff, m map[int64]map[int]Heap) {
	for _, sm := range m {
		for _, h := range sm {
			ch <- h.Result()
		}
	}
}

func (s *Service) getDiffs(c context.Context, date time.Time) (ds map[int64][]*model.Diff, err error) {
	lastMonth := time.Date(date.Year(), date.Month()-1, 1, 0, 0, 0, 0, time.Local)
	last, err := s.getR(c, lastMonth)
	if err != nil {
		return
	}

	curMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	cur, err := s.getR(c, curMonth)
	if err != nil {
		return
	}
	ds = diff(last, cur)
	return
}

func sections(ctype int) (ss []*section) {
	switch ctype {
	case TotalType:
		ss = initSec(10, 600)
	case CreativeType, InfluenceType, CreditType:
		ss = initSec(10, 200)
	}
	return
}

type section struct {
	start int64
	end   int64
	tips  string
}

func initSec(c, total int64) (ss []*section) {
	size := total / c
	for index := int64(0); index < c; index++ {
		start := index * size
		end := start + size - 1
		if index == c-1 {
			end = total
		}
		sec := &section{
			start: start,
			end:   end,
			tips:  fmt.Sprintf("\"%d-%d\"", start, end),
		}
		ss = append(ss, sec)
	}
	return
}

func getSection(score int64, ss []*section) (index int, tips string) {
	for index, section := range ss {
		if score >= section.start && score <= section.end {
			return index, section.tips
		}
	}
	return
}

func getHeap(order string, ctype int) Heap {
	if order == "asc" {
		return &AscHeap{heap: make([]*model.Diff, 0), ctype: ctype}
	}
	return &DescHeap{heap: make([]*model.Diff, 0), ctype: ctype}
}

// map[tag_id]map[section][]diffs
func classify(ds map[int64][]*model.Diff, ctype int, order string, ss []*section) (m map[int64]map[int]Heap) {
	m = make(map[int64]map[int]Heap)
	for tagID, diffs := range ds {
		for _, diff := range diffs {
			sec, tips := getSection(diff.GetScore(ctype), ss)
			// need clone
			_diff := clone(diff)
			_diff.Section = sec
			_diff.CType = ctype
			_diff.Tips = tips
			if sm, ok := m[tagID]; ok {
				if _, ok = sm[sec]; ok {
					sm[sec].Put(_diff)
				} else {
					h := getHeap(order, ctype)
					h.Put(_diff)
					sm[sec] = h
				}
			} else {
				h := getHeap(order, ctype)
				h.Put(_diff)
				sm := make(map[int]Heap)
				sm[sec] = h
				m[tagID] = sm
			}
		}
	}
	return
}

func clone(a *model.Diff) *model.Diff {
	return &model.Diff{
		MID:             a.MID,
		MagneticScore:   a.MagneticScore,
		CreativityScore: a.CreativityScore,
		InfluenceScore:  a.InfluenceScore,
		CreditScore:     a.CreditScore,
		MagneticDiff:    a.MagneticDiff,
		CreativityDiff:  a.CreativityDiff,
		InfluenceDiff:   a.InfluenceDiff,
		CreditDiff:      a.CreditDiff,
		TotalAvs:        a.TotalAvs,
		Fans:            a.Fans,
		TagID:           a.TagID,
		Date:            a.Date,
	}
}

// ds map[tag_id][]*model.Diff
func diff(last map[int64]*model.Rating, cur map[int64]*model.Rating) (ds map[int64][]*model.Diff) {
	ds = make(map[int64][]*model.Diff)
	for mid, r := range cur {
		if _, ok := last[mid]; !ok {
			continue
		}
		lr := last[mid]
		diff := &model.Diff{
			MID:             mid,
			TagID:           r.TagID,
			MagneticScore:   r.MagneticScore,
			CreativityScore: r.CreativityScore,
			InfluenceScore:  r.InfluenceScore,
			CreditScore:     r.CreditScore,
			MagneticDiff:    int(r.MagneticScore - lr.MagneticScore),
			CreativityDiff:  int(r.CreativityScore - lr.CreativityScore),
			InfluenceDiff:   int(r.InfluenceScore - lr.InfluenceScore),
			CreditDiff:      int(r.CreditScore - lr.CreditScore),
			Date:            r.Date,
		}
		if _, ok := ds[diff.TagID]; ok {
			ds[diff.TagID] = append(ds[diff.TagID], diff)
		} else {
			ds[diff.TagID] = []*model.Diff{diff}
		}
	}
	return
}

// GetR m[mid]*model.Rating
func (s *Service) getR(c context.Context, date time.Time) (m map[int64]*model.Rating, err error) {
	var (
		g        errgroup.Group
		routines = 5 // 4 core + 1
		ch       = make(chan []*model.Rating, routines)
	)

	m = make(map[int64]*model.Rating)
	offset, end, total, err := s.RatingOffEnd(c, date)
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	t := end - offset
	section := (t - t%routines) / routines
	for i := 0; i < routines; i++ {
		begin := section*i + offset
		over := begin + section
		if i == routines-1 {
			over = end
		}
		g.Go(func() (err error) {
			err = s.Ratings(c, date, begin, over, ch)
			if err != nil {
				log.Error("get rating infos error(%v)", err)
			}
			return
		})
	}

	g.Go(func() (err error) {
	Loop:
		for rs := range ch {
			for _, r := range rs {
				m[r.MID] = r
			}
			if len(m) == total {
				break Loop
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		log.Error("get rating wait error(%v)", err)
	}
	return
}

// BatchInsertDiffs batch insert diffs
func (s *Service) batchInsertDiffs(c context.Context, table string, wch chan []*model.Diff) (err error) {
	var (
		buff    = make([]*model.Diff, _limit)
		buffEnd = 0
	)
	for ds := range wch {
		for _, d := range ds {
			buff[buffEnd] = d
			buffEnd++
			if buffEnd >= _limit {
				values := diffValues(buff[:buffEnd])
				buffEnd = 0
				_, err = s.dao.InsertTrend(c, table, values)
				if err != nil {
					return
				}
			}
		}
		if buffEnd > 0 {
			values := diffValues(buff[:buffEnd])
			buffEnd = 0
			_, err = s.dao.InsertTrend(c, table, values)
		}
	}
	return
}

func diffValues(ds []*model.Diff) (values string) {
	var buf bytes.Buffer
	for _, r := range ds {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(r.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.CreativityScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(r.CreativityDiff))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.InfluenceScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(r.InfluenceDiff))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.CreditScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(r.CreditDiff))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(r.MagneticScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(r.MagneticDiff))
		buf.WriteByte(',')
		buf.WriteString(fmt.Sprintf("'%s'", r.Date.Time().Format(_layout)))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(r.CType))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(r.Section))
		buf.WriteByte(',')
		buf.WriteString(fmt.Sprintf("'%s'", r.Tips))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}

// DelTrends del trends
func (s *Service) DelTrends(c context.Context, table string) (err error) {
	for {
		var rows int64
		rows, err = s.dao.DelTrend(c, table, _limit)
		if err != nil {
			return
		}
		if rows == 0 {
			break
		}
	}
	return
}
