package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/job/main/growup/model"

	"go-common/library/log"
)

var (
	_archiveTable = "archive"
	_musicTable   = "music"
	_actionUpdate = "update"
	_actionInsert = "insert"
)

func (s *Service) archiveConsume(ctx context.Context) {
	var (
		msgs = s.arcSub.Messages()
		err  error
	)
	log.Info("archiveConsume start")
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.arcSub.Messages closed", err)
			return
		}
		msg.Commit()
		archive := &model.ArchiveMsg{}
		if err = json.Unmarshal(msg.Value, archive); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if archive.Table == _musicTable {
			go s.addBgmWhiteList(archive.Action, archive.New, archive.Old)
		}
		if archive.Table == _archiveTable && archive.Action == _actionUpdate {
			go s.checkArchiveState(archive.New, archive.Old)
		}
	}
}

// (action == insert && new.state = 0) || (action = update && new.state = 0 && old.state < 0)
func (s *Service) addBgmWhiteList(action string, newMsg, oldMsg []byte) {
	nw := &model.BgmSub{}
	if err := json.Unmarshal(newMsg, nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", newMsg, err)
		return
	}
	old := &model.BgmSub{}
	if action == _actionUpdate {
		if err := json.Unmarshal(oldMsg, old); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", oldMsg, err)
			return
		}
	}
	if (action == _actionInsert && nw.State == 0) || (action == _actionUpdate && nw.State == 0 && old.State < 0) {
		log.Info("addBgmWhiteList mid(%d)", nw.MID)
		_, err := s.dao.InsertBgmWhiteList(context.Background(), nw.MID)
		if err != nil {
			log.Error(" s.dao.InsertBgmWhiteList(%d) error(%v)", nw.MID, err)
		}
	}
}

// new.state>=0 && new.Copyright =2 && old.Copyright == 1   原创变转载
// new.state>=0 && new.Copyright =1 && old.Copyright == 2   转载变原创
func (s *Service) checkArchiveState(newMsg, oldMsg []byte) {
	nw := &model.ArchiveSub{}
	if err := json.Unmarshal(newMsg, nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", newMsg, err)
		return
	}
	old := &model.ArchiveSub{}
	if err := json.Unmarshal(oldMsg, old); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", oldMsg, err)
		return
	}
	// 原创变转载
	if nw.State >= 0 && nw.Copyright == 2 && old.Copyright == 1 {
		log.Info("checkArchiveState get 1 avid(%d) mid(%d)", nw.ID, nw.MID)
		s.avBreachPre(context.Background(), nw.ID, nw.MID, 1)
	}

	// 转载变原创
	if nw.State >= 0 && nw.Copyright == 1 && old.Copyright == 2 {
		log.Info("checkArchiveState get 0 avid(%d) mid(%d)", nw.ID, nw.MID)
		s.avBreachPre(context.Background(), nw.ID, nw.MID, 0)
	}
}

func (s *Service) avBreachPre(c context.Context, aid, mid int64, state int) (err error) {
	if aid == 0 || mid == 0 {
		return
	}
	accState, err := s.dao.GetUpStateByMID(c, mid)
	if err != nil {
		log.Error(" s.dao.GetUpStateByMID(%d) error(%v)", mid, err)
		return
	}
	if accState != 3 {
		return
	}
	// status == 1 insert av_breach_pre, status == 0 update state if exist
	if state == 1 {
		val := fmt.Sprintf("%d,%d,'%s',0,1", aid, mid, time.Now().Format(_layout))
		_, err = s.dao.InsertAvBreachPre(c, val)
		if err != nil {
			log.Error("s.dao.InsertAvBreachPre error(%v)", err)
		}
	} else if state == 0 {
		_, err = s.dao.UpdateAvBreachPre(c, aid, 0, time.Now().AddDate(0, 0, -1).Format(_layout), state)
		if err != nil {
			log.Error("s.dao.UpdateAvBreachPre error(%v)", err)
		}
	}
	return
}
