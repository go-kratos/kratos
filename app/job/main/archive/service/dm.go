package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	dmmdl "go-common/app/interface/main/dm2/model"
	"go-common/app/job/main/archive/model/dm"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_type             = "archive"
	_table            = "dm_subject_"
	_update           = "update"
	_subjectTypeForAv = 1
)

func (s *Service) dmConsumer() {
	defer s.waiter.Done()
	for {
		var (
			msg   *databus.Message
			ok    bool
			err   error
			canal = &dm.Canal{}
		)
		if msg, ok = <-s.dmSub.Messages(); !ok || s.closeSub {
			log.Error("s.dmSub Closed")
			return
		}
		msg.Commit()
		if err = json.Unmarshal(msg.Value, canal); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, canal)
			continue
		}
		// not dm_subject_
		if !strings.HasPrefix(canal.Table, _table) {
			log.Warn("table(%s) message(%s) skiped", canal.Table, msg.Value)
			continue
		}
		// not update
		if canal.Action != _update {
			log.Warn("action(%s) message(%s) skiped", canal.Action, msg.Value)
			continue
		}
		var subject *dm.Subject
		if err = json.Unmarshal(canal.New, &subject); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", canal.New, err)
			continue
		}
		// type must be av
		if subject.Type != _subjectTypeForAv {
			log.Warn("subject type(%s) is not av message(%+v)", subject.Type, subject)
			continue
		}
		s.dmMu.Lock()
		s.dmCids[subject.CID] = struct{}{}
		s.dmMu.Unlock()
	}
}

func (s *Service) dmCounter() {
	defer s.waiter.Done()
	for {
		time.Sleep(5 * time.Second)
		s.dmMu.Lock()
		cm := s.dmCids
		s.dmCids = make(map[int64]struct{})
		s.dmMu.Unlock()
		var (
			aids    []int64
			err     error
			c       = context.TODO()
			am      = make(map[int64][]int64)
			allCids []int64
		)
		for cid := range cm {
			if aids, err = s.archiveDao.Aids(c, cid); err != nil {
				log.Error("s.archiveDao.Aids(%d) error(%v)", err)
				continue
			}
			for _, aid := range aids {
				if _, ok := am[aid]; ok {
					continue
				}
				var pages []*api.Page
				if pages, err = s.arcServices[0].Page3(c, &archive.ArgAid2{Aid: aid}); err != nil {
					log.Error("s.arcServices[0].Page3(%d) error(%v)", aid, err)
					continue
				}
				for _, p := range pages {
					am[aid] = append(am[aid], p.Cid)
					allCids = append(allCids, p.Cid)
				}
			}
		}
		var (
			times    int
			argCount = 100
			cids     []int64
			cDmCount = make(map[int64]int64)
		)
		if len(allCids)%argCount == 0 {
			times = len(allCids) / argCount
		} else {
			times = len(allCids)/argCount + 1
		}
		for i := 0; i < times; i++ {
			if i == times-1 {
				cids = allCids[i*argCount:]
			} else {
				cids = allCids[i*argCount : (i+1)*argCount]
			}
			var sm map[int64]*dmmdl.SubjectInfo
			if sm, err = s.dm2RPC.SubjectInfos(c, &dmmdl.ArgOids{Type: 1, Oids: cids}); err != nil {
				log.Error("s.dm2RPC.SubjectInfos(%v) error(%v)", cids, err)
				continue
			}
			for cid, s := range sm {
				cDmCount[cid] = int64(s.Count)
			}
		}
	L:
		for aid, cids := range am {
			var sum int64
			for _, cid := range cids {
				var (
					cnt int64
					ok  bool
				)
				if cnt, ok = cDmCount[cid]; !ok {
					log.Error("dm cid(%d) no count", cid)
					break L
				}
				sum += cnt
			}
			dMsg := &dm.Count{ID: aid, Count: sum, Type: _type, Timestamp: time.Now().Unix()}
			if err = s.dmPub.Send(c, strconv.FormatInt(aid, 10), dMsg); err != nil {
				log.Error("s.dmPub.Send error(%v)", err)
				continue
			}
			log.Info("s.dmPub.Send(%+v) success", dMsg)
		}
		if s.closeSub {
			return
		}
	}
}
