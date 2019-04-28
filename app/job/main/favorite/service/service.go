package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/app/job/main/favorite/conf"
	favDao "go-common/app/job/main/favorite/dao/fav"
	musicDao "go-common/app/job/main/favorite/dao/music"
	pubDao "go-common/app/job/main/favorite/dao/pub"
	statDao "go-common/app/job/main/favorite/dao/stat"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	arcmdl "go-common/app/service/main/archive/model/archive"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	coinmdl "go-common/app/service/main/coin/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
)

// Service favorite service.
type Service struct {
	c      *conf.Config
	waiter *sync.WaitGroup
	// fav
	cleanCDTime int64
	// dao
	pubDao  *pubDao.Dao
	statDao *statDao.Dao
	favDao  *favDao.Dao
	// databus
	consumer     *databus.Databus
	playStatSub  *databus.Databus
	favStatSub   *databus.Databus
	shareStatSub *databus.Databus
	procChan     []chan *favmdl.Message
	// rpc
	coinRPC *coinrpc.Service
	arcRPC  *arcrpc.Service2
	artRPC  *artrpc.Service
	// cache chan
	cache     *fanout.Fanout
	statMerge *statMerge
	musicDao  *musicDao.Dao
}

type statMerge struct {
	Business int
	Target   int64
	Sources  map[int64]bool
}

// New new a service and return.
func New(c *conf.Config) (s *Service) {
	if c.Fav.Proc <= 0 {
		c.Fav.Proc = 32
	}
	s = &Service{
		c:      c,
		waiter: new(sync.WaitGroup),
		// fav
		cleanCDTime: int64(time.Duration(c.Fav.CleanCDTime) / time.Second),
		// dao
		favDao:   favDao.New(c),
		pubDao:   pubDao.New(c),
		statDao:  statDao.New(c),
		musicDao: musicDao.New(c),
		// databus
		consumer: databus.New(c.JobDatabus),
		procChan: make([]chan *favmdl.Message, c.Fav.Proc),
		// stat databus
		playStatSub:  databus.New(c.MediaListCntDatabus),
		favStatSub:   databus.New(c.FavStatDatabus),
		shareStatSub: databus.New(c.ShareStatDatabus),
		// rpc
		coinRPC: coinrpc.New(c.RPCClient2.Coin),
		artRPC:  artrpc.New(c.RPCClient2.Article),
		arcRPC:  arcrpc.New2(c.RPCClient2.Archive),
		// cache chan
		cache: fanout.New("cache"),
	}

	if c.StatMerge != nil {
		s.statMerge = &statMerge{
			Business: c.StatMerge.Business,
			Target:   c.StatMerge.Target,
			Sources:  make(map[int64]bool),
		}
		for _, id := range c.StatMerge.Sources {
			s.statMerge.Sources[id] = true
		}
	}

	for i := int64(0); i < c.Fav.Proc; i++ {
		ch := make(chan *favmdl.Message, 128)
		s.procChan[i] = ch
		s.waiter.Add(1)
		go s.jobproc(ch)
	}

	s.waiter.Add(1)
	go s.consumeStat()

	s.waiter.Add(1)
	go s.consumeproc()
	return
}

func (s *Service) consumeproc() {
	offsets := make(map[int32]int64, 9)
	defer func() {
		log.Info("end databus msg offsets:%v", offsets)
		s.waiter.Done()
	}()
	for {
		msg, ok := <-s.consumer.Messages()
		if !ok {
			log.Info("consumeproc exit")
			for _, c := range s.procChan {
				close(c)
			}
			return
		}
		if _, ok := offsets[msg.Partition]; !ok {
			log.Info("begin databus msg offsets:%v", offsets)
		}
		offsets[msg.Partition] = msg.Offset

		msg.Commit()
		m := &favmdl.Message{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal() error(%v)", err)
			continue
		}
		if m.Mid <= 0 {
			log.Warn("m.Mid shuld not be equal or lesser than zero，m:%+v", m)
			continue
		}
		log.Info("consumer topic:%s, partitionId:%d, offset:%d, Key:%s, Value:%s Mid:%d Proc:%d", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value, m.Mid, s.c.Fav.Proc)
		s.procChan[m.Mid%s.c.Fav.Proc] <- m
	}
}

func (s *Service) jobproc(ch chan *favmdl.Message) {
	defer s.waiter.Done()
	for {
		m, ok := <-ch
		if !ok {
			log.Info("jobproc exit")
			return
		}
		switch m.Field {
		case favmdl.FieldResource:
			if err := s.upResource(context.Background(), m); err != nil {
				log.Error("upResource(%v) error(%v)", m, err)
				continue
			}
		default:
		}
	}
}

// Close close.
func (s *Service) Close() (err error) {
	if err = s.consumer.Close(); err != nil {
		log.Error("s.consumer.Close() error(%v)", err)
		return
	}
	return s.favDao.Close()
}

// Wait wait.
func (s *Service) Wait() {
	s.waiter.Wait()
}

// Ping ping method for server check
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.favDao.Ping(c); err != nil {
		log.Error("s.favDao.Ping error(%v)", err)
		return
	}
	return
}

// ArcRPC find archive by rpc
func (s *Service) archiveRPC(c context.Context, aid int64) (a *api.Arc, err error) {
	argAid := &arcmdl.ArgAid2{
		Aid: aid,
	}
	if a, err = s.arcRPC.Archive3(c, argAid); err != nil {
		log.Error("arcRPC.Archive3(%v, archive), err(%v)", argAid, err)
	}
	return
}

// AddCoinRpc check user whether or not banned to post
func (s *Service) addCoinRPC(c context.Context, mid int64, coin float64, reason string) (err error) {
	if _, err = s.coinRPC.ModifyCoin(c, &coinmdl.ArgModifyCoin{Mid: mid, Count: coin, Reason: reason}); err != nil {
		log.Error("coinRPC.ModifyCoin(%v, %v), err(%v)", mid, coin, err)
	}
	return
}

// articleRPC find aritile by rpc
func (s *Service) articleRPC(c context.Context, aid int64) (a map[int64]*artmdl.Meta, err error) {
	argAid := &artmdl.ArgAids{
		Aids: []int64{aid},
	}
	if a, err = s.artRPC.ArticleMetas(c, argAid); err != nil {
		log.Error("d.artRPC.ArticleMetas(%+v), error(%v)", argAid, err)
	}
	return
}

// ArcsRPC find archives by rpc.
func (s *Service) ArcsRPC(c context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	argAids := &arcmdl.ArgAids2{
		Aids: aids,
	}
	if as, err = s.arcRPC.Archives3(c, argAids); err != nil {
		log.Error("s.arcRPC.Archives3(%v), error(%v)", argAids, err)
	}
	return
}

func (s *Service) mergeTarget(business int, aid int64) int64 {
	if s.statMerge != nil && s.statMerge.Business == business && s.statMerge.Sources[aid] {
		return s.statMerge.Target
	}
	return 0
}
