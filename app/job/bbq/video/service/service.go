package service

import (
	"context"
	"flag"
	"fmt"
	"go-common/library/conf/env"
	"go-common/library/queue/databus"
	"reflect"
	"time"

	"go-common/app/job/bbq/video/conf"
	"go-common/app/job/bbq/video/dao"
	"go-common/library/log"

	topic "go-common/app/service/bbq/topic/api"

	"github.com/robfig/cron"
)

var (
	srvName string
)

// Service struct
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	searchChan  chan string
	videoSub    *databus.Databus
	videoRep    *databus.Databus
	bvcSub      *databus.Databus
	scheFunc    map[string]func()
	topicClient topic.TopicClient
}

func init() {
	flag.StringVar(&srvName, "srv", "", "service name")
}

func newTopicClient() topic.TopicClient {
	topicClient, err := topic.NewClient(nil)
	if err != nil {
		log.Errorw(context.Background(), "log", "get topic client fail")
		panic(err)
	}
	return topicClient
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		dao:         dao.New(c),
		searchChan:  make(chan string, 1),
		topicClient: newTopicClient(),
	}
	s.scheFunc = s.initScheduleFunc()
	if srvName != "" {
		switch srvName {
		case "test":
			s.Test()
		case "syncsv2es":
			s.taskSyncVideo2ES()
		case "syncuserdmg":
			s.taskSyncUserDmg()
		case "rminvalides":
			s.taskRmInvalidES()
		case "regcmtall":
			s.AutoRegAll(context.Background())
		case "syncuserbase":
			s.taskSyncUsrBaseFromVideo(context.Background())
		case "syncuserbasic":
			s.taskSyncPegasusUserBasic()
		case "syncsearchvideo":
			s.SyncVideo2Search()
		case "syncsearchuser":
			s.SyncUser2Search()
		case "syncsearchsug":
			s.SyncSug2Search()
		case "upubface":
			s.UpdateUsrBaseFace()
		case "SysMsgTask":
			s.SysMsgTask()
		case "UserProfile":
			s.UserProfileUpdate()
		case "pushbvc":
			s.commitCID()
		case "cmscheckback":
			s.TransToCheckBack()
		}
		return s
	}
	//初始化databus
	s.initDatabus()
	//启动相关rountine
	s.launchCor()
	//定时任务启动
	if env.DeployEnv == env.DeployEnvProd {
		s.runScheduler(c.Scheduler)
	}
	return s
}

//runScheduler .1
func (s *Service) runScheduler(c *conf.Scheduler) {
	sche := cron.New()

	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)
	for i := 0; i < v.NumField(); i++ {
		//排除配置为空的任务
		if job := v.Field(i).String(); job != "" {
			fn := t.Field(i).Name
			//从映射集中取出对应函数
			if f, ok := s.scheFunc[fn]; !ok {
				fmt.Printf("skip[%s]\n", fn)
				continue
			} else {
				fmt.Printf("run[%s]\n", fn)
				if err := sche.AddFunc(job, f); err != nil {
					panic(err)
				}
			}
		}
	}
	sche.Start()
}

//Test 测试
func (s *Service) Test() {
	log.Info("HeartBeat:%s", time.Now())
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// initScheduleFunc 任务结构与函数映射关系
// map {key: conf.Scheduler结构体中字段名 value: 对应执行的函数}
func (s *Service) initScheduleFunc() map[string]func() {
	return map[string]func(){
		"Test":             s.Test,
		"CheckVideo2ES":    s.taskSyncVideo2ES,
		"SyncUserDmg":      s.taskSyncUserDmg,
		"SyncUpUserDmg":    s.taskSyncUpUserDmg,
		"CheckVideo":       s.taskCheckVideo,
		"CheckVideoSt":     s.taskCheckVideoStatistics,
		"CheckVideoStHv":   s.taskCheckVideoStatisticsHive,
		"CheckVideoTag":    s.taskCheckVideoTag,
		"CheckTag":         s.taskCheckTag,
		"SyncUsrSta":       s.taskSyncUsrStaFromHive,
		"SysMsgTask":       s.SysMsgTask,
		"UserProfileBbq":   s.UserProfileUpdate,
		"TransToReview":    s.TransToReview,
		"TransToCheckBack": s.TransToCheckBack,
	}
}

func (s *Service) launchCor() {
	time.Sleep(time.Second * 3)
	if env.DeployEnv == env.DeployEnvProd {
		go s.SyncSearch()
	}
	go s.videoBinlogSub()
	go s.videoRepositoryBinlogSub()
	go s.BvcTransSub()
}

func (s *Service) initDatabus() {
	s.videoSub = databus.New(conf.Conf.Databus["videosub"])
	s.videoRep = databus.New(conf.Conf.Databus["videorep"])
	s.bvcSub = databus.New(conf.Conf.Databus["bvcsub"])
}
