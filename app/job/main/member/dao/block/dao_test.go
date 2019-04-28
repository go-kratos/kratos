package block

import (
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/job/main/member/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestTool(t *testing.T) {
	Convey("tool", t, func() {
		var (
			mids = []int64{1, 2, 3, 46333, 35858}
		)
		str := midsToParam(mids)
		So(str, ShouldEqual, "1,2,3,46333,35858")
	})
}

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.member-job")
		flag.Set("conf_token", "VEc5eqZNZHGQi6fsx7J6lJTqOGR9SnEO")
		flag.Set("tree_id", "2134")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/member-job-dev.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	config := conf.Conf
	//d.New().BlockImpl()
	d = New(conf.Conf, memcache.NewPool(config.Memcache.Config), xsql.NewMySQL(config.BlockDB), bm.NewClient(config.HTTPClient), nil)
	d.httpClient.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
