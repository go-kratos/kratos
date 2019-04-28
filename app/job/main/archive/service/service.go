package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	dm2rpc "go-common/app/interface/main/dm2/rpc/client"
	"go-common/app/job/main/archive/conf"
	"go-common/app/job/main/archive/dao/archive"
	"go-common/app/job/main/archive/dao/email"
	"go-common/app/job/main/archive/dao/monitor"
	"go-common/app/job/main/archive/dao/reply"
	"go-common/app/job/main/archive/dao/result"
	dbusmdl "go-common/app/job/main/archive/model/databus"
	resmdl "go-common/app/job/main/archive/model/result"
	"go-common/app/job/main/archive/model/retry"
	accgrpc "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	xredis "go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service service
type Service struct {
	c                *conf.Config
	closeRetry       bool
	closeSub         bool
	archiveDao       *archive.Dao
	emailDao         *email.Dao
	monitorDao       *monitor.Dao
	replyDao         *reply.Dao
	resultDao        *result.Dao
	redis            *xredis.Pool
	waiter           sync.WaitGroup
	videoupSub       *databus.Databus
	archiveResultPub *databus.Databus
	dmPub            *databus.Databus
	dmSub            *databus.Databus
	cacheSub         *databus.Databus
	accountNotifySub *databus.Databus
	sfTpsCache       map[int16]int16
	adtTpsCache      map[int16]struct{}
	arcServices      []*arcrpc.Service2
	accGRPC          accgrpc.AccountClient
	dm2RPC           *dm2rpc.Service
	// databus channel
	videoupAids []chan int64
	pgcAids     chan int64
	// dm count
	dmCids    map[int64]struct{}
	dmMu      sync.Mutex
	notifyMid map[int64]struct{}
	notifyMu  sync.Mutex
}

// New is archive service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		archiveDao:       archive.New(c),
		emailDao:         email.New(c),
		monitorDao:       monitor.New(c),
		replyDao:         reply.New(c),
		resultDao:        result.New(c),
		dm2RPC:           dm2rpc.New(c.Dm2RPC),
		videoupSub:       databus.New(c.VideoupSub),
		dmSub:            databus.New(c.DmSub),
		dmPub:            databus.New(c.DmPub),
		archiveResultPub: databus.New(c.ArchiveResultPub),
		cacheSub:         databus.New(c.CacheSub),
		accountNotifySub: databus.New(c.AccountNotifySub),
		redis:            xredis.NewPool(c.Redis),
		pgcAids:          make(chan int64, 1024),
		dmCids:           make(map[int64]struct{}),
		notifyMid:        make(map[int64]struct{}),
		arcServices:      make([]*arcrpc.Service2, 0),
	}
	var err error
	if s.accGRPC, err = accgrpc.NewClient(nil); err != nil {
		panic(fmt.Sprintf("account.service grpc not found!!!!!!!!!!!! error(%v)", err))
	}
	for _, sc := range s.c.ArchiveServices {
		s.arcServices = append(s.arcServices, arcrpc.New2(sc))
	}
	for i := 0; i < s.c.ChanSize; i++ {
		s.videoupAids = append(s.videoupAids, make(chan int64, 1024))
		s.waiter.Add(1)
		go s.consumerVideoup(i)
		s.waiter.Add(1)
		go s.pgcConsumer()
	}
	s.loadType()
	go s.cacheproc()
	// sync archive_result db!!!!!!!
	s.waiter.Add(1)
	go s.videoupConsumer()
	s.waiter.Add(1)
	go s.dmConsumer()
	// check consumer
	go s.checkConsume()
	s.waiter.Add(1)
	go s.retryproc()
	s.waiter.Add(1)
	go s.dmCounter()
	s.waiter.Add(1)
	go s.cachesubproc()
	s.waiter.Add(1)
	go s.accountNotifyproc()
	s.waiter.Add(1)
	go s.clearMidCache()
	return s
}

func (s *Service) sendNotify(upInfo *resmdl.ArchiveUpInfo) {
	var (
		nw  []byte
		old []byte
		err error
		msg *dbusmdl.Message
		c   = context.TODO()
		rt  = &retry.Info{}
	)
	if nw, err = json.Marshal(upInfo.Nw); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", upInfo.Nw, err)
		return
	}
	if old, err = json.Marshal(upInfo.Old); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", upInfo.Old, err)
		return
	}
	msg = &dbusmdl.Message{Action: upInfo.Action, Table: upInfo.Table, New: nw, Old: old}
	if err = s.archiveResultPub.Send(c, strconv.FormatInt(upInfo.Nw.AID, 10), msg); err != nil {
		log.Error("s.archiveResultPub.Send(%+v) error(%v)", msg, err)
		rt.Action = retry.FailDatabus
		rt.Data.Aid = upInfo.Nw.AID
		rt.Data.DatabusMsg = upInfo
		s.PushFail(c, rt)
		return
	}
	msgStr, _ := json.Marshal(msg)
	log.Info("sendNotify(%s) successed", msgStr)
}

func (s *Service) loadType() {
	tpm, err := s.archiveDao.TypeMapping(context.TODO())
	if err != nil {
		log.Error("s.dede.TypeMapping error(%v)", err)
		return
	}
	s.sfTpsCache = tpm
	// audit types
	adt, err := s.archiveDao.AuditTypesConf(context.TODO())
	if err != nil {
		log.Error("s.dede.AuditTypesConf error(%v)", err)
		return
	}
	s.adtTpsCache = adt
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(1 * time.Minute)
		s.loadType()
	}
}

// check consumer stat
func (s *Service) checkConsume() {
	if s.c.Env != "pro" {
		return
	}
	for {
		time.Sleep(1 * time.Minute)
		for i := 0; i < s.c.ChanSize; i++ {
			if l := len(s.videoupAids[i]); l > s.c.MonitorSize {
				s.monitorDao.Send(context.TODO(), s.c.WeChantUsers, fmt.Sprintf("archive-job报警了啊\n UGC的chan太大了！！！\n s.videoupAids[%d] size(%d) is too large\n 是不是有人在刷数据！！！！", i, l), s.c.WeChatToken, s.c.WeChatSecret)
			}
		}
		if l := len(s.pgcAids); l > s.c.MonitorSize {
			s.monitorDao.Send(context.TODO(), s.c.WeChantUsers, fmt.Sprintf("archive-job报警了啊\n PGC的chan太大了！！！\n chan size(%d) is too large \n 是不是有人在刷数据！！！！", l), s.c.WeChatToken, s.c.WeChatSecret)
		}
	}
}

// Close kafaka consumer close.
func (s *Service) Close() (err error) {
	s.closeSub = true
	time.Sleep(2 * time.Second)
	s.videoupSub.Close()
	s.closeRetry = true
	s.waiter.Wait()
	return
}
