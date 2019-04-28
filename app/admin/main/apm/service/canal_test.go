package service

import (
	"context"
	"net/url"
	"testing"

	"go-common/app/admin/main/apm/model/canal"

	"github.com/BurntSushi/toml"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	cookie      = "username=fengshanshan; _AJSESSIONID=ee7c557c29e9b3f405b21c49a9aaa72a; sven-apm=9fc2bd4165a90a21df5803b4abc5cf7a81b7564e5119a8a547fd763bda757ceb"
	_jsonstring = `[{ "schema":"123","table":[  {"name":"abc","primarykey":["order_id","new_id"],"omitfield":["new","old"]} , {"name":"def","primarykey":["order_id","new_id"],"omitfield":["new","old"] } ,{"name":"sfg","primarykey":["order_id","new_id"],"omitfield":["new","old"]} ],"infoc":{"taskID":"000846","proto":"tcp","addr":"172.19.100.20:5401","reporterAddr":"172.19.40.195:6200"}},{ "schema":"456","table":[  {"name":"abc" ,"primarykey":["order_id","new_id"],"omitfield":["new","old"]} , {"name":"def" } ,{"name":"sfg"} ],"databus": { "group": "LiveTime-LiveLive-P","addr": "172.16.33.158:6205"}}]`
)

func TestService_GetAllErrors(t *testing.T) {
	Convey("test getAll errors", t, func() {
		//svr = New(conf.Conf)
		errS, err := svr.GetAllErrors(context.Background())
		So(err, ShouldBeNil)
		t.Log("getAllErrors=", errS)
	})
}

func TestGetCanalInstance(t *testing.T) {
	Convey("TestGetCanalInstance", t, func() {
		host, err := svr.getCanalInstance(context.Background())
		So(err, ShouldBeNil)
		t.Log("canalinstance=", host)
	})
}

func TestCheckMaster(t *testing.T) {
	Convey("TestCheckMaster", t, func() {
		req := &canal.ConfigReq{
			Addr:     "172.16.33.205:3310",
			User:     "gosync",
			Password: "xmJ4KlRuXzc9UfNerOGHP0LpaS26VqoM",
		}
		err := svr.CheckMaster(context.Background(), req)
		So(err, ShouldNotBeNil)
		t.Log("err=", err)
	})
}

func TestScanByAddrFromConfig(t *testing.T) {
	Convey("test ScanByAddrFromConfig", t, func() {

		//svr = New(conf.Conf)
		res, err := svr.ScanByAddrFromConfig(context.Background(), "172.16.33.205:3308", cookie)
		var document canal.Document
		data, _ := toml.Decode(res.Comment, &document)
		t.Log("toml=", data)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)

	})
}
func TestScanInfo(t *testing.T) {
	Convey("TestScanInfo", t, func() {
		v := &canal.ScanReq{
			Addr: "172.16.33.223:3308",
		}
		res, err := svr.GetScanInfo(context.Background(), v, "fss", cookie)
		t.Logf("toml=%+v", res.Document)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)

	})
}

