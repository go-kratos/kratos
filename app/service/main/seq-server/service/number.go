package service

import (
	"context"
	"time"

	"go-common/app/service/main/seq-server/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_maxSeq = 1024
)

// tilNextMillis spin wait till next millisecond.
func tilNextMillis(lastTimestamp int64) (now int64) {
	now = time.Now().UnixNano() / int64(time.Millisecond)
	for now <= lastTimestamp {
		now = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return now
}

// ID get id
func (s *Service) ID(c context.Context, businessID int64, token string) (id int64, err error) {
	s.bl.RLock()
	b, ok := s.bs[businessID]
	s.bl.RUnlock()
	if !ok || b.Token != token {
		err = ecode.NothingFound
		log.Error("businessID(%d) not found!", businessID)
		return
	}
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	now := time.Now().UnixNano() / int64(time.Millisecond)
	if b.LastTimestamp > now {
		err = ecode.ServerErr
		log.Error("clock is moving backwards. Rejecting requests until %d.", b.LastTimestamp)
		return
	}
	if b.LastTimestamp == now {
		b.CurSeq++
		if b.CurSeq >= _maxSeq {
			b.CurSeq = 0
			now = tilNextMillis(b.LastTimestamp)
		}
	} else {
		b.CurSeq = 0
	}
	b.LastTimestamp = now
	id = b.Perch<<62 | (now-b.BenchTime)<<19 | b.CurSeq<<6 | s.svrNum<<s.idcLen | s.idc
	return
}

// CheckVersion check server version
func (s *Service) CheckVersion(c context.Context) (ver *model.SeqVersion, err error) {
	ver = &model.SeqVersion{IDC: s.idc, SvrNum: s.svrNum, SvrTime: time.Now().Unix()}
	return
}
