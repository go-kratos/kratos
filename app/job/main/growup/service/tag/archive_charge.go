package tag

import (
	"context"

	model "go-common/app/job/main/growup/model/tag"
)

// getAvDailyCharge get av_daily_charge
func (s *Service) getAvDailyCharge(c context.Context, month, query string) (avs []*model.ArchiveCharge, err error) {
	if query == "" {
		return
	}
	from, limit := 0, 2000
	for {
		var av []*model.ArchiveCharge
		av, err = s.dao.GetAvDailyCharge(c, month, query, from, limit)
		if err != nil {
			return
		}
		for _, a := range av {
			if a.IncCharge > 0 && a.IsDeleted == 0 {
				avs = append(avs, a)
			}
		}
		if len(av) < limit {
			break
		}
		from += limit
	}
	return
}

// getCmDailyCharge get column_daily_charge
func (s *Service) getCmDailyCharge(c context.Context, query string) (cms []*model.ArchiveCharge, err error) {
	if query == "" {
		return
	}
	from, limit := 0, 2000
	for {
		var cm []*model.ArchiveCharge
		cm, err = s.dao.GetCmDailyCharge(c, query, from, limit)
		if err != nil {
			return
		}
		for _, c := range cm {
			if c.IncCharge > 0 && c.IsDeleted == 0 {
				cms = append(cms, c)
			}
		}
		if len(cm) < limit {
			break
		}
		from += limit
	}
	return
}

func (s *Service) getBackgroundMusic(c context.Context, limit int64) (bgms []*model.ArchiveCharge, err error) {
	var id int64
	bm := make(map[int64]bool)
	bgms = make([]*model.ArchiveCharge, 0)
	for {
		var bs []*model.ArchiveCharge
		bs, id, err = s.dao.GetBgm(c, id, limit)
		if err != nil {
			return
		}
		if len(bs) == 0 {
			break
		}
		for _, b := range bs {
			if _, ok := bm[b.AID]; !ok {
				bm[b.AID] = true
				bgms = append(bgms, b)
			}
		}
	}
	return
}
