package chuanglan

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/sms/conf"
	"go-common/app/service/main/sms/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	cl      *Client
	sendLog = &model.ModelSend{
		Content: "卍 您的账号正在哔哩哔哩2018动画角色人气大赏活动中进行领票操作，验证码为5467当日有效",
	}
)

func init() {
	dir, _ := filepath.Abs("../../cmd/sms-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	cl = NewClient(conf.Conf)
}

func TestGetPid(t *testing.T) {
	Convey("test chuanglan get pid", t, func() {
		pID := cl.GetPid()
		So(pID, ShouldEqual, model.ProviderChuangLan)
	})
}

func TestSendSms(t *testing.T) {
	Convey("test chuanglan send sms", t, func() {
		sendLog.Mobile = ""
		msgid, err := cl.SendSms(context.TODO(), sendLog)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}

func TestSendActSms(t *testing.T) {
	Convey("test ChuangLan SendActSms", t, func() {
		sendLog.Mobile = ""
		msgid, err := cl.SendActSms(context.TODO(), sendLog)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}

func TestSendBatchActSms(t *testing.T) {
	Convey("test ChuangLan sendBatchActSms", t, func() {
		sendLog.Mobile = "" // 187****3870,189****1728
		msgid, err := cl.SendBatchActSms(context.TODO(), sendLog)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}

func TestSendInternationalSms(t *testing.T) {
	Convey("test ChuangLan SendInternationalSms", t, func() {
		sendLog.Mobile = "" // 5344295506
		sendLog.Country = "1"
		msgid, err := cl.SendInternationalSms(context.TODO(), sendLog)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}
