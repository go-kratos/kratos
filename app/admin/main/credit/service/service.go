package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/admin/main/credit/conf"
	accdao "go-common/app/admin/main/credit/dao/account"
	blockedao "go-common/app/admin/main/credit/dao/blocked"
	managerdao "go-common/app/admin/main/credit/dao/manager"
	msgdao "go-common/app/admin/main/credit/dao/msg"
	relationdao "go-common/app/admin/main/credit/dao/relation"
	searchdao "go-common/app/admin/main/credit/dao/search"
	uploaddao "go-common/app/admin/main/credit/dao/upload"
	blkmodel "go-common/app/admin/main/credit/model/blocked"
	coinclient "go-common/app/service/main/coin/api"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// Service struct of service.
type Service struct {
	// conf
	c             *conf.Config
	accDao        *accdao.Dao
	blockedDao    *blockedao.Dao
	managerDao    *managerdao.Dao
	searchDao     *searchdao.Dao
	uploadDao     *uploaddao.Dao
	msgDao        *msgdao.Dao
	RelationDao   *relationdao.Dao
	DB            *gorm.DB
	ReadDB        *gorm.DB
	Search        *searchdao.Dao
	caseConfCache map[string]string
	Managers      map[int64]string
	caseReasons   map[int]string
	MsgCh         chan *blkmodel.SysMsg
	notifych      chan func()
	stop          chan struct{}
	// wait
	wg         sync.WaitGroup
	coinClient coinclient.CoinClient
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		blockedDao:  blockedao.New(c),
		managerDao:  managerdao.New(c),
		accDao:      accdao.New(c),
		searchDao:   searchdao.New(c),
		uploadDao:   uploaddao.New(c),
		RelationDao: relationdao.New(c),
		msgDao:      msgdao.New(c),
		MsgCh:       make(chan *blkmodel.SysMsg, c.ChanSize.SysMsg),
		notifych:    make(chan func(), c.ChanSize.SysMsg),
		stop:        make(chan struct{}),
	}
	s.DB = s.blockedDao.DB
	s.ReadDB = s.blockedDao.ReadDB
	s.Search = s.searchDao
	var err error
	if s.coinClient, err = coinclient.NewClient(c.CoinClient); err != nil {
		panic(err)
	}
	s.loadConfig()
	s.loadManager()
	go s.cacheproc()
	go s.loadManagerproc()
	go s.msgproc()
	s.wg.Add(1)
	go s.notifyproc()
	return s
}

func (s *Service) cacheproc() {
	for {
		time.Sleep(3 * time.Minute)
		s.loadConfig()
	}
}

func (s *Service) loadManagerproc() {
	for {
		time.Sleep(1 * time.Hour)
		s.loadManager()
	}
}

func (s *Service) loadConfig() {
	cc, err := s.CaseConf(context.TODO())
	if err != nil {
		log.Error("s.CaseConf error(%v)", err)
		return
	}
	s.caseConfCache = cc
	cr, err := s.CaseReason(context.TODO())
	if err != nil {
		log.Error("s.CaseReason error(%v)", err)
		return
	}
	s.caseReasons = cr
}

func (s *Service) loadManager() {
	managers, err := s.managerDao.Managers(context.TODO())
	if err != nil {
		log.Error("s.Managers error(%v)", err)
		return
	}
	s.Managers = managers
}

// AddNotify .
func (s *Service) AddNotify(f func()) {
	select {
	case s.notifych <- f:
	default:
		log.Warn("addNotify chan full")
	}
}

// notifyproc nofity clear cache
func (s *Service) notifyproc() {
	defer s.wg.Done()
	for {
		f, ok := <-s.notifych
		if !ok {
			log.Warn("s.notifych chan is close")
			return
		}
		f()
	}
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) error {
	return s.blockedDao.Ping(c)
}

// Close dao.
func (s *Service) Close() {
	s.blockedDao.Close()
	close(s.MsgCh)
	close(s.notifych)
	close(s.stop)
	time.Sleep(1 * time.Second)
	s.wg.Wait()
}
