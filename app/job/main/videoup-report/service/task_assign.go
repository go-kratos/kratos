package service

import (
	"context"
	"math/rand"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/task"
	"go-common/library/log"
	"go-common/library/xstr"
)

// 设置指派任务
func (s *Service) setTaskAssign(c context.Context, a *archive.Archive, v *archive.Video) (t *task.Task) {
	var (
		now    = time.Now()
		arruid = []int64{}
		mapUID = make(map[int64]int64)
	)

	t = &task.Task{
		Pool:    task.PoolForFirst,
		Aid:     a.ID,
		Cid:     v.Cid,
		Subject: task.SubjectForNormal,
		AdminID: int64(0),
		UID:     int64(0),
		State:   task.StateForTaskDefault,
	}

	for _, tc := range s.assignCache {
		log.Info("task doing(%v) aid(%d) type(%d) filename(%s) tc(%d)", t, a.ID, a.TypeID, v.Filename, tc.ID)
		if tc.STime.After(now) || tc.ETime.Before(now) {
			log.Error("task time is error stime(%v) etime(%v)", tc.STime, tc.ETime)
			continue
		}
		var midOk, tidOk, durationOk = true, true, true
		if len(tc.MIDs) > 0 {
			if _, midOk = tc.MIDs[a.Mid]; !midOk {
				log.Info("task mid(%d) wrong", a.Mid)
			}
		}
		if len(tc.TIDs) > 0 {
			if _, tidOk = tc.TIDs[a.TypeID]; !tidOk {
				log.Info("task type(%d) wrong", a.TypeID)
			}
		}
		if tc.MinDuration != tc.MaxDuration && (v.Duration < tc.MinDuration || v.Duration > tc.MaxDuration) {
			log.Error("task minDur(%d) maxDur(%d) wrong", tc.MinDuration, tc.MaxDuration)
			durationOk = false
		}
		if midOk && tidOk && durationOk {
			for _, uid := range tc.UIDs {
				if _, ok := mapUID[uid]; !ok {
					mapUID[uid] = tc.AdminID
					arruid = append(arruid, uid)
				}
			}
		}
	}
	if len(arruid) > 0 {
		uids, err := s.arc.ConsumerOnline(c, xstr.JoinInts(arruid))
		if err != nil || len(uids) == 0 {
			log.Warn("task s.arc.ConsumerOnline(%v) (%v) err(%v)", arruid, uids, err)
			return
		}
		if len(uids) == 1 {
			t.UID = uids[0]
		} else {
			inx := rand.Intn(len(uids) - 1)
			log.Info("task uids(%v) rand inx(%d)", uids, inx)
			t.UID = uids[inx]
		}
	}

	if t.UID != 0 {
		t.Subject = task.SubjectForTask
		t.AdminID = mapUID[t.UID] // 命中多个指派者配置，选择其中一个就行
		t.State = task.StateForTaskDefault
	}

	return
}

// 指派配置
func (s *Service) assignConf(c context.Context) (tcs map[int64]*task.AssignConfig, err error) {
	var ids []int64
	if tcs, err = s.arc.AssignConfigs(context.TODO()); err != nil {
		log.Error("s.arc.AssignConfigs(%v) error(%v)", err)
		return
	}

	for k, v := range tcs {
		if !v.ETime.IsZero() && v.ETime.Before(time.Now()) {
			delete(tcs, k)
			ids = append(ids, k)
		}
	}
	if len(ids) > 0 {
		log.Info("task config(%v) 指派配置已过期,自动失效", ids)
		s.arc.DelAssignConfs(c, ids)
	}
	return
}
