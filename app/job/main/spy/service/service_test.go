package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/spy/conf"
	"go-common/app/job/main/spy/model"
	"go-common/library/cache/redis"

	"github.com/robfig/cron"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	c          = context.Background()
	s          *Service
	testCron   = "*/5 * * * * ?"
	testCyTime = 5000
)

func init() {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/spy-job-dev.toml")
	flag.Set("conf", dir)
	if err = conf.Init(); err != nil {
		panic(err)
	}
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

func CleanCache() {
	pool := redis.NewPool(conf.Conf.Redis.Config)
	pool.Get(c).Do("FLUSHDB")
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

func Test_cycleblock(t *testing.T) {
	Convey("Test_cycleblock cron", t, WithService(func(s *Service) {
		fmt.Println("Test_cycleblock start ")
		tx, err := s.dao.BeginTran(c)
		So(err, ShouldBeNil)
		ui := &model.UserInfo{Mid: testBlockMid, State: model.StateNormal}
		err = s.dao.TxUpdateUserState(c, tx, ui)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)

		lastBlockNo := s.lastBlockNo(s.c.Property.Block.CycleTimes)
		t := cron.New()
		err = t.AddFunc(testCron, s.cycleblock)
		if err != nil {
			panic(err)
		}
		t.Start()

		Convey("Test_cycleblock user info 1", WithService(func(s *Service) {
			var ui *model.UserInfo
			ui, err = s.dao.UserInfo(context.TODO(), testBlockMid)
			So(err, ShouldBeNil)
			So(ui.State == model.StateNormal, ShouldBeTrue)
		}))

		err = s.dao.AddBlockCache(c, testBlockMid, testLowScore, lastBlockNo)
		So(err, ShouldBeNil)
		mids, err := s.blockUsers(c, lastBlockNo)
		So(err, ShouldBeNil)
		So(mids, ShouldContain, testBlockMid)

		time.Sleep(5000 * time.Millisecond)
		time.Sleep(time.Duration(s.blockWaitTick))

		Convey("Test_cycleblock user info 2 ", WithService(func(s *Service) {
			ui, err := s.dao.UserInfo(context.TODO(), testBlockMid)
			So(err, ShouldBeNil)
			So(ui.State == model.StateBlock, ShouldBeTrue)
		}))
		fmt.Println("Test_cycleblock end ")
	}))
}

//  go test  -test.v -test.run TestStat
func TestStat(t *testing.T) {
	Convey(" UpdateStatData ", t, WithService(func(s *Service) {
		err := s.UpdateStatData(c, &model.SpyStatMessage{
			TargetMid: 1,
			TargetID:  1,
			EventName: "init_user_info",
			Type:      model.IncreaseStat,
			Quantity:  2,
			Time:      time.Now().Unix(),
			UUID:      "123456789qweasdzxcccccccc",
		})
		So(err, ShouldBeNil)
	}))
	Convey(" UpdateStatData 2", t, WithService(func(s *Service) {
		err := s.UpdateStatData(c, &model.SpyStatMessage{
			TargetMid: 1,
			TargetID:  1,
			EventName: "auto_block",
			Type:      model.ResetStat,
			Quantity:  2,
			Time:      time.Now().Unix(),
			UUID:      "123456789qweasdzxcccccccc",
		})
		So(err, ShouldBeNil)
	}))
}
