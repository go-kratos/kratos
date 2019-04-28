package manager

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/up/conf"
	"go-common/app/service/main/up/dao"
	"go-common/app/service/main/up/dao/global"
	"go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", dao.AppID)
		flag.Set("conf_token", dao.UatToken)
		flag.Set("tree_id", dao.TreeID)
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/up-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	global.Init(conf.Conf)
	d = New(conf.Conf)
	d.HTTPClient = blademaster.NewClient(conf.Conf.HTTPClient.Normal)
	m.Run()
	os.Exit(0)
}

func TestManagerPing(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestManagerprepareAndExec(t *testing.T) {
	var (
		c      = context.TODO()
		db     = d.managerDB
		sqlstr = _selectByID
		args   = interface{}(0)
	)
	convey.Convey("prepareAndExec", t, func(ctx convey.C) {
		res, err := prepareAndExec(c, db, sqlstr, args)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerprepareAndQuery(t *testing.T) {
	var (
		c      = context.TODO()
		db     = d.managerDB
		sqlstr = _selectByID
		args   = interface{}(0)
	)
	convey.Convey("prepareAndQuery", t, func(ctx convey.C) {
		rows, err := prepareAndQuery(c, db, sqlstr, args)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerGetUNamesByUids(t *testing.T) {
	var (
		c    = context.TODO()
		uids = []int64{100}
	)
	convey.Convey("GetUNamesByUids", t, func(ctx convey.C) {
		res, err := d.GetUNamesByUids(c, uids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestManagerGetUIDByNames(t *testing.T) {
	var (
		c     = context.TODO()
		names = []string{"wangzhe01"}
	)
	convey.Convey("GetUIDByNames", t, func(ctx convey.C) {
		res, err := d.GetUIDByNames(c, names)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
