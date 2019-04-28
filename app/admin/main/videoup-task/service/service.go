package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/admin/main/videoup-task/conf"
	"go-common/app/admin/main/videoup-task/dao"
	"go-common/app/admin/main/videoup-task/model"
	account "go-common/app/service/main/account/api"
	upsrpc "go-common/app/service/main/up/api/v1"
	"go-common/library/database/elastic"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// MemberCache .
type MemberCache struct {
	sync.RWMutex
	uptime time.Time
	ms     map[int64]*model.MemberStat
}

//Service service
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	es  *elastic.Elastic

	httpClient  *bm.Client
	memberCache *MemberCache
	// cache
	typeCache   map[int16]*model.Type
	reviewCache *model.ReviewCache
	twConCache  map[int8]map[int64]*model.WCItem
	upperCache  map[int8]map[int64]struct{}
	typeCache2  map[int16][]int64 // 记录每个一级分区下的二级分区

	//grpc
	accRPC account.AccountClient
	upsRPC upsrpc.UpClient
}

//New new service
func New(conf *conf.Config) (svr *Service) {
	svr = &Service{
		c:   conf,
		dao: dao.New(conf),
		es:  elastic.NewElastic(nil),

		httpClient:  bm.NewClient(conf.HTTPClient),
		memberCache: &MemberCache{},
		// cache
		reviewCache: model.NewRC(),
	}
	var err error
	if svr.accRPC, err = account.NewClient(conf.GRPC.AccRPC); err != nil {
		panic(err)
	}
	if svr.upsRPC, err = upsrpc.NewClient(conf.GRPC.UpsRPC); err != nil {
		panic(err)
	}

	go svr.memberproc()
	svr.loadConf()
	go svr.cacheproc()
	go svr.delProc()
	svr.loadRC()
	go svr.loadRCproc()

	return
}

//Close close
func (s *Service) Close() {
	s.dao.Close()
}

//Ping ping
func (s *Service) Ping(ctx context.Context) (err error) {
	err = s.dao.Ping(ctx)
	return
}

func (s *Service) loadConf() {
	var (
		err        error
		tpm        map[int16]*model.Type
		tpm2       map[int16][]int64
		twConCache map[int8]map[int64]*model.WCItem
		upm        map[int8]map[int64]struct{}
	)
	if tpm, err = s.dao.TypeMapping(context.TODO()); err != nil {
		log.Error("s.dao.TypeMapping error(%v)", err)
		return
	}
	s.typeCache = tpm

	tpm2 = make(map[int16][]int64)
	for id, tmod := range tpm {
		if tmod.PID == 0 {
			if _, ok := tpm2[id]; !ok {
				tpm2[id] = []int64{}
			}
			continue
		}
		arrid, ok := tpm2[tmod.PID]
		if !ok {
			tpm2[tmod.PID] = []int64{int64(id)}
		} else {
			tpm2[tmod.PID] = append(arrid, int64(id))
		}
	}
	s.typeCache2 = tpm2

	if twConCache, err = s.weightConf(context.TODO()); err != nil {
		log.Error("s.weightConf error(%v)", err)
		return
	}
	s.twConCache = twConCache

	upm, err = s.upSpecial(context.TODO())
	if err != nil {
		log.Error("s.upSpecial error(%v)", err)
		return
	}
	s.upperCache = upm
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(3 * time.Minute)
		s.loadConf()
	}
}
