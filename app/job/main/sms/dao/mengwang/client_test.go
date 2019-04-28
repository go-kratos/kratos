package mengwang

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
	mw  *Client
	sl  = &model.ModelSend{Mobile: "", Content: "卍 测试短信，验证码为5467当日有效 https://search.bilibili.com/all?keyword=你好"} // 17621660828
	isl = &model.ModelSend{Country: "852", Mobile: "", Content: "卍 您的账号正在哔哩哔哩2017动画角色人气大赏活动中进行领票操作，验证码为5467当日有效"} // 00852 69529378 梦网的香港测试号
)

func init() {
	dir, _ := filepath.Abs("../../cmd/sms-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	mw = NewClient(conf.Conf)
}

func TestSendSms(t *testing.T) {
	Convey("mengwang send sms", t, func() {
		msgid, err := mw.SendSms(context.TODO(), sl)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}

func TestSendActSms(t *testing.T) {
	Convey("mengwang send act sms", t, func() {
		msgid, err := mw.SendActSms(context.TODO(), sl)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}

func TestSendBatchSms(t *testing.T) {
	Convey("mengwang send batch sms", t, func() {
		msl := sl
		msl.Mobile = "" // 17621660828,17621660828
		msgid, err := mw.SendActSms(context.TODO(), msl)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}

func TestSendInternationalSms(t *testing.T) {
	Convey("mengwang send international sms", t, func() {
		msgid, err := mw.SendInternationalSms(context.TODO(), isl)
		So(err, ShouldBeNil)
		t.Logf("msgid(%s)", msgid)
	})
}

func TestCallback(t *testing.T) {
	Convey("mengwang callback", t, func() {
		// callbacks, err := mw.Callback(context.Background(), conf.Conf.Pconf.MengWangSmsUser, conf.Conf.Pconf.MengWangSmsPwd, conf.Conf.Pconf.MengWangSmsCallbackURL, 5)
		// So(err, ShouldBeNil)
		// t.Logf("callbacks(%d)", len(callbacks))
		// for _, v := range callbacks {
		// 	t.Logf("callback(%+v)", v)
		// }
	})
}
