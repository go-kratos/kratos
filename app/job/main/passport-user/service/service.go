package service

import (
	"context"

	"go-common/app/job/main/passport-user/conf"
	"go-common/app/job/main/passport-user/dao"
	"go-common/app/job/main/passport-user/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
)

const (
	_asoAccountTable             = "aso_account"
	_asoAccountInfoTable         = "aso_account_info"
	_asoAccountSnsTable          = "aso_account_sns"
	_asoAccountRegOriginTable    = "aso_account_reg_origin"
	_insertAction                = "insert"
	_updateAction                = "update"
	_deleteAction                = "delete"
	_asoAccountInfoSharding      = 30
	_asoAccountRegOriginSharding = 20
	_mySQLErrCodeDuplicateEntry  = 1062
	_retry                       = 3
)

// Service service.
type Service struct {
	c *conf.Config
	d *dao.Dao
	// aso binlog consumer
	asoBinLogConsumer *databus.Databus
	group             *databusutil.Group
	// fullSync chan
	asoAccountChan     []chan *model.OriginAccount
	asoAccountInfoChan []chan *model.OriginAccountInfo
	asoAccountRegChan  []chan *model.OriginAccountReg
	asoAccountSnsChan  []chan *model.OriginAccountSns
	asoTelBindLogChan  []chan *model.UserTel

	countryMap map[int64]string
	aesKey     []byte
	salt       []byte
	//cron       *cron.Cron

	ch chan func()
}

// New new a service instance.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                  c,
		d:                  dao.New(c),
		asoAccountChan:     make([]chan *model.OriginAccount, c.FullSync.AsoAccount.ChanNum),
		asoAccountInfoChan: make([]chan *model.OriginAccountInfo, c.FullSync.AsoAccountInfo.ChanNum),
		asoAccountRegChan:  make([]chan *model.OriginAccountReg, c.FullSync.AsoAccountReg.ChanNum),
		asoAccountSnsChan:  make([]chan *model.OriginAccountSns, c.FullSync.AsoAccountSns.ChanNum),
		asoTelBindLogChan:  make([]chan *model.UserTel, c.FullSync.AsoTelBindLog.ChanNum),
		ch:                 make(chan func(), 1024000),
	}
	go s.cacheproc()
	if c.FullSync.AsoCountryCodeSwitch {
		err := s.syncAsoCountryCode()
		if err != nil {
			log.Error("fail to sync AsoCountryCode")
			panic(err)
		}
	}
	aesKey, err := s.d.AesKey(context.Background())
	if err != nil || aesKey == "" {
		log.Error("fail to get rsaKey")
		panic(err)
	}
	s.aesKey = []byte(aesKey)
	salt, err := s.d.Salt(context.Background())
	if err != nil || salt == "" {
		log.Error("fail to get salt")
		panic(err)
	}
	s.salt = []byte(salt)
	s.countryMap, err = s.d.CountryCodeMap(context.Background())
	if err != nil || len(s.countryMap) == 0 {
		log.Error("fail to get country map")
		panic(err)
	}
	if c.FullSync.AsoAccount.Switch {
		for i := 0; i < c.FullSync.AsoAccount.ChanNum; i++ {
			ch := make(chan *model.OriginAccount, c.FullSync.AsoAccount.ChanSize)
			s.asoAccountChan[i] = ch
			go s.asoAccountConsume(ch)
		}
		go s.getAsoAccount(c.FullSync.AsoAccount.Start, c.FullSync.AsoAccount.End, c.FullSync.AsoAccount.Count)
	}
	if c.FullSync.AsoAccountInfo.Switch {
		for i := 0; i < c.FullSync.AsoAccountInfo.ChanNum; i++ {
			ch := make(chan *model.OriginAccountInfo, c.FullSync.AsoAccountInfo.ChanSize)
			s.asoAccountInfoChan[i] = ch
			go s.asoAccountInfoConsume(ch)
		}
		go s.getAsoAccountInfo(c.FullSync.AsoAccountInfo.Start, c.FullSync.AsoAccountInfo.End, c.FullSync.AsoAccountInfo.Count)
	}
	if c.FullSync.AsoAccountReg.Switch {
		for i := 0; i < c.FullSync.AsoAccountReg.ChanNum; i++ {
			ch := make(chan *model.OriginAccountReg, c.FullSync.AsoAccountReg.ChanSize)
			s.asoAccountRegChan[i] = ch
			go s.asoAccountRegConsume(ch)
		}
		go s.getAsoAccountReg(c.FullSync.AsoAccountReg.Start, c.FullSync.AsoAccountReg.End, c.FullSync.AsoAccountReg.Count)
	}
	if c.FullSync.AsoAccountSns.Switch {
		for i := 0; i < c.FullSync.AsoAccountSns.ChanNum; i++ {
			ch := make(chan *model.OriginAccountSns, c.FullSync.AsoAccountSns.ChanSize)
			s.asoAccountSnsChan[i] = ch
			go s.asoAccountSnsConsume(ch)
		}
		go s.getAsoAccountSns(c.FullSync.AsoAccountSns.Start, c.FullSync.AsoAccountSns.End, c.FullSync.AsoAccountSns.Count)
	}
	if c.FullSync.AsoTelBindLog.Switch {
		for i := 0; i < c.FullSync.AsoTelBindLog.ChanNum; i++ {
			ch := make(chan *model.UserTel, c.FullSync.AsoTelBindLog.ChanSize)
			s.asoTelBindLogChan[i] = ch
			go s.asoTelBindLogConsume(ch)
		}
		go s.getAsoTelBindLog(c.FullSync.AsoTelBindLog.Start, c.FullSync.AsoTelBindLog.End, c.FullSync.AsoTelBindLog.Count)
	}
	if c.IncSync.Switch {
		s.asoBinLogConsumer = databus.New(c.DataBus.AsoBinLogSub)
		s.group = databusutil.NewGroup(
			c.DatabusUtil,
			s.asoBinLogConsumer.Messages(),
		)
		s.consumeproc()
	}
	//if c.Scheduler.Switch {
	//	s.cron = cron.New()
	//	if err := s.cron.AddFunc(c.Scheduler.EmailDuplicateCron, s.checkEmailDuplicateJob); err != nil {
	//		panic(err)
	//	}
	//	if err := s.cron.AddFunc(c.Scheduler.TelDuplicateCron, s.checkTelDuplicateJob); err != nil {
	//		panic(err)
	//	}
	//	s.cron.Start()
	//}
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return s.d.Ping(c)
}

// Close close service, including databus and outer service.
func (s *Service) Close() (err error) {
	if err = s.group.Close(); err != nil {
		log.Error("s.group.Close() error(%v)", err)
	}
	s.d.Close()
	return
}

func (s *Service) addCache(f func()) {
	select {
	case s.ch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.ch
		f()
	}
}
