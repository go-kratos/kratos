package service

import (
	"context"
	"time"

	coinrpc "go-common/app/service/main/coin/api/gorpc"
	memrpc "go-common/app/service/main/member/api/gorpc"
	pointrpc "go-common/app/service/main/point/rpc/client"
	"go-common/app/service/main/usersuit/conf"
	inviteDao "go-common/app/service/main/usersuit/dao/invite"
	medalDao "go-common/app/service/main/usersuit/dao/medal"
	pendantDao "go-common/app/service/main/usersuit/dao/pendant"
	"go-common/app/service/main/usersuit/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct of service.
type Service struct {
	d          *inviteDao.Dao
	pendantDao *pendantDao.Dao
	medalDao   *medalDao.Dao
	// conf
	c                 *conf.Config
	missch            chan func()
	notifych          chan func()
	memRPC            *memrpc.Service
	coinRPC           *coinrpc.Service
	pointRPC          *pointrpc.Service
	merchantID        string
	merchantProductID string
	callBackURL       string
	groupInfo         []*model.PendantGroupInfo
	pendantInfo       []*model.Pendant
	groupMap          map[int64]*model.PendantGroupInfo
	pendantMap        map[int64]*model.Pendant
	medalInfoAll      map[int64]*model.MedalInfo
	medalGroupAll     []*model.MedalGroup
	accountNotifyPub  *databus.Databus
	gidMap            map[int64][]int64
	pidMap            map[int64]int64
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                 c,
		d:                 inviteDao.New(c),
		pendantDao:        pendantDao.New(c),
		medalDao:          medalDao.New(c),
		missch:            make(chan func(), 10240),
		notifych:          make(chan func(), 1024000),
		memRPC:            memrpc.New(c.GORPCClient.Member),
		coinRPC:           coinrpc.New(c.GORPCClient.Coin),
		pointRPC:          pointrpc.New(c.GORPCClient.Point),
		merchantID:        c.PayInfo.MerchantID,
		merchantProductID: c.PayInfo.MerchantProductID,
		callBackURL:       c.PayInfo.CallBackURL,
		pendantMap:        make(map[int64]*model.Pendant),
		accountNotifyPub:  databus.New(conf.Conf.AccountNotify),
		gidMap:            make(map[int64][]int64),
		pidMap:            make(map[int64]int64),
	}
	s.loadGidRefPid()
	s.loadMedal()
	s.refreshGroup(context.TODO())
	go s.loadgidrefpidproc()
	go s.cacheproc()
	go s.refreshGroup(context.TODO())
	go s.taskproc()
	go s.notifyproc()
	go s.loadmedalproc()
	return
}

func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}

func (s *Service) addNotify(f func()) {
	select {
	case s.notifych <- f:
	default:
		log.Warn("addNotify chan full")
	}
}

// notifyproc nofity clear cache
func (s *Service) notifyproc() {
	for {
		f := <-s.notifych
		f()
	}
}

// Close dao.
func (s *Service) Close() {
	s.d.Close()
	s.pendantDao.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.d.Ping(c); err != nil {
		return
	}
	return
}

// Task regular get appkey info
func (s *Service) taskproc() {
	for {
		time.Sleep(2 * time.Minute)
		s.refreshGroup(context.TODO())
	}
}

// rerefreshGroup when group info changed
func (s *Service) refreshGroup(c context.Context) {
	var err error
	s.groupInfo, err = s.AllGroupInfo(c)
	if err != nil {
		log.Error("s.AllGroupInfo error(%+v)", err)
		return
	}
	s.pendantInfo, err = s.PendantAll(c)
	if err != nil {
		log.Error("s.PendantAll error(%+v)", err)
		return
	}
	groupMap := make(map[int64]*model.PendantGroupInfo)
	pendantMap := make(map[int64]*model.Pendant)
	for _, v := range s.groupInfo {
		if _, ok := groupMap[v.ID]; !ok {
			groupMap[v.ID] = v
		}
	}
	for _, v := range s.pendantInfo {
		if _, ok := pendantMap[v.ID]; !ok {
			pendantMap[v.ID] = v
		}
	}
	s.pendantMap = pendantMap
	s.groupMap = groupMap
}

func (s *Service) loadGidRefPid() {
	var (
		err    error
		c      = context.TODO()
		pidMap map[int64]int64
		gidMap map[int64][]int64
	)
	if gidMap, pidMap, err = s.pendantDao.GIDRefPID(c); err != nil {
		log.Error("s.pendantDao.GIDRefPID error(%+v)", err)
		return
	}
	s.gidMap = gidMap
	s.pidMap = pidMap
}

func (s *Service) loadmedalproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadMedal panic(%v)", x)
			go s.loadMedal()
			log.Info("service.loadMedal recover")
		}
	}()
	for {
		time.Sleep(2 * time.Minute)
		s.loadMedal()
	}
}

func (s *Service) loadgidrefpidproc() {
	for {
		time.Sleep(2 * time.Minute)
		s.loadGidRefPid()
	}
}
