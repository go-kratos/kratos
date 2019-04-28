package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/passport-user-compare/conf"
	"go-common/app/job/main/passport-user-compare/dao"
	"go-common/app/job/main/passport-user-compare/model"
	"go-common/library/log"

	"github.com/robfig/cron"
)

const (
	publicKeyConst             int8  = 0
	privateKeyConst            int8  = 1
	aesKeyConst                int8  = 2
	md5KeyConst                int8  = 3
	timeFormat                       = "2006-01-02 15:04:05"
	pwdErrorType                     = 1
	statusErrorType                  = 2
	telErrorType                     = 3
	mailErrorType                    = 4
	safeErrorType                    = 5
	sinaErrorType                    = 6
	qqErrorType                      = 7
	notExistUserBase                 = 8
	notExistUserTel                  = 9
	notExistUserMail                 = 10
	notExistUserSafeQuestion         = 11
	notExistUserThirdBind            = 12
	userIDErrorType                  = 13
	notExistUserRegOriginType        = 14
	userRegOriginErrorType           = 15
	insertAction                     = "insert"
	updateAction                     = "update"
	fullDivSegment             int64 = 10
	platformSina               int64 = 1
	platformQQ                 int64 = 2
	filterStart                int64 = 1536572121
	filterEnd                  int64 = 1536616436
	mySQLErrCodeDuplicateEntry       = 1062
)

var (
	privateKey             = ""
	publicKey              = ""
	aesKey                 = ""
	md5slat                = ""
	dynamicAccountStat     = make(map[string]int64, 7)
	dynamicAccountInfoStat = make(map[string]int64, 2)
	dynamicAccountRegStat  = make(map[string]int64, 2)
	accountStat            = make(map[string]int64, 7)
	accountInfoStat        = make(map[string]int64, 2)
	accountSnsStat         = make(map[string]int64, 3)
	loc                    = time.Now().Location()
)

// Service service.
type Service struct {
	c                 *conf.Config
	d                 *dao.Dao
	cron              *cron.Cron
	fullFixChan       chan *model.ErrorFix
	incFixChan        chan *model.ErrorFix
	mu                sync.Mutex
	countryMap        map[int64]string
	dataFixSwitch     bool
	incrDataFixSwitch bool
}

// New new service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                 c,
		d:                 dao.New(c),
		cron:              cron.New(),
		dataFixSwitch:     c.DataFixSwitch,
		incrDataFixSwitch: c.IncrDataFixSwitch,
		fullFixChan:       make(chan *model.ErrorFix, 1024),
		incFixChan:        make(chan *model.ErrorFix, 2048),
	}
	var err error
	s.countryMap, err = s.d.QueryCountryCode(context.Background())
	if err != nil || len(s.countryMap) == 0 {
		log.Error("fail to get country map")
		panic(err)
	}

	s.initSecret()
	if s.c.FullTask.Switch {
		log.Info("start full compare and fixed")
		go s.fullFixed(s.fullFixChan)
		if err := s.cron.AddFunc(s.c.FullTask.CronFullStr, s.fullCompareAndFix); err != nil {
			panic(err)
		}
		s.cron.Start()
	}

	if s.c.IncTask.Switch {
		log.Info("start inc compare and fixed")
		go s.incFixed(s.incFixChan)
		if err := s.cron.AddFunc(s.c.IncTask.CronIncStr, s.incCompareAndFix); err != nil {
			panic(err)
		}
		s.cron.Start()
	}

	if s.c.DuplicateTask.Switch {
		log.Info("start check duplicate job")
		if err := s.cron.AddFunc(s.c.DuplicateTask.DuplicateCron, s.checkEmailDuplicateJob); err != nil {
			panic(err)
		}
		if err := s.cron.AddFunc(s.c.DuplicateTask.DuplicateCron, s.checkTelDuplicateJob); err != nil {
			panic(err)
		}
		s.cron.Start()
	}
	if s.c.FixEmailVerifiedSwitch {
		go s.fixEmailVerified()
	}
	return s
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close close service, including databus and outer service.
func (s *Service) Close() (err error) {
	if err = s.d.Close(); err != nil {
		panic(err)
	}
	return
}

func (s *Service) initSecret() {
	res, err := s.d.LoadSecret(context.Background())
	if err != nil {
		panic(err)
	}
	for _, r := range res {
		if publicKeyConst == r.KeyType {
			publicKey = r.Key
		}
		if privateKeyConst == r.KeyType {
			privateKey = r.Key
		}
		if aesKeyConst == r.KeyType {
			aesKey = r.Key
		}
		if md5KeyConst == r.KeyType {
			md5slat = r.Key
		}
	}
	if len(aesKey) == 0 {
		panic("load secret error")
	}
}
