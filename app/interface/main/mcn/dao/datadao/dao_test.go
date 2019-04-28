package datadao

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/dao/global"

	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.mcn-interface")
		flag.Set("conf_token", "49e4671bafbf93059aeb602685052ca0")
		flag.Set("tree_id", "58909")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/mcn-interface.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/mcn-interface.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	global.Init(conf.Conf)
	d = New(conf.Conf)
	d.Client.SetTransport(gock.DefaultTransport)
	var result = `{
"code": 200,
"msg": "success",
"result": [ ]
}`
	gock.New("http://berserker.bilibili.co/avenger/api").Get("/").AddMatcher(
		func(request *http.Request, request2 *gock.Request) (bool, error) {
			return true, nil
		}).Persist().Reply(200).JSON(result)
	defer gock.OffAll()
	os.Exit(m.Run())
}

// func httpMock(method, url string) *gock.Request {
// 	r := gock.New(url)
// 	r.Method = strings.ToUpper(method)
// 	return r
// }
