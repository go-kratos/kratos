package weeklyhonor

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"go-common/app/interface/main/creative/conf"
	upDao "go-common/app/interface/main/creative/dao/up"
	honDao "go-common/app/interface/main/creative/dao/weeklyhonor"
	honmdl "go-common/app/interface/main/creative/model/weeklyhonor"
	"go-common/app/interface/main/creative/service"
	upmdl "go-common/app/service/main/up/model"
	xtime "go-common/library/time"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	err := conf.Init()
	if err != nil {
		fmt.Printf("conf.Init() error(%v)", err)
	}
	rpcdaos := service.NewRPCDaos(conf.Conf)
	s = New(conf.Conf, rpcdaos)
}

func WithService(f func(s *Service)) func() {
	return func() {
		convey.Reset(func() {})
		f(s)
	}
}

func TestService_WeeklyHonor(t *testing.T) {
	var (
		mid      int64 = 1
		uid      int64 = 2
		token          = s.genToken(2)
		mockStat       = honmdl.HonorStat{
			Play:      100,
			PlayLastW: 100,
			Fans:      100,
			FansInc:   2000,
			CoinInc:   2000,
			Rank3:     50,
		}
		mockHls = map[int]*honmdl.HonorLog{
			12: {
				ID:    1,
				MID:   uid,
				HID:   12,
				Count: 1,
				MTime: xtime.Time(time.Now().AddDate(0, 0, -5).Unix()),
			},
			23: {
				ID:    1,
				MID:   uid,
				HID:   12,
				Count: 1,
				MTime: xtime.Time(honmdl.LatestSunday().Add(time.Hour).Unix()),
			},
		}
	)
	convey.Convey("test", t, WithService(func(s *Service) {
		// mock
		s.c.HonorDegradeSwitch = false
		monkey.PatchInstanceMethod(reflect.TypeOf(s.up), "UpInfo", func(_ *upDao.Dao, _ context.Context, _ int64, _ int, _ string) (*upmdl.UpInfo, error) {
			return &upmdl.UpInfo{IsAuthor: 1}, nil
		})
		monkey.PatchInstanceMethod(reflect.TypeOf(s.honDao), "StatMC", func(_ *honDao.Dao, _ context.Context, _ int64, _ string) (*honmdl.HonorStat, error) {
			return &mockStat, nil
		})
		monkey.PatchInstanceMethod(reflect.TypeOf(s.honDao), "HonorMC", func(_ *honDao.Dao, _ context.Context, _ int64, _ string) (*honmdl.HonorLog, error) {
			return nil, nil
		})
		monkey.PatchInstanceMethod(reflect.TypeOf(s.honDao), "HonorLogs", func(_ *honDao.Dao, _ context.Context, _ int64) (map[int]*honmdl.HonorLog, error) {
			return mockHls, nil
		})
		defer monkey.UnpatchAll()
		// test guest visit
		honor, err := s.WeeklyHonor(context.Background(), mid, uid, token)
		convey.ShouldBeNil(err)
		convey.So(honor.HID, convey.ShouldEqual, 20)
		convey.So(honor.SubState, convey.ShouldEqual, 0)
		// test up visit
		monkey.PatchInstanceMethod(reflect.TypeOf(s.honDao), "GetUpSwitch", func(_ *honDao.Dao, _ context.Context, _ int64) (uint8, error) {
			return 1, nil
		})
		honor, err = s.WeeklyHonor(context.Background(), uid, uid, token)
		convey.ShouldBeNil(err)
		convey.So(honor.HID, convey.ShouldEqual, 20)
		convey.So(honor.SubState, convey.ShouldEqual, 1)
	}))
}
