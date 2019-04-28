package relation

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"go-common/app/admin/main/credit/conf"

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

func TestDaoRPCStats(t *testing.T) {
	Convey("return someting", t, func() {
		res, err := d.RPCStats(context.TODO(), []int64{88889018, 2089809})
		for _, re := range res {
			fmt.Printf("re:%#v \n", re)
		}
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
