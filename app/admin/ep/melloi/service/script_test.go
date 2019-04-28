package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/ep/melloi/model"
)

var (
	script = model.Script{
		ID:         1,
		Type:       1,
		TestName:   "testName",
		ThreadsSum: 10,
		LoadTime:   10,
		ReadyTime:  100,
		ProcType:   "https",
		URL:        "live.bilibili.com/as/xxx",
		Domain:     "live.bilibili.com",
		Port:       "80",
		Path:       "/x/v2/search?actionKey=appkey&appkey=27eb53fc9058f8c3&build=6790&device=phone&duration=0&from_source=app_search&highlight=1&keyword=${test}&access_key=${access_key}",
		Method:     "Get",
		UpdateBy:   "hujianping",
		JmeterLog:  "/data/jmeter-log/test/ep/melloi/",
		ResJtl:     "/data/jmeter-log/test/ep/melloi/",
	}
	fileWrite = false
)

func Test_Script(t *testing.T) {
	Convey("add script", t, func() {
		_, _, err := s.AddScript(&script, fileWrite)
		So(err, ShouldBeNil)
	})
	Convey("query script", t, func() {
		_, err := s.QueryScripts(&script, 1, 1)
		So(err, ShouldBeNil)
	})
	Convey("count query script", t, func() {
		cs := s.CountQueryScripts(&script)
		So(cs, ShouldNotBeNil)
	})
	Convey("delete script", t, func() {
		err := s.DelScript(script.ID)
		So(err, ShouldBeNil)
	})
	Convey("update script", t, func() {
		_, err := s.UpdateScript(&script)
		So(err, ShouldBeNil)
	})
	Convey("add jmeter sample", t, func() {
		_, err := s.AddJmeterSample(&script)
		So(err, ShouldBeNil)
	})
}
