package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// GetStats .
func (s *Service) GetStats(c context.Context) (stats []*model.MemberStat, err error) {
	s.memberCache.RLock()
	defer s.memberCache.RUnlock()

	st := s.memberCache.uptime
	if st.IsZero() {
		st = time.Now().Add(-30 * time.Minute)
	}
	rst, err := s.memberStats(c, st, time.Now())
	if err != nil {
		log.Error("s.MemberStats error(%v)", err)
		err = nil
	}
	for uid, omst := range s.memberCache.ms {
		if nmst, ok := rst[uid]; ok {
			mst := &model.MemberStat{UID: uid}
			mst.DispatchCount = omst.DispatchCount + nmst.DispatchCount
			mst.ReleaseCount = omst.ReleaseCount + nmst.ReleaseCount
			mst.SubmitCount = omst.SubmitCount + nmst.SubmitCount
			mst.OSubmitCount = omst.OSubmitCount + nmst.OSubmitCount
			mst.NSubmitCount = omst.NSubmitCount + nmst.NSubmitCount
			mst.BelongCount = omst.BelongCount + nmst.BelongCount
			mst.PassCount = omst.PassCount + nmst.PassCount
			mst.NormalCount = omst.NormalCount + nmst.NormalCount
			mst.SubjectCount = omst.SubjectCount + nmst.SubjectCount

			mst.SumDu = omst.SumDu + nmst.SumDu
			mst.SumDuration = fmt.Sprintf("%.2d:%.2d:%.2d", mst.SumDu/3600, (mst.SumDu%3600)/60, (mst.SumDu%3600)%60)

			var CompleteRate, PassRate = 100.0, 100.0
			if mst.SubmitCount < mst.BelongCount {
				CompleteRate = float64(mst.SubmitCount) / float64(mst.BelongCount) * 100.0
			}
			if mst.PassCount < mst.SubmitCount {
				PassRate = float64(mst.PassCount) / float64(mst.SubmitCount) * 100.0
			}
			mst.CompleteRate = fmt.Sprintf("%.2f%%", CompleteRate)
			mst.PassRate = fmt.Sprintf("%.2f%%", PassRate)

			if mst.NSubmitCount == 0 {
				mst.AvgUtime = "00:00:00"
			} else {
				mst.AvgUt = (omst.AvgUt*float64(omst.NSubmitCount) + nmst.AvgUt*float64(nmst.NSubmitCount)) / float64(mst.NSubmitCount)
				mst.AvgUtime = fmt.Sprintf("%.2d:%.2d:%.2d", int64(mst.AvgUt)/3600, (int64(mst.AvgUt)%3600)/60, (int64(mst.AvgUt)%3600)%60)
			}
			delete(rst, uid)

			stats = append(stats, mst)
		} else {
			stats = append(stats, omst)
		}
	}
	if len(rst) > 0 {
		for _, nmst := range rst {
			stats = append(stats, nmst)
		}
	}

	if len(stats) > 0 {
		wg, ctx := errgroup.WithContext(c)
		wg.Go(func() error {
			if err := s.mulIDtoName(ctx, stats, s.lastInTime, "UID", "InTime"); err != nil {
				log.Error("mulIDtoName s.lastInTime error(%v)", err)
			}
			return nil
		})
		wg.Go(func() error {
			if err := s.mulIDtoName(ctx, stats, s.lastOutTime, "UID", "QuitTime"); err != nil {
				log.Error("mulIDtoName s.lastOutTime error(%v)", err)
			}
			return nil
		})
		wg.Wait()
	}
	for _, st := range stats {
		if st.QuitTime <= st.InTime {
			st.QuitTime = ""
		}
	}

	return
}

// MemberStats 审核人员统计数据[通过旧一审提交。会导致同一个任务被多次完成]
func (s *Service) memberStats(c context.Context, st, et time.Time) (stats map[int64]*model.MemberStat, err error) {
	var (
		mx   sync.Mutex
		uids []int64
	)

	stats = make(map[int64]*model.MemberStat)

	if uids, err = s.dao.ActiveUids(c, st, et); err != nil || len(uids) == 0 {
		return
	}
	log.Info("MemberStats,st(%s),et(%s) count(%d) uids(%v) ", st.String(), et.String(), len(uids), uids)

	wg := errgroup.Group{}
	for _, uid := range uids {
		id := uid
		wg.Go(func() error {
			st, e := s.singleStat(context.TODO(), id, st, et)
			if e != nil || st == nil {
				log.Error("s.singleStat(%d) error(%v)", id, e)
				return nil
			}
			mx.Lock()
			stats[st.UID] = st
			mx.Unlock()
			return nil
		})
	}

	err = wg.Wait()
	return
}

func (s *Service) singleStat(c context.Context, uid int64, stime, etime time.Time) (stat *model.MemberStat, err error) {
	var (
		SumDuration int64
		AvgUtime    float64
	)
	stat = &model.MemberStat{UID: uid}

	mapAction, err := s.dao.ActionCountByUID(c, uid, stime, etime)
	if err != nil {
		return
	}
	stat.OSubmitCount = mapAction[model.ActionOldSubmit]
	stat.NSubmitCount = mapAction[model.ActionSubmit]
	stat.SubmitCount = stat.OSubmitCount + stat.NSubmitCount
	stat.DispatchCount = mapAction[model.ActionDispatch]
	stat.ReleaseCount = mapAction[model.ActionRelease]

	if stat.PassCount, err = s.dao.PassCountByUID(c, uid, stime, etime); err != nil {
		return
	}
	if stat.SubjectCount, err = s.dao.SubjectCountByUID(c, uid, stime, etime); err != nil {
		return
	}
	if SumDuration, err = s.dao.SumDurationByUID(c, uid, stime, etime); err != nil {
		return
	}
	if AvgUtime, err = s.dao.AvgUtimeByUID(c, uid, stime, etime); err != nil {
		return
	}

	stat.BelongCount = stat.DispatchCount - stat.ReleaseCount

	var CompleteRate, PassRate = 100.0, 100.0
	if stat.SubmitCount < stat.BelongCount {
		CompleteRate = float64(stat.SubmitCount) / float64(stat.BelongCount) * 100.0
	}
	if stat.PassCount < stat.SubmitCount {
		PassRate = float64(stat.PassCount) / float64(stat.SubmitCount) * 100.0
	}
	stat.CompleteRate = fmt.Sprintf("%.2f%%", CompleteRate)
	stat.PassRate = fmt.Sprintf("%.2f%%", PassRate)

	stat.NormalCount = stat.SubmitCount - stat.SubjectCount

	stat.SumDu = SumDuration
	stat.AvgUt = AvgUtime
	stat.SumDuration = fmt.Sprintf("%.2d:%.2d:%.2d", SumDuration/3600, (SumDuration%3600)/60, (SumDuration%3600)%60)
	stat.AvgUtime = fmt.Sprintf("%.2d:%.2d:%.2d", int64(AvgUtime)/3600, (int64(AvgUtime)%3600)/60, (int64(AvgUtime)%3600)%60)
	return
}

func (s *Service) memberproc() {
	for {
		st := time.Now()
		stats, err := s.memberStats(context.TODO(), st.Add(-24*time.Hour), st)
		if err != nil {
			log.Error("s.MemberStats error(%v)", err)
		} else {
			s.memberCache.Lock()
			s.memberCache.uptime = st
			s.memberCache.ms = stats
			s.memberCache.Unlock()
		}
		log.Info("s.MemberStats ut(%.2f)", time.Since(st).Seconds())
		time.Sleep(30 * time.Minute)
	}
}
