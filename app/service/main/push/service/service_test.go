package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"go-common/app/service/main/push/conf"
	"go-common/app/service/main/push/model"
	"go-common/library/cache/redis"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	srv        *Service
	commonTask = &model.Task{
		ID:         "t123456",
		Job:        432674328,
		APPID:      model.APPIDBBPhone,
		BusinessID: 1,
		Platform:   []int{model.PlatformIPhone},
		Title:      "test task",
		Summary:    "test task summary",
		LinkType:   1,
		LinkValue:  "5",
		Build:      make(map[int]*model.Build),
		Sound:      1,
		Vibration:  1,
		MidFile:    "/path/to/file",
		Progress:   &model.Progress{},
		PushTime:   xtime.Time(time.Now().Unix()),
		ExpireTime: xtime.Time(time.Now().Unix() + 99999),
		Status:     model.TaskStatusPrepared,
	}
	commonReport = &model.Report{
		APPID:        model.APPIDBBPhone,
		PlatformID:   model.PlatformIPhone,
		Mid:          910819,
		Buvid:        "b",
		DeviceToken:  "dt",
		Build:        2233,
		TimeZone:     8,
		NotifySwitch: model.SwitchOn,
	}
)

func init() {
	dir, _ := filepath.Abs("../cmd/push-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(srv)
	}
}

func CleanCache() {
	c := context.TODO()
	pool := redis.NewPool(conf.Conf.Redis.Config)
	pool.Get(c).Do("FLUSHDB")
}

func Test_Setting(t *testing.T) {
	Convey("setting", t, WithService(func(s *Service) {
		var (
			c   = context.Background()
			mid = int64(91221505)
		)
		err := s.SetSetting(c, mid, model.UserSettingArchive, model.SwitchOff)
		So(err, ShouldBeNil)

		setting, err := s.Setting(c, mid)
		So(err, ShouldBeNil)
		st := model.Settings
		st[model.UserSettingArchive] = model.SwitchOff
		So(setting, ShouldResemble, st)
	}))

	Convey("get default setting", t, WithService(func(s *Service) {
		setting, err := s.Setting(context.TODO(), 8888888888888888)
		So(err, ShouldBeNil)
		So(setting, ShouldResemble, model.Settings)
	}))
}

func Test_Tokens(t *testing.T) {
	Convey("tokensByMids", t, WithService(func(s *Service) {
		mids := []int64{91221505, 91819}
		ts, missed, err := s.tokensByMids(context.TODO(), commonTask, mids)
		So(err, ShouldBeNil)
		t.Logf("misssed: %d", len(missed))
		t.Logf("tokens: %d", len(ts))
	}))
}

func Test_Pushs(t *testing.T) {
	Convey("pushs", t, WithService(func(s *Service) {
		task := &model.Task{
			ID:          "0",
			Type:        model.TaskTypeBusiness,
			APPID:       model.APPIDBBPhone,
			Platform:    []int{1, 2, 3},
			Title:       "test title",
			Summary:     "test summary",
			LinkType:    model.LinkTypeBrowser,
			LinkValue:   "https://www.bilibili.com/",
			Sound:       model.SwitchOn,
			Vibration:   model.SwitchOn,
			PushTime:    xtime.Time(time.Now().Unix()),
			ExpireTime:  xtime.Time(time.Now().Add(7 * 24 * time.Hour).Unix()),
			PassThrough: 0,
		}
		mids := []int64{2089809}
		err := s.Pushs(context.Background(), task, mids)
		So(err, ShouldBeNil)
	}))
}

func Test_Report(t *testing.T) {
	Convey("add report", t, WithService(func(s *Service) {
		err := s.AddReport(context.Background(), commonReport)
		So(err, ShouldBeNil)
	}))
}

func Test_AddUserReportCache(t *testing.T) {
	Convey("add user report cache", t, WithService(func(s *Service) {
		r1, r2 := *commonReport, *commonReport
		r1.DeviceToken = "device token 1"
		r2.DeviceToken = "device token 2"
		rs := []*model.Report{&r1, &r2}
		err := s.AddUserReportCache(context.TODO(), 91221505, rs)
		So(err, ShouldBeNil)
	}))
}

func Test_distribute(t *testing.T) {
	Convey("distribute tokens", t, WithService(func(s *Service) {
		task := *commonTask
		task.Build = map[int]*model.Build{2: {Build: 6190, Condition: "gte"}}
		rs := map[int64][]*model.Report{91819: {
			&model.Report{APPID: 1, PlatformID: 2, Mid: 91819, DeviceToken: "dt1", Build: 6190, NotifySwitch: 1},
			&model.Report{APPID: 1, PlatformID: 2, Mid: 91819, DeviceToken: "dt2", Build: 6191, NotifySwitch: 1},
		}}
		res := s.distribute(&task, rs)
		fmt.Println(len(res[2]))
		// fmt.Printf("%+v", res[2][0])
	}))
}

func Benchmark_Callback(b *testing.B) {
	Convey("callback", b, WithService(func(s *Service) {
		// for n := 0; n < b.N; n++ {
		// 	s.CallbackClick(context.TODO(), &model.Callback{
		// 		Type: model.CallbackTypeClick,
		// 	})
		// }
	}))
}

func Test_SinglePush(t *testing.T) {
	Convey("single push", t, WithService(func(s *Service) {
		s.SinglePush(context.TODO(), "", &model.Task{}, 91221505)
	}))
}

func Test_AddMidProgress(t *testing.T) {
	ctx := context.Background()
	Convey("add mid progress", t, WithService(func(s *Service) {
		id, err := s.dao.AddTask(ctx, commonTask)
		So(err, ShouldBeNil)
		idStr := strconv.FormatInt(id, 10)
		err = s.AddMidProgress(ctx, idStr, 100, 80)
		So(err, ShouldBeNil)
		task, err := s.dao.Task(ctx, idStr)
		So(err, ShouldBeNil)
		So(task, ShouldNotBeNil)
		So(task.Progress.MidTotal, ShouldEqual, 100)
		So(task.Progress.MidValid, ShouldEqual, 80)
	}))
}

func Test_dispatchByScheme(t *testing.T) {
	Convey("dispatch by scheme", t, func() {
		items := []*model.PushItem{
			{Platform: model.PlatformIPhone, Token: "token iphone", Build: 123456},
			{Platform: model.PlatformHuawei, Token: "token huawei", Build: 123456},
			{Platform: model.PlatformXiaomi, Token: "token xiaomi", Build: 888888888},
			{Platform: model.PlatformIPhone, Token: "token iphone 2", Build: 888888888},
		}
		res := dispatchByScheme(model.LinkTypeLive, "1,0", items)
		So(len(res), ShouldEqual, 2)
	})
}
