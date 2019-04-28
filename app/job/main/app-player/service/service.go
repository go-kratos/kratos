package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	appmdl "go-common/app/interface/main/app-player/model/archive"
	"go-common/app/job/main/app-player/conf"
	"go-common/app/job/main/app-player/dao"
	"go-common/app/job/main/app-player/model"
	arcrpc "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_updateAct    = "update"
	_insertAct    = "insert"
	_tableArchive = "archive"
)

// Service is service.
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	arcRPC arcrpc.ArchiveClient
	// sub
	archiveNotifySub *databus.Databus
	waiter           sync.WaitGroup
	closed           bool
}

// New new a service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		dao:              dao.New(c),
		archiveNotifySub: databus.New(c.ArchiveNotifySub),
		closed:           false,
	}
	var err error
	s.arcRPC, err = arcrpc.NewClient(nil)
	if err != nil {
		panic(fmt.Sprintf("archive NewClient error(%v)", err))
	}
	s.waiter.Add(1)
	go s.arcConsumeproc()
	s.waiter.Add(1)
	go s.retryproc()
	return
}

// Close Databus consumer close.
func (s *Service) Close() {
	s.closed = true
	s.archiveNotifySub.Close()
	s.waiter.Wait()
}

// Ping is
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.PingMc(c); err != nil {
		return
	}
	return
}

// arcConsumeproc consumer archive
func (s *Service) arcConsumeproc() {
	var (
		msg *databus.Message
		ok  bool
		err error
	)
	msgs := s.archiveNotifySub.Messages()
	for {
		if msg, ok = <-msgs; !ok {
			log.Info("arc databus Consumer exit")
			break
		}
		log.Info("got databus message(%s)", msg.Value)
		msg.Commit()
		var ms = &model.Message{}
		if err = json.Unmarshal(msg.Value, ms); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		switch ms.Table {
		case _tableArchive:
			s.archiveUpdate(ms.Action, ms.New)
		}
	}
	s.waiter.Done()
}

func (s *Service) archiveUpdate(action string, nwMsg []byte) {
	nw := &model.ArcMsg
	if err := json.Unmarshal(nwMsg, nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nwMsg, err)
		return
	}
	switch action {
	case _updateAct, _insertAct:
		s.upArcCache(nw.Aid)
	}
}

func (s *Service) upArcCache(aid int64) {
	var (
		view *arcrpc.ViewReply
		arc  *appmdl.Info
		cids []int64
		err  error
	)
	defer func() {
		if err != nil {
			retry := &model.Retry{Aid: aid}
			s.dao.PushList(context.Background(), retry)
			log.Warn("upArcCache fail(%+v)", retry)
		}
	}()
	c := context.Background()
	if view, err = s.arcRPC.View(c, &arcrpc.ViewRequest{Aid: aid}); err != nil {
		if ecode.Cause(err).Equal(ecode.NothingFound) {
			err = nil
			return
		}
		log.Error("s.arcRPC.View3(%d) error(%v)", aid, err)
		return
	}
	for _, p := range view.Pages {
		cids = append(cids, p.Cid)
	}
	arc = &appmdl.Info{
		Aid:       aid,
		Cids:      cids,
		State:     view.Arc.State,
		Mid:       view.Arc.Author.Mid,
		Attribute: view.Arc.Attribute,
	}
	if err = s.dao.AddArchiveCache(context.Background(), aid, arc); err == nil {
		log.Info("update view cahce aid(%d) success", aid)
	}
}

func (s *Service) retryproc() {
	for {
		if s.closed {
			break
		}
		var (
			retry = &model.Retry{}
			bs    []byte
			err   error
		)
		if bs, err = s.dao.PopList(context.Background()); err != nil || len(bs) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		if err = json.Unmarshal(bs, retry); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", bs, err)
			continue
		}
		log.Info("retry data(%+v) start", retry)
		if retry.Aid != 0 {
			s.upArcCache(retry.Aid)
		}
		log.Info("retry data(%+v) end", retry)
	}
	s.waiter.Done()
}
