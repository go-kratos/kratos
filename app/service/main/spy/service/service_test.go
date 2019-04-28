package service

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/model"
	"go-common/library/cache/redis"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.TODO()
	// 35858 : bind nothing
	// 22333 : bind mail only
	// 46333 : bind both
	// 27515232 : bind both
	mids = []int64{27515232}
)

func TestMain(m *testing.M) {
	flag.Parse()
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	if s == nil {
		s = New(conf.Conf)
	}
	os.Exit(m.Run())
}

func TestDelCache(t *testing.T) {
	var err error
	for _, mid := range mids {
		if err = s.dao.DelInfoCache(c, mid); err != nil {
			t.Fatal(err)
		}
	}
}

func TestUserInfo(t *testing.T) {
	var (
		ip  = "127.0.0.1"
		err error
		ui  *model.UserInfo
	)

	for _, mid := range mids {
		if ui, err = s.UserInfo(c, mid, ip); err != nil {
			t.Fatal(err)
		}
		if ui == nil {
			t.Fatal()
		}
		t.Logf("userinfo (%+v)", ui)
	}
}

func TestHandleEvent(t *testing.T) {
	var (
		arg struct {
			Test string
		}
		err           error
		ui            *model.UserInfo
		preEventScore int8
		preBaseScore  int8
		fac           *model.Factor
	)
	arg.Test = "test arg"

	for _, mid := range mids {
		var fakeEventMsg = &model.EventMessage{
			IP:        "127.0.0.1",
			Service:   "spy_service",
			Event:     "bind_nothing",
			ActiveMid: mid,
			TargetMid: mid,
			Result:    "test_result",
			Effect:    "test_effect",
			RiskLevel: 1,
		}
		fakeEventMsg.Args = arg
		if ui, err = s.UserInfo(c, fakeEventMsg.ActiveMid, fakeEventMsg.IP); err != nil {
			t.Fatal(err)
		}
		if ui == nil {
			t.Fatal()
		}
		preBaseScore = ui.BaseScore
		preEventScore = ui.EventScore

		if err = s.HandleEvent(c, fakeEventMsg); err != nil {
			t.Fatal(err)
		}
		if ui, err = s.UserInfo(c, fakeEventMsg.ActiveMid, fakeEventMsg.IP); err != nil {
			t.Fatal(err)
		}
		if ui.State == model.StateBlock {
			continue
		}
		if ui.BaseScore != preBaseScore {
			t.Fatal()
		}
		if fac, err = s.factor(c, fakeEventMsg.Service, fakeEventMsg.Event, 1); err != nil {
			t.Fatal(err)
		}
		if ui.EventScore != int8(float64(preEventScore)*fac.FactorVal) {
			t.Fatal()
		}
		if ui.Score != int8(float64(ui.BaseScore)*float64(ui.EventScore)/float64(100)) {
			t.Fatal()
		}
		if err = s.dao.DelInfoCache(c, mid); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second * 1)
}

func CleanCache() {
	pool := redis.NewPool(conf.Conf.Redis.Config)
	pool.Get(c).Do("FLUSHDB")
}

func init() {
	dir, _ := filepath.Abs("../cmd/spy-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(s)
	}
}

func Test_LoadSystemConfig(t *testing.T) {
	Convey("Test_LoadSystemConfig had data", t, WithService(func(s *Service) {
		fmt.Println(s.spyConfig)
		So(s.spyConfig, ShouldContainKey, model.LimitBlockCount)
		So(s.spyConfig, ShouldContainKey, model.LessBlockScore)
		So(s.spyConfig, ShouldContainKey, model.AutoBlock)
	}))
}

// go test  -test.v -test.run TestStatList
func TestStatList(t *testing.T) {
	var (
		mid int64 = 1
		id  int64 = 1
	)
	Convey("TestStatList ", t, WithService(func(s *Service) {
		stat, err := s.StatByID(c, mid, id)
		So(err, ShouldBeNil)
		fmt.Println("stat", stat)
	}))
}

// go test  -test.v -test.run TestTelLevel
func TestTelLevel(t *testing.T) {
	var (
		mid int64 = 15555180
	)
	Convey("TestTelLevel no err", t, WithService(func(s *Service) {
		level, err := s.TelRiskLevel(c, mid, "127.0.0.1")
		fmt.Println("level", level)
		So(err, ShouldBeNil)
	}))
}

// go test  -test.v -test.run TestTelLevel
func TestService_StatByIDGroupEvent(t *testing.T) {
	var (
		mid int64 = 7455
	)
	Convey("StatByIDGroupEvent no err", t, WithService(func(s *Service) {
		res, err := s.StatByIDGroupEvent(c, mid, 0)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}
