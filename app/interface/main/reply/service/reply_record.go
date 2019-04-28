package service

import (
	"context"
	"text/template"

	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/log"
	"go-common/library/xstr"
)

var _emptyRecords = make([]*model.Record, 0)

// Records return reply record from es.
func (s *Service) Records(c context.Context, types []int64, mid, stime, etime int64, order, sort string, pn, ps int32) (res []*model.Record, total int32, err error) {
	var midAts []int64
	if res, total, err = s.search.RecordPaginate(c, types, mid, stime, etime, order, sort, pn, ps); err != nil {
		log.Error("s.search.RecordPaginate(%d,%d,%d,%d,%s,%s) error(%v)", mid, sort, pn, ps, stime, etime, err)
		return
	}
	if res == nil {
		res = _emptyRecords
		return
	}
	for _, r := range res {
		r.Message = template.HTMLEscapeString(r.Message)
		if len(r.Ats) == 0 {
			continue
		}
		var ats []int64
		if ats, err = xstr.SplitInts(r.Ats); err != nil {
			log.Error("xstr.SplitInts(%s) error(%v)", r.Ats, err)
			err = nil
		}
		midAts = append(midAts, ats...)
	}
	if len(midAts) == 0 {
		return
	}
	accMap, _ := s.getAccInfo(c, midAts)
	for _, r := range res {
		r.FillAts(accMap)
	}
	return
}