func TestGetBuildID(t *testing.T) {
	Convey("test getBuildID", t, func() {

		res, err := svr.getBuildID(context.Background(), cookie)
		t.Log("buidID=", res)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestGetConfigID(t *testing.T) {
	Convey("test getConfigID", t, func() {
		res, err := svr.getConfigID(context.Background(), "13", cookie)
		t.Log("ConfigID=", res)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestGetConfigValue(t *testing.T) {
	Convey("test getConfigValue", t, func() {
		res, err := svr.getConfigValue(context.Background(), "112", cookie)
		t.Log("ConfigValue=", res)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestConfigByName(t *testing.T) {
	Convey("test configsByName", t, func() {
		params := "172.16.33.205:3308"
		res, err := svr.getConfigsByName(context.Background(), params, cookie)
		t.Logf("configsByName:%+v,%+v", res, err)
		So(res, ShouldNotBeEmpty)
	})
}

func TestProcessCanalInfo(t *testing.T) {
	Convey("TestProcessCanalInfo", t, func() {
		v := &canal.ConfigReq{
			Addr:     "127.0.0.1:8001",
			User:     "admin",
			Password: "admin",
			// Project:   "main.web-svr",
			// Leader:    "fss",
			Databases: _jsonstring,
			Mark:      "sfs",
		}
		_ = svr.ProcessCanalInfo(context.Background(), v, "fengshanshan")
	})
}

func TestJointConfigInfo(t *testing.T) {
	Convey("TestJointConfigInfo", t, func() {
		v := &canal.ConfigReq{
			Addr:          "172.16.33.999:3308",
			User:          "admin",
			Password:      "admin",
			MonitorPeriod: "2h",
			Project:       "main.web-svr",
			Leader:        "fss",
			Databases:     _jsonstring,
			Mark:          "sfs",
		}
		comment, _ := svr.jointConfigInfo(context.Background(), v, cookie)
		t.Logf("comment:%s", comment)
	})
}

func TestCreateConfig(t *testing.T) {
	Convey("test createConfig", t, func() {
		params := url.Values{}
		params.Set("comment", "23232ew")
		params.Set("name", "demo1214s2.toml")
		params.Set("state", "2")
		params.Set("mark", "demo")
		params.Set("from", "0")
		params.Set("user", "sada")
		res, err := svr.createConfig(context.Background(), params, cookie)
		if err != nil {
			So(err, ShouldContain, "643")
		}
		So(res, ShouldBeNil)

		t.Log("createConfig=", res)
	})
}

func TestUpdateConfig(t *testing.T) {
	Convey("test updateConfig", t, func() {
		params := url.Values{}
		params.Set("config_id", "364")
		params.Set("mtime", "1528975519")
		params.Set("state", "2")
		params.Set("mark", "update")
		params.Set("comment", "dsdfasdsads")
		res, err := svr.updateConfig(context.Background(), params, cookie)
		if err != nil {
			So(err, ShouldContain, "643")
		}
		So(res, ShouldBeNil)

		t.Log("updateConfig=", res)
	})
}

func TestGetGroupInfo(t *testing.T) {
	Convey("test getGroupInfo", t, func() {
		res, _, err := svr.getGroupInfo("test-group")
		t.Log("GroupInfo=", res)
		So(err, ShouldBeNil)
	})
}

func TestGetAppInfo(t *testing.T) {
	Convey("test getAppInfo", t, func() {
		res, err := svr.getAppInfo("test-group")
		t.Log("AppInfo=", res)
		So(err, ShouldBeNil)
	})
}

func TestGetAction(t *testing.T) {
	Convey("test getAction", t, func() {
		res := svr.getAction("LiveTimeS-LiveLivePÍ")
		t.Log("Action=", res)
	})
}

func TestGetTableInfo(t *testing.T) {
	Convey("test getTableInfo", t, func() {
		table := "172.16.33.12,172.16.33.48"
		res, err := svr.TableInfo(table)
		t.Log("table=", res)
		So(err, ShouldBeNil)
	})
}

func TestGetDatabusInfo(t *testing.T) {
	Convey("test getDatabusInfo", t, func() {
		tab1 := &canal.Table{
			Name:       "abc",
			Primarykey: []string{"new", "old"},
			Omitfield:  []string{"new", "old"},
		}
		tables := []*canal.Table{tab1}
		res, err := svr.databusInfo("LiveTime-LiveLive-P", "172.16.33.22", "relation", tables)
		t.Logf("res:%+v", res)
		So(err, ShouldBeNil)
	})
}

func TestUpdateProcessTag(t *testing.T) {
	Convey("test UpdateProcessTag", t, func() {
		err := svr.UpdateProcessTag(context.Background(), 355, cookie)
		So(err, ShouldBeNil)
	})
}

func TestGetServerID(t *testing.T) {
	Convey("test TestGetServerID", t, func() {
		sid, err := svr.getServerID("172.16.448.789:9000")
		t.Logf("sid:%v\n", sid)
		So(sid, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
func TestServiceSendWechatMessage(t *testing.T) {
	Convey("SendWechatMessage", t, func() {
		var (
			c        = context.Background()
			addr     = "127.0.0.1:8000"
			aType    = canal.TypeMap[canal.TypeReview]
			result   = canal.TypeMap[canal.ReviewSuccess]
			sender   = "fengshanshan"
			receiver = []string{"fengshanshan"}
			note     = "test"
		)
		err := svr.SendWechatMessage(c, addr, aType, result, sender, note, receiver)
		So(err, ShouldBeNil)
	})
}
