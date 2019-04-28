package service

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"
	"golang.org/x/sync/errgroup"

	"go-common/app/job/main/growup/model"
)

// account_type, 0: All; 1: UGC; 2: PGC
// o: old state; n: new state
func (s *Service) expiredCheck(accType int, o int, n int, table string) (err error) {
	c := context.TODO()
	var m map[int64]xtime.Time
	if accType == 0 {
		m, err = s.dao.MIDsByState(c, o, table)
		if err != nil {
			log.Error("s.dao.MIDsByState error(%v)", err)
			return
		}
	}

	if accType == 1 || accType == 2 {
		m, err = s.dao.MIDsByStateType(c, accType, o, table)
		if err != nil {
			log.Error("s.dao.MIDsByStateType error(%v)", err)
			return
		}
	}

	var freed []int64
	now := time.Now().Unix()
	for mid, exp := range m {
		if now > int64(exp) {
			freed = append(freed, mid)
		}
	}
	if len(freed) == 0 {
		log.Info("no mid cooldown")
		return
	}
	_, err = s.dao.UpdateAccountState(c, n, freed, table)
	return
}

// UpdateUpInfo update up_info_video
func (s *Service) UpdateUpInfo(c context.Context) (err error) {
	ch := make(chan []int64, 100)
	var group errgroup.Group
	group.Go(func() (err error) {
		defer close(ch)
		return s.mids(c, ch)
	})
	group.Go(func() (err error) {
		return s.update(c, ch)
	})
	return group.Wait()
}

func (s *Service) mids(c context.Context, ch chan []int64) (err error) {
	var id int64
	var ms []int64
	for {
		id, ms, err = s.dao.MIDs(c, id, 2000)
		if err != nil {
			return
		}
		if len(ms) == 0 {
			break
		}
		ch <- ms
	}
	return
}

func (s *Service) update(c context.Context, ch chan []int64) (err error) {
	for mids := range ch {
		var bs []*model.UpBaseInfo
		bs, err = s.dao.GetUpBaseInfo(c, mids)
		if err != nil {
			return
		}
		// batch update bs
		values := baseInfoValues(bs)
		_, err = s.dao.UpdateUpInfo(c, values)
		if err != nil {
			return
		}
	}
	return
}

func baseInfoValues(bs []*model.UpBaseInfo) (values string) {
	var buf bytes.Buffer
	for _, b := range bs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(b.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(b.Fans))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(b.TotalPlayCount))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(b.OriginalArchiveCount))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(b.Avs))
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
