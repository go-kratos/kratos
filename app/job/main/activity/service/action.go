package service

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/job/main/activity/model/like"
	"go-common/library/log"
)

// actionDealProc .
func (s *Service) actionDealProc(i int) {
	defer s.waiter.Done()
	var (
		ch = s.subActionCh[i]
		sm = s.actionSM[i]
		ls *like.LastTmStat
	)
	for {
		ms, ok := <-ch
		if !ok {
			s.multiUpActDB(i, sm)
			log.Warn("s.actionDealProc(%d) quit", i)
			return
		}
		if ls, ok = sm[ms.Lid]; !ok {
			ls = &like.LastTmStat{Last: time.Now().Unix()}
			sm[ms.Lid] = ls
			// the first time update db.
			s.updateActDB([]int64{ms.Lid})
		}
		if time.Now().Unix()-ls.Last > 60 {
			s.updateActDB([]int64{ms.Lid})
			delete(sm, ms.Lid)
		}
		log.Info("s.actionDealProc(%d) lid:%d time:%d", i, ms.Lid, ls.Last)
	}
}

// updateActDB batch to deal like_extend.
func (s *Service) updateActDB(lids []int64) {
	var (
		c         = context.Background()
		insertExt []*like.Extend
	)
	if len(lids) == 0 {
		return
	}
	lidLike, err := s.dao.BatchLikeActSum(c, lids)
	if err != nil {
		log.Error("s.dao.BatchLikeActSum(%v) error(%+v)", lids, err)
		return
	}
	insertExt = make([]*like.Extend, 0, len(lids))
	for _, v := range lids {
		if _, ok := lidLike[v]; ok {
			insertExt = append(insertExt, &like.Extend{Lid: v, Like: lidLike[v]})
		} else {
			log.Warn("s.updateActDB() data has not found")
		}
	}
	if len(insertExt) == 0 {
		return
	}
	s.BatchInsertLikeExtend(c, insertExt)
}

// multiUpActDB division sm data .
func (s *Service) multiUpActDB(yu int, sm map[int64]*like.LastTmStat) {
	var (
		i         int
		startLids = [1000]int64{}
		lids      = startLids[:0]
	)
	log.Info("start close(%d) multiUpActDB start", yu)
	for lid := range sm {
		lids = append(lids, lid)
		i++
		if i%1000 == 0 {
			s.updateActDB(lids)
			lids = startLids[:0]
		}
	}
	if len(lids) > 0 {
		s.updateActDB(lids)
	}
	log.Info("start close(%d) multiUpActDB end", yu)
}

// BatchInsertLikeExtend batch insert like_extend table.
func (s *Service) BatchInsertLikeExtend(c context.Context, extends []*like.Extend) (res int64, err error) {
	var buf bytes.Buffer
	cnt := 0
	rows := int64(0)
	for _, v := range extends {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(v.Lid, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(v.Like, 10))
		buf.WriteString("),")
		cnt++
		if cnt%500 == 0 {
			buf.Truncate(buf.Len() - 1)
			if rows, err = s.dao.AddExtend(c, buf.String()); err != nil {
				log.Error("s.dao.dealAddExtend() error(%+v)", err)
				return
			}
			res += rows
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
		if rows, err = s.dao.AddExtend(c, buf.String()); err != nil {
			log.Error("s.dao.dealAddExtend() error(%+v)", err)
			return
		}
		res += rows
	}
	return
}

// actionProc .
func (s *Service) actionProc(c context.Context, msg json.RawMessage) (err error) {
	var (
		act = new(like.Action)
	)
	if err = json.Unmarshal(msg, act); err != nil {
		log.Error("actionProc json.Unmarshal(%s) error(%v)", msg, err)
		return
	}
	s.subActionCh[act.Lid%_sharding] <- act
	return
}
