package service

import (
	"context"
	"encoding/json"

	"go-common/app/service/main/rank/model"
	"go-common/library/log"
)

func (s *Service) consumeStatView() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.statViewSub.Messages()
		if !ok {
			log.Info("consumeproc exit")
			close(s.procChan)
			return
		}
		msg.Commit()
		m := &model.StatViewMsg{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal() error(%v)", err)
			continue
		}
		// log.Info("consumer topic:%s, Key:%s, Value:%s ", msg.Topic, msg.Key, msg.Value)
		s.procChan <- m
	}
}

func (s *Service) consumeArchive() {
	defer s.waiter.Done()
	var err error
	for {
		msg, ok := <-s.archiveSub.Messages()
		if !ok {
			log.Error("s.archiveSub.Messages channel closed")
			return
		}
		m := &model.CanalMsg{}
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg, err)
			continue
		}
		log.Info("consumeArchive topic:%s, Key:%s, Value:%s, ", msg.Topic, msg.Key, msg.Value)
		switch m.Table {
		case "archive":
			s.setArchiveMeta(context.Background(), m)
		default:
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit msg(%v) error(%v)", msg, err)
		}
	}
}

func (s *Service) consumeArchiveTv() {
	defer s.waiter.Done()
	var err error
	for {
		msg, ok := <-s.archiveTvSub.Messages()
		if !ok {
			log.Error("s.archiveTvSub.Messages channel closed")
			return
		}
		m := &model.CanalMsg{}
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg, err)
			continue
		}
		// log.Info("consumeArchiveTv topic:%s, Key:%s, Value:%s ", msg.Topic, msg.Key, msg.Value)
		switch m.Table {
		case "ugc_archive":
			s.setTv(context.Background(), m)
		default:
		}
		if err = msg.Commit(); err != nil {
			log.Error("commit msg(%v) error(%v)", msg, err)
		}
	}
}

func (s *Service) batchProc(ch chan *model.StatViewMsg) {
	defer s.waiter.Done()
	for {
		m, ok := <-ch
		if !ok {
			log.Info("jobproc exit")
			return
		}
		switch m.Type {
		case "archive":
			s.field(m.ID).Click = m.Count
		default:
		}
	}
}

func (s *Service) setArchiveMeta(c context.Context, m *model.CanalMsg) {
	o := &model.ArchiveMeta{}
	if err := json.Unmarshal(m.New, o); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", m, err)
		return
	}
	switch m.Action {
	case model.SyncInsert:
		s.field(o.Aid).Flag = model.FlagExist
		s.field(o.Aid).Oid = o.Aid
		s.field(o.Aid).Pubtime = o.SetPubtime()
		typeMap, err := s.dao.ArchiveTypes(c, []int64{o.Typeid})
		if err != nil {
			log.Error("s.dao.ArchiveTypes(%d)", o.Typeid, err)
			return
		}
		if v, ok := typeMap[o.Typeid]; ok {
			s.field(o.Aid).Pid = v.SetPid()
		}
	case model.SyncUpdate:
		oo := &model.ArchiveMeta{}
		if err := json.Unmarshal(m.Old, oo); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", m, err)
			return
		}
		s.field(o.Aid).Flag = model.FlagExist
		if o.Typeid != oo.Typeid {
			typeMap, err := s.dao.ArchiveTypes(c, []int64{o.Typeid})
			if err != nil {
				log.Error("s.dao.ArchiveTypes(%d)", o.Typeid, err)
				return
			}
			if v, ok := typeMap[o.Typeid]; ok {
				s.field(o.Aid).Pid = v.SetPid()
			}
		}
		if o.Pubtime != oo.Pubtime {
			s.field(o.Aid).Pubtime = o.SetPubtime()
		}
	}
}

func (s *Service) setTv(c context.Context, m *model.CanalMsg) {
	o := &model.ArchiveMeta{}
	if err := json.Unmarshal(m.New, o); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", m, err)
		return
	}
	switch m.Action {
	case model.SyncInsert:
		s.field(o.Aid).Result = o.Result
		s.field(o.Aid).Deleted = o.Deleted
		s.field(o.Aid).Valid = o.Valid
	case model.SyncUpdate:
		oo := &model.ArchiveMeta{}
		if err := json.Unmarshal(m.Old, oo); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", m, err)
			return
		}
		if o.Result != oo.Result {
			s.field(o.Aid).Result = o.Result
		}
		if o.Deleted != oo.Deleted {
			s.field(o.Aid).Deleted = o.Deleted
		}
		if o.Valid != oo.Valid {
			s.field(o.Aid).Valid = o.Valid
		}
	}
}
