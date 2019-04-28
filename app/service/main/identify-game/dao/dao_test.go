package dao

import (
	"flag"
	"go-common/app/service/main/identify-game/conf"
	"os"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "dev" {
		flag.Set("app_id", "main.account.identify-game")
		flag.Set("conf_token", "ba87634784874cad05e941bcabd55512")
		flag.Set("tree_id", "identify-game")
		flag.Set("tree_id", "11799")
		flag.Set("conf_version", "server-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.client.SetTransport(gock.DefaultTransport)
	m.Run()
	os.Exit(0)
}

//func httpMock(method, url string) *gock.Request {
//	r := gock.New(url)
//	r.Method = strings.ToUpper(method)
//	return r
//}
