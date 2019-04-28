package gorpc

import (
	"context"
	"testing"

	pushsrv "go-common/app/service/main/push/api/gorpc"
	"go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	// _noArg = &struct{}{}
	// _noRes = &struct{}{}
	ctx = context.TODO()
)

func WithRPC(f func(client *pushsrv.Service)) func() {
	return func() {
		client := pushsrv.New(nil)
		f(client)
	}
}

func Test_AddReport(t *testing.T) {
	Convey("AddReport", t, WithRPC(func(client *pushsrv.Service) {
		arg := &model.ArgReport{
			APPID:        1,
			PlatformID:   1,
			Mid:          1,
			Buvid:        "b",
			DeviceToken:  "d",
			Build:        8080,
			TimeZone:     8,
			NotifySwitch: 1,
		}
		err := client.AddReport(ctx, arg)
		So(err, ShouldBeNil)
	}))
}

func Test_Setting(t *testing.T) {
	Convey("get setting", t, WithRPC(func(client *pushsrv.Service) {
		arg := &model.ArgMid{Mid: 88888888}
		res, err := client.Setting(ctx, arg)
		So(err, ShouldBeNil)
		t.Logf("setting(%v)", res)
	}))

	Convey("set setting", t, WithRPC(func(client *pushsrv.Service) {
		arg := &model.ArgSetting{Mid: 999999999, Type: model.UserSettingArchive, Value: model.SwitchOff}
		err := client.SetSetting(ctx, arg)
		So(err, ShouldBeNil)

		argMid := &model.ArgMid{Mid: 999999999}
		res, err := client.Setting(ctx, argMid)
		So(err, ShouldBeNil)
		t.Logf("setting(%v)", res)
	}))
}

func TestAddUserReportCache(t *testing.T) {
	Convey("AddUserReportCache", t, WithRPC(func(client *pushsrv.Service) {
		arg := &model.ArgUserReports{Mid: 123456, Reports: []*model.Report{{
			APPID:       1,
			PlatformID:  1,
			Mid:         123456,
			DeviceToken: "dtrpc",
		}}}
		err := client.AddUserReportCache(context.Background(), arg)
		So(err, ShouldBeNil)
	}))
}

func TestAddTokensCache(t *testing.T) {
	Convey("AddTokensCache", t, WithRPC(func(client *pushsrv.Service) {
		arg := &model.ArgReports{Reports: []*model.Report{{
			APPID:       1,
			PlatformID:  1,
			Mid:         123456,
			DeviceToken: "dtrpc",
		}, {
			APPID:       1,
			PlatformID:  1,
			Mid:         123456,
			DeviceToken: "dtrpc2",
		}}}
		err := client.AddTokensCache(context.Background(), arg)
		So(err, ShouldBeNil)
	}))
}
