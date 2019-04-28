package service

import (
	"context"

	"go-common/app/service/main/passport-sns/conf"
	"go-common/app/service/main/passport-sns/dao"
	"go-common/app/service/main/passport-sns/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct of service
type Service struct {
	c         *conf.Config
	d         *dao.Dao
	snsLogPub *databus.Databus
	job       *fanout.Fanout
	cache     *fanout.Fanout

	AppMap          map[int]map[string]*model.SnsApps
	PlatformList    []int
	PlatformStrList []string
}

// New create new service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:               c,
		d:               dao.New(c),
		snsLogPub:       databus.New(c.DataBus.SnsLogPub),
		job:             fanout.New("job", fanout.Worker(10), fanout.Buffer(10240)),
		cache:           fanout.New("cache", fanout.Worker(10), fanout.Buffer(10240)),
		PlatformList:    []int{model.PlatformQQ, model.PlatformWEIBO},
		PlatformStrList: []string{model.PlatformQQStr, model.PlatformWEIBOStr},
	}
	var err error
	s.AppMap, err = s.initAppMap()
	if err != nil || s.AppMap == nil {
		log.Error("fail to get appMap")
		panic(err)
	}
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close dao.
func (s *Service) Close() {
	s.d.Close()
}

func (s *Service) initAppMap() (res map[int]map[string]*model.SnsApps, err error) {
	snsApps, err := s.d.SnsApps(context.Background())
	if err != nil {
		return
	}
	res = make(map[int]map[string]*model.SnsApps)
	platformMap := make(map[int][]*model.SnsApps)
	for _, app := range snsApps {
		platformMap[app.Platform] = append(platformMap[app.Platform], app)
	}
	for k, v := range platformMap {
		kMap := make(map[string]*model.SnsApps)
		for _, app := range v {
			kMap[app.AppID] = app
		}
		res[k] = kMap
	}
	return
}

func parsePlatform(platform string) int {
	switch platform {
	case model.PlatformQQStr:
		return model.PlatformQQ
	case model.PlatformWEIBOStr:
		return model.PlatformWEIBO
	}
	return 0
}

func parsePlatformStr(platform int) string {
	switch platform {
	case model.PlatformQQ:
		return model.PlatformQQStr
	case model.PlatformWEIBO:
		return model.PlatformWEIBOStr
	}
	return ""
}

func (s *Service) isAppID(platform int, appID string) bool {
	return s.AppMap[platform][appID] != nil
}

func platformToMidBindErr(platform int) error {
	switch platform {
	case model.PlatformQQ:
		return ecode.PassportSnsMidAlreadyBindQQ
	case model.PlatformWEIBO:
		return ecode.PassportSnsMidAlreadyBindWEIBO
	}
	return ecode.ServerErr
}

func platformToSnsBindErr(platform int) error {
	switch platform {
	case model.PlatformQQ:
		return ecode.PassportSnsQQAlreadyBind
	case model.PlatformWEIBO:
		return ecode.PassportSnsWEIBOAlreadyBind
	}
	return ecode.ServerErr
}
