package service

import (
	"context"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// UpMcnDataSummaryCron .
func (s *Service) UpMcnDataSummaryCron() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
	}()
	var (
		err  error
		sids []int64
		msid map[int64]int64
		mmc  map[int64]int64
		mup  map[int64][]int64
		c    = context.TODO()
	)
	if msid, sids, err = s.dao.McnSignMids(c); err != nil {
		log.Error("s.dao.McnSignMids error(%+v)", err)
		return
	}
	if len(sids) == 0 {
		log.Warn("mcn sign data summary empty!")
		return
	}
	if mmc, err = s.dao.McnUPCount(c, sids); err != nil {
		log.Error("s.dao.McnUPCount(%s) error(%+v)", xstr.JoinInts(sids), err)
		return
	}
	if mup, err = s.dao.McnUPMids(c, sids); err != nil {
		log.Error("s.dao.McnUPMids(%s) error(%+v)", xstr.JoinInts(sids), err)
		return
	}
	for sid, smid := range msid {
		var (
			upOK, upMidOK bool
			upNums        int64
			upMids        []int64
			totalFans     int64
			now           = time.Now()
			gDate         = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
		)
		if upNums, upOK = mmc[sid]; !upOK {
			upNums = 0
		}
		if upMids, upMidOK = mup[sid]; upMidOK {
			if len(upMids) == 0 {
				totalFans = 0
			} else {
				if totalFans, err = s.dao.CrmUpMidsSum(c, upMids); err != nil {
					log.Error("s.dao.CrmUpMidsSum(%s) error(%+v)", xstr.JoinInts(upMids), err)
					err = nil
					totalFans = 0
				}
			}
		} else {
			totalFans = 0
		}
		if err = s.dao.AddMcnDataSummary(c, smid, sid, upNums, totalFans, xtime.Time(gDate.Unix())); err != nil {
			log.Error("s.dao.UpMcnUpStateOP(%d,%d,%d,%d,%+v) error(%+v)", smid, sid, upNums, totalFans, xtime.Time(gDate.Unix()), err)
			continue
		}
		log.Info("mcnMid(%d) signID(%d) upNum(%d) totalFans(%d) date(%+v) add data summary table", smid, sid, upNums, totalFans, xtime.Time(gDate.Unix()))
	}
}
