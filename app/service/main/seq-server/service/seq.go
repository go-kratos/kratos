package service

import (
	"context"
	"fmt"

	"go-common/app/service/main/seq-server/model"
	"go-common/library/ecode"
)

const _nextStepRetry = 3

func (s *Service) nextStep(c context.Context, b *model.Business) (err error) {
	var (
		n             int
		lastSeq, rows int64
	)
	for {
		if lastSeq, err = s.db.MaxSeq(c, b.ID); err != nil {
			return
		}
		if rows, err = s.db.UpMaxSeq(c, b.ID, lastSeq+b.Step, lastSeq); err != nil {
			return
		}
		if rows > 0 {
			b.CurSeq = lastSeq
			b.MaxSeq = lastSeq + b.Step
			break
		}
		if n++; n > _nextStepRetry {
			err = fmt.Errorf("get the next step failed(id:%d maxseq:%d step:%d)", b.ID, b.MaxSeq, b.Step)
			return
		}
	}
	return
}

// ID32 get id int32.
func (s *Service) ID32(c context.Context, businessID int64, token string) (id int32, err error) {
	s.bl.RLock()
	b, ok := s.bs[businessID]
	s.bl.RUnlock()
	if !ok || b.Token != token {
		err = ecode.NothingFound
		return
	}
	b.Mutex.Lock()
	// NOTE: make sure curSeq begin with maxSeq when start from 0
	if b.CurSeq == 0 || b.CurSeq+1 > b.MaxSeq {
		if err = s.nextStep(c, b); err != nil {
			b.Mutex.Unlock()
			return
		}
	}
	b.CurSeq++
	id = int32(b.CurSeq)
	b.Mutex.Unlock()
	return
}

// UpMaxSeq update max seq by buisinessID and token.
func (s *Service) UpMaxSeq(c context.Context, businessID, maxSeq, step int64, token string) (err error) {
	rows, err := s.db.UpMaxSeqToken(c, businessID, maxSeq, step, token)
	if err != nil {
		return
	}
	if rows <= 0 {
		err = fmt.Errorf("update maxSeq failed(businessID:%d maxSeq:%d)", businessID, maxSeq)
	}
	return
}
