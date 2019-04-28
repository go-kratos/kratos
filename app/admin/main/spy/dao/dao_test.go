package dao

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"go-common/app/admin/main/spy/conf"
	"go-common/app/admin/main/spy/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dataMID   = int64(15555180)
	noDataMID = int64(1)
	d         *Dao
)

const (
	_cleanFactorSQL  = "delete from spy_factor where nick_name = ? AND service_id = ? AND event_id = ? AND risk_level = ?"
	_cleanEventSQL   = "delete from spy_event where name = ? AND service_id = ? AND status = ?"
	_cleanServiceSQL = "delete from spy_service where name = ? AND status = ?"
	_cleanGroupSQL   = "delete from spy_factor_group where name = ?"
)

func CleanMysql() {
	ctx := context.Background()
	d.db.Exec(ctx, _cleanFactorSQL, "test", 1, 1, 2)
	d.db.Exec(ctx, _cleanEventSQL, "test", 1, 1)
	d.db.Exec(ctx, _cleanServiceSQL, "test", 1)
	d.db.Exec(ctx, _cleanGroupSQL, "test")
}

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.spy-admin")
		flag.Set("conf_token", "bc3d60c2bb2b08a1b690b004a1953d3c")
		flag.Set("tree_id", "2857")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/spy-admin-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	CleanMysql()
	m.Run()
	os.Exit(0)
}

func WithMysql(f func(d *Dao)) func() {
	return func() {
		f(d)
	}
}

func Test_Mysql(t *testing.T) {
	Convey("get user info", t, WithMysql(func(d *Dao) {
		res, err := d.Info(context.TODO(), dataMID)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("get user info no data ", t, WithMysql(func(d *Dao) {
		res, err := d.Info(context.TODO(), noDataMID)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("get event history", t, WithMysql(func(d *Dao) {
		hpConf := &model.HisParamReq{
			Mid: 46333,
			Ps:  10,
			Pn:  1,
		}
		res, err := d.HistoryPage(context.TODO(), hpConf)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("get event history count", t, WithMysql(func(d *Dao) {
		hpConf := &model.HisParamReq{
			Mid: 46333,
			Ps:  10,
			Pn:  1,
		}
		res, err := d.HistoryPageTotalC(context.TODO(), hpConf)
		So(err, ShouldBeNil)
		fmt.Printf("history count : %d\n", res)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("get setting list", t, WithMysql(func(d *Dao) {
		res, err := d.SettingList(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("update setting", t, WithMysql(func(d *Dao) {
		list, err := d.SettingList(context.TODO())
		So(err, ShouldBeNil)
		setting := list[0]
		res, err := d.UpdateSetting(context.TODO(), setting.Val, setting.Property)
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
	Convey(" add factor ", t, WithMysql(func(d *Dao) {
		res, err := d.AddFactor(context.TODO(), &model.Factor{
			NickName:   "test",
			ServiceID:  int64(1),
			EventID:    int64(1),
			GroupID:    int64(1),
			RiskLevel:  int8(2),
			FactorVal:  float32(1),
			CategoryID: int8(1),
			CTime:      time.Now(),
			MTime:      time.Now(),
		})
		So(err, ShouldBeNil)
		So(res == 1, ShouldBeTrue)
	}))
	Convey(" add event ", t, WithMysql(func(d *Dao) {
		res, err := d.AddEvent(context.TODO(), &model.Event{
			Name:      "test",
			NickName:  "nickname",
			ServiceID: 1,
			Status:    1,
			CTime:     time.Now(),
			MTime:     time.Now(),
		})
		So(err, ShouldBeNil)
		So(res == 1, ShouldBeTrue)
	}))
	Convey(" add service ", t, WithMysql(func(d *Dao) {
		res, err := d.AddService(context.TODO(), &model.Service{
			Name:     "test",
			NickName: "nickname",
			Status:   1,
			CTime:    time.Now(),
			MTime:    time.Now(),
		})
		So(err, ShouldBeNil)
		So(res == 1, ShouldBeTrue)
	}))
	Convey(" add group ", t, WithMysql(func(d *Dao) {
		res, err := d.AddGroup(context.TODO(), &model.FactorGroup{
			Name:  "test",
			CTime: time.Now(),
		})
		So(err, ShouldBeNil)
		So(res == 1, ShouldBeTrue)
	}))
}
