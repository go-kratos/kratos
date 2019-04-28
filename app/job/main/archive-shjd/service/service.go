package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-common/app/job/main/archive-shjd/conf"
	"go-common/app/job/main/archive-shjd/model"
	accrpc "go-common/app/service/main/account/rpc/client"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

const (
	_tableArchive = "archive"
	_tableVideo   = "archive_video"
	_actionInsert = "insert"
	_actionUpdate = "update"
	_actionDelete = "delete"
	_sharding     = 10
)

type lastTmStat struct {
	last int64
	stat *api.Stat
}

// Service service
type Service struct {
	c             *conf.Config
	waiter        sync.WaitGroup
	canal         *databus.Databus
	canalChan     chan *model.Message
	subMap        map[string]*databus.Databus
	subView       *databus.Databus
	subDm         *databus.Databus
	subReply      *databus.Databus
	subFav        *databus.Databus
	subCoin       *databus.Databus
	subShare      *databus.Databus
	subRank       *databus.Databus
	subLike       *databus.Databus
	notifyPub     *databus.Databus
	accountNotify *databus.Databus
	subStatCh     []chan *model.StatMsg
	arcRPCs       map[string]*arcrpc.Service2
	accRPC        *accrpc.Service3
	notifyMid     map[int64]struct{}
	notifyMu      sync.Mutex
	rds           *redis.Pool
	statSM        []map[int64]*lastTmStat
	close         bool
}

// New is archive service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		canal:     databus.New(c.Databus),
		canalChan: make(chan *model.Message, 10240),
		rds:       redis.NewPool(c.Redis),
		subMap:    make(map[string]*databus.Databus),
		// databus
		subView:       databus.New(c.ViewSub),
		subDm:         databus.New(c.DmSub),
		subReply:      databus.New(c.ReplySub),
		subFav:        databus.New(c.FavSub),
		subCoin:       databus.New(c.CoinSub),
		subShare:      databus.New(c.ShareSub),
		subRank:       databus.New(c.RankSub),
		subLike:       databus.New(c.LikeSub),
		notifyPub:     databus.New(c.NotifyPub),
		accountNotify: databus.New(c.AccountNotify),
		notifyMid:     make(map[int64]struct{}, 10240),
		accRPC:        accrpc.New3(nil),
	}
	s.arcRPCs = make(map[string]*arcrpc.Service2)
	for _, cc := range c.ArchiveRPCs {
		s.arcRPCs[cc.Cluster] = arcrpc.New2(cc)
	}
	s.subMap[model.TypeForView] = s.subView
	s.subMap[model.TypeForDm] = s.subDm
	s.subMap[model.TypeForReply] = s.subReply
	s.subMap[model.TypeForFav] = s.subFav
	s.subMap[model.TypeForCoin] = s.subCoin
	s.subMap[model.TypeForShare] = s.subShare
	s.subMap[model.TypeForRank] = s.subRank
	s.subMap[model.TypeForLike] = s.subLike
	for i := 0; i < _sharding; i++ {
		s.waiter.Add(1)
		go s.canalChanproc()
		s.subStatCh = append(s.subStatCh, make(chan *model.StatMsg, 10240))
		s.statSM = append(s.statSM, map[int64]*lastTmStat{})
		s.waiter.Add(1)
		go s.statDealproc(i)
	}
	for k, d := range s.subMap {
		s.waiter.Add(1)
		go s.consumerproc(k, d)
	}
	s.waiter.Add(1)
	go s.canalproc()
	s.waiter.Add(1)
	go s.retryconsumer()
	s.waiter.Add(1)
	go s.accountNotifyproc()
	s.waiter.Add(1)
	go s.clearMidCache()
	return s
}

func (s *Service) canalChanproc() {
	defer s.waiter.Done()
	for {
		m, ok := <-s.canalChan
		if !ok {
			log.Info("canalChanproc closed")
			return
		}
		log.Info("got canal message table(%s) action(%s) old(%s) new(%s)", m.Table, m.Action, m.Old, m.New)
		var err error
		switch m.Table {
		case _tableArchive:
			var (
				old *model.Archive
				nw  *model.Archive
			)
			switch m.Action {
			case _actionInsert:
				if err = json.Unmarshal(m.New, &nw); err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
					continue
				}
			case _actionUpdate:
				if err = json.Unmarshal(m.Old, &old); err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", m.Old, err)
					continue
				}
				if err = json.Unmarshal(m.New, &nw); err != nil {
					log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
					continue
				}
			default:
				log.Warn("got unknow action(%s)", m.Action)
				continue
			}
			s.UpdateCache(old, nw, m.Action)
		case _tableVideo:
			var video *model.Video
			if err = json.Unmarshal(m.New, &video); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", m.New, err)
				continue
			}
			switch m.Action {
			case _actionInsert, _actionUpdate:
				err = s.UpdateVideoCache(video.AID, video.CID)
			case _actionDelete:
				err = s.DelteVideoCache(video.AID, video.CID)
			default:
				bs, _ := json.Marshal(m)
				log.Error("unknow action(%s) message(%s)", m.Action, bs)
			}
			if err != nil {
				log.Error("%+v", err)
				continue
			}
		default:
			log.Warn("table(%s) skiped", m.Table)
		}
	}
}

func (s *Service) canalproc() {
	defer s.waiter.Done()
	msgs := s.canal.Messages()
	for {
		msg, ok := <-msgs
		if !ok || s.close {
			close(s.canalChan)
			log.Info("s.closed databus canal")
			return
		}
		var (
			m   = &model.Message{}
			err error
		)
		msg.Commit()
		log.Info("got message(%s)", msg.Value)
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		s.canalChan <- m
	}
}

// Ping check status
func (s *Service) Ping() (err error) {
	conn := s.rds.Get(context.TODO())
	defer conn.Close()
	if _, err = conn.Do("SET", "PING", "PONG"); err != nil {
		err = errors.Wrap(err, "redis ping")
		return
	}
	return
}

// Close is
func (s *Service) Close() (err error) {
	s.close = true
	time.Sleep(5 * time.Second)
	s.canal.Close()
	s.waiter.Wait()
	return
}
