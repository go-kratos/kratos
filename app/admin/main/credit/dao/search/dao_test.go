package search

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/credit/conf"
	"go-common/app/admin/main/credit/model/blocked"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.credit-admin")
		flag.Set("conf_appid", "main.account-law.credit-admin")
		flag.Set("conf_token", "eKmbn2M4jvSyyjMEOywLFOQlX5ggRG9x")
		flag.Set("tree_id", "5885")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/convey-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestDao_Blocked(t *testing.T) {
	arg := &blocked.ArgBlockedSearch{Order: "id"}
	Convey("return someting", t, func() {
		ids, pagers, err := d.Blocked(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(ids), ShouldBeGreaterThanOrEqualTo, 0)
		So(pagers, ShouldNotBeNil)
	})
}

func TestDao_Publish(t *testing.T) {
	arg := &blocked.ArgPublishSearch{Order: "id"}
	Convey("return someting", t, func() {
		ids, pagers, err := d.Publish(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(ids), ShouldBeGreaterThanOrEqualTo, 0)
		So(pagers, ShouldNotBeNil)
	})
}

func Test_Case(t *testing.T) {
	arg := &blocked.ArgCaseSearch{Order: "id"}
	Convey("return someting", t, func() {
		ids, pagers, err := d.Case(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(ids), ShouldBeGreaterThanOrEqualTo, 0)
		So(pagers, ShouldNotBeNil)
	})
}

func Test_Jury(t *testing.T) {
	arg := &blocked.ArgJurySearch{Order: "id"}
	Convey("return someting", t, func() {
		ids, pagers, err := d.Jury(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(ids), ShouldBeGreaterThanOrEqualTo, 0)
		So(pagers, ShouldNotBeNil)
	})
}

func Test_Opinion(t *testing.T) {
	arg := &blocked.ArgOpinionSearch{Order: "id"}
	Convey("return someting", t, func() {
		ids, pagers, err := d.Opinion(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(ids), ShouldBeGreaterThanOrEqualTo, 0)
		So(pagers, ShouldNotBeNil)
	})
}

func Test_KPIPoint(t *testing.T) {
	arg := &blocked.ArgKpiPointSearch{Order: "id", Start: "2006-01-02 15:04:05", End: "2019-01-02 15:04:05"}
	Convey("return someting", t, func() {
		ids, pagers, err := d.KPIPoint(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(len(ids), ShouldBeGreaterThanOrEqualTo, 0)
		So(pagers, ShouldNotBeNil)
	})
}

func Test_BlockedUpdate(t *testing.T) {
	single := map[string]interface{}{
		"id":      1,
		"oper_id": 1,
		"status":  1,
		"black":   1,
	}
	var multiple []interface{}
	multiple = append(multiple, single)
	Convey("return someting", t, func() {
		err := d.SearchUpdate(context.TODO(), "block_info", "blocked_info", multiple)
		So(err, ShouldBeNil)
	})
}
