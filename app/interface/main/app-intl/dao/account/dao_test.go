package account

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-intl/conf"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-intl")
		flag.Set("conf_token", "02007e8d0f77d31baee89acb5ce6d3ac")
		flag.Set("tree_id", "64518")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-intl-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestRelations3(t *testing.T) {
	Convey(t.Name(), t, func() {
		d.Relations3(context.Background(), []int64{112}, 123)
	})
}

func TestIsAttention(t *testing.T) {
	Convey(t.Name(), t, func() {
		d.IsAttention(context.Background(), []int64{12}, 23)
	})
}

func TestCard3(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Card3(context.Background(), 12)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestCards3(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Cards3(context.Background(), []int64{1})
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestFollowing3(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Following3(context.Background(), 12, 13)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestInfos3(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Infos3(context.Background(), []int64{12})
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
