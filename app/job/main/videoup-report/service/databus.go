package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/manager"
	"go-common/library/log"
)

var (
	_archive   = "archive"
	_video     = "archive_video"
	_insertAct = "insert"
	_updateAct = "update"
	//_delete    = "delete"
)

// consumer binlog
func (s *Service) arcCanalConsume() {
	defer s.waiter.Done()
	var (
		msgs = s.archiveSub.Messages()
		err  error
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.archiveSub.Message closed", err)
			return
		}
		msg.Commit()
		m := &archive.Message{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		log.Info("arcCanalConsume msg(%s)", msg.Value)
		log.Info("arcCanalConsume topic(%s) partition(%d) offset(%d)  commit start", msg.Topic, msg.Partition, msg.Offset)
		if msg.Offset >= s.c.BeginOffset {
			log.Info("arcCanalConsume offset(%d) is hit BeginOffset(%d) and start track data", msg.Offset, s.c.BeginOffset)
			if m.Table == _archiveTable {
				s.putArcChan(m.Action, m.New, m.Old)
			}
			if m.Table == _videoTable {
				s.putVideoChan(m.Action, m.New, m.Old)
			}
		} else {
			log.Info("arcCanalConsume offset(%d) not hit BeginOffset(%d) and pass", msg.Offset, s.c.BeginOffset)
		}
		//todo 异步消费
		if m.Table == _video && m.Action == _updateAct {
			s.hdlVideoUpdateBinLog(m.New, m.Old)
		}
		if m.Table == _archive {
			s.hdlArchiveMessage(m.Action, m.New, m.Old)
		}
	}
}

func (s *Service) videoupConsumer() {
	defer s.waiter.Done()
	var (
		msgs = s.videoupSub.Messages()
		err  error
		c    = context.TODO()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.videoupSub.Messages closed")
			return
		}
		msg.Commit()
		m := &archive.VideoupMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			return
		}
		log.Info("videoupMessage key(%s) value(%s) partition(%d) offset(%d) route(%s) commit start", msg.Key, msg.Value, msg.Partition, msg.Offset, m.Route)
		switch m.Route {
		case archive.RoutePostFirstRound:
			err = s.postFirstRound(c, m)
		case archive.RouteSecondRound:
			err = s.secondRound(c, m)
		case archive.RouteAddArchive:
			err = s.addArchive(c, m)
		case archive.RouteModifyArchive:
			err = s.modifyArchive(c, m)
		case archive.RouteAutoOpen:
			err = s.autoOpen(c, m)
		case archive.RouteDelayOpen:
			err = s.delayOpen(c, m)
		default:
			log.Warn("videoupConsumer unknown message route(%s)", m.Route)
		}
		if err == nil {
			log.Info("videoupMessage key(%s) value(%s) partition(%d) offset(%d) end", msg.Key, msg.Value, msg.Partition, msg.Offset)
		} else {
			log.Error("videoupMessage key(%s) value(%s) partition(%d) offset(%d) error(%v)", msg.Key, msg.Value, msg.Partition, msg.Offset, err)
		}
	}
}

// managerDBConsume 消费manager binlog
func (s *Service) managerDBConsume() {
	defer s.waiter.Done()
	var (
		err  error
		msgs = s.ManagerDBSub.Messages()
	)
	for {
		msg, open := <-msgs
		if !open {
			log.Info("managerDBConsume s.arcResultSub.Messages is closed")
			return
		}
		if msg == nil {
			continue
		}
		msg.Commit()
		log.Info("managerDBConsume consume key(%s) offset(%d) value(%s)", msg.Key, msg.Offset, string(msg.Value))

		m := &manager.BinMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("managerDBConsume json.Unmarshal error(%v)", err)
			continue
		}
		switch m.Table {
		case _upsTable:
			s.hdlManagerUpsBinlog(m)
		}
	}
}
