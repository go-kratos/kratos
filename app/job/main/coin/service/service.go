package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-common/app/job/main/coin/conf"
	"go-common/app/job/main/coin/dao"
	"go-common/app/job/main/coin/model"
	accrpc "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	coinmdl "go-common/app/service/main/coin/model"
	memrpc "go-common/app/service/main/member/api"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
)

// Service coin job service.
type Service struct {
	coinDao  *dao.Dao
	c        *conf.Config
	waiter   *sync.WaitGroup
	accRPC   accrpc.AccountClient
	memRPC   memrpc.MemberClient
	arcRPC   *arcrpc.Service2
	coinRPC  *coinrpc.Service
	databus  *databus.Databus
	group    *databusutil.Group
	expGroup *databusutil.Group
}

// New new and return service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		coinDao: dao.New(c),
		c:       c,
		waiter:  new(sync.WaitGroup),
		arcRPC:  arcrpc.New2(c.ArchiveRPC),
		coinRPC: coinrpc.New(c.CoinRPC),
	}
	var err error
	if s.memRPC, err = memrpc.NewClient(c.MemRPC); err != nil {
		panic(err)
	}
	if s.accRPC, err = accrpc.NewClient(c.AccountRPC); err != nil {
		panic(err)
	}
	s.databus = databus.New(c.Databus)
	g := databusutil.NewGroup(c.Databusutil, databus.New(c.LoginDatabus).Messages())
	g.New = newMsg
	g.Split = split
	g.Do = s.awardDo
	g.Start()
	s.group = g
	eg := databusutil.NewGroup(c.Databusutil, databus.New(c.ExpDatabus).Messages())
	eg.New = newExpMsg
	eg.Split = split
	eg.Do = s.awardDo
	eg.Start()
	s.expGroup = eg

	s.waiter.Add(1)
	go s.consumeproc()
	go s.settleproc()
	return
}

func (s *Service) consumeproc() {
	defer s.waiter.Done()
	var (
		msg    *databus.Message
		err    error
		ok     bool
		period *model.CoinSettlePeriod
		ctx    = context.TODO()
	)
	for {
		if msg, ok = <-s.databus.Messages(); !ok {
			log.Error("s.databus.Message err(%v)", err)
			return
		}
		r := &coinmdl.Record{}
		if err = json.Unmarshal([]byte(msg.Value), r); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			dao.PromError("msg:JSON")
			continue
		}
		var ip = r.IPV6
		for i := 0; i < 3; i++ {
			if err = s.addCoinExp(ctx, r.Mid, r.AvType, r.Multiply, ip); err != nil {
				time.Sleep(time.Millisecond * 50)
				continue
			}
			break
		}
		if err != nil {
			log.Errorv(ctx, log.KV("log", "fix: addCoinExp"), log.KV("mid", r.Mid), log.KV("type", r.AvType), log.KV("num", r.Multiply), log.KV("ip", ip))
			continue
		}
		if err = s.updateAddCoin(ctx, r); err != nil {
			continue
		}
		at := time.Unix(r.Timestamp, 0)
		if period, err = s.coinDao.HitSettlePeriod(ctx, at); err != nil {
			log.Errorv(ctx, log.KV("log", "s.coinDao.HitCoinPeriod"), log.KV("record", r), log.KV("err", err))
			dao.PromError("service:HitSettlePeriod")
			continue
		}
		for i := 0; ; i++ {
			if err = s.coinDao.UpsertSettle(ctx, period.ID, r.Up, r.Aid, r.AvType, r.Multiply, time.Now()); err != nil {
				log.Error("s.coinDao.UpsertCoinSettle(%d, %d, %d) error(%v)", r.Up, r.Aid, r.Multiply)
				dao.PromError("service:UpsertSettle")
				i++
				if i > 5 {
					// if env.DeployEnv == env.DeployEnvProd {
					// 	s.moni.Sms(ctx, s.c.Sms.Phone, s.c.Sms.Token, "coin-job upsetSettle fail for 5 time")
					// }
					break
				}
				continue
			}
			break
		}
		log.Info("key: %s,partion:%d,offset:%d success %s", msg.Key, msg.Partition, msg.Offset, msg.Value)
		err = msg.Commit()
		if err != nil {
			log.Error("msg.Commit partition:%d offset:%d err %v", msg.Partition, msg.Offset, err)
			dao.PromError("service:msgCommit")
		}
	}
}

// Close close service.
func (s *Service) Close() {
	s.group.Close()
	s.expGroup.Close()
	s.databus.Close()
}

// Wait wait routine unitl all close.
func (s *Service) Wait() {
	s.waiter.Wait()
}

// Ping check service health.
func (s *Service) Ping(c context.Context) error {
	return s.coinDao.Ping(c)
}

func (s *Service) updateAddCoin(c context.Context, record *coinmdl.Record) (err error) {
	if record == nil {
		return
	}
	if err = s.coinRPC.UpdateAddCoin(c, record); err != nil {
		log.Errorv(c, log.KV("log", "UpdateAddCoin"), log.KV("mid", record.Mid), log.KV("record", record))
		dao.PromError("service:updateAddCoin")
	}
	return
}
