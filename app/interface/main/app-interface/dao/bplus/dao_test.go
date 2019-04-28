package bplus

import (
	"context"
	"flag"
	"go-common/app/interface/main/app-interface/conf"
	"os"
	"strings"
	"testing"
	// "go-common/app/interface/main/app-interface/model/search"

	"gopkg.in/h2non/gock.v1"
)

var (
	dao *Dao
)

func ctx() context.Context {
	return context.Background()
}

// TestMain dao ut.
func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-interface")
		flag.Set("conf_token", "1mWvdEwZHmCYGoXJCVIdszBOPVdtpXb3")
		flag.Set("tree_id", "2688")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-interface-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	dao = New(conf.Conf)
	os.Exit(m.Run())
	// time.Sleep(time.Second)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
