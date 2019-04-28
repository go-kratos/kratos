package upcrmservice

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/dao/email"
	"go-common/app/admin/main/up/dao/manager"
	"go-common/app/admin/main/up/dao/upcrm"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/app/admin/main/up/service/cache"
	"go-common/app/admin/main/up/service/data"
	"go-common/app/admin/main/up/util"
	"go-common/app/admin/main/up/util/timerqueue"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

//Service upcrm service
type Service struct {
	crmdb         *upcrm.Dao
	mng           *manager.Dao
	httpClient    *bm.Client
	mailService   *email.Dao
	uprankCache   map[int][]*upcrmmodel.UpRankInfo
	lastCacheDate time.Time
	dataService   *data.Service
}

//New create service
func New(c *conf.Config) (svc *Service) {
	svc = &Service{
		crmdb:       upcrm.New(c),
		mng:         manager.New(c),
		mailService: email.New(c),
		dataService: data.New(c),
	}
	svc.initSvc()
	return svc
}

func (s *Service) initSvc() {
	addScheduleWithConf(conf.Conf.TimeConf.RefreshUpRankTime, s.RefreshCache, "03:00:00", "refresh up rank data")
	go s.RefreshCache(time.Now())

	cache.LoadCache()
	var cacheInterval = 60 * time.Minute
	util.GlobalTimer.ScheduleRepeat(timerqueue.NewTimerWrapper(cache.RefreshUpTypeAsync), time.Now().Add(cacheInterval), cacheInterval)
	util.GlobalTimer.ScheduleRepeat(timerqueue.NewTimerWrapper(cache.ClearTagCache), time.Now().Add(cacheInterval), cacheInterval)
}

func addScheduleWithConf(scheduleTime string, timerFunc timerqueue.TimerFunc, defaultTime string, desc string) {
	if scheduleTime == "" {
		scheduleTime = defaultTime
	}
	var next, err = util.GetNextPeriodTime(scheduleTime, time.Hour*24, time.Now())
	if err != nil {
		panic(fmt.Sprintf("config for time fail, err=%+v", err))
	}
	log.Info("[%s] next period time is %+v", desc, next)
	util.GlobalTimer.ScheduleRepeat(timerqueue.NewTimerWrapper(timerFunc), next, time.Hour*24)
}

//SetHTTPClient set client
func (s *Service) SetHTTPClient(client *bm.Client) {
	s.httpClient = client
	s.crmdb.SetHTTPClient(client)
}

//DataService data service
func (s *Service) DataService() *data.Service {
	return s.dataService
}

//RefreshCache refresh cache
func (s *Service) RefreshCache(tm time.Time) {
	latestDate, err := s.crmdb.GetUpRankLatestDate()
	if err != nil {
		log.Error("get latest rank time from db fail, err=%+v", err)
		return
	}
	if latestDate == s.lastCacheDate {
		log.Info("no need to refresh cache, latest cache date=%v", latestDate)
		return
	}

	s.refreshUpRankDate(latestDate)
}

//TestGetViewBase test get view
func (s *Service) TestGetViewBase(c context.Context, arg *upcrmmodel.TestGetViewBaseArgs) (res interface{}, err error) {
	var dataMap = make(map[string]interface{})
	dataMap["info"], err = s.dataService.GetViewData(c, arg.Mid)
	res = dataMap
	return
}
