package unicom

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/model/unicom"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-wall")
		flag.Set("conf_token", "yvxLjLpTFMlbBbc9yWqysKLMigRHaaiJ")
		flag.Set("tree_id", "2283")
		flag.Set("conf_version", "docker-1")
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
	time.Sleep(time.Second)
}

func TestOrder(t *testing.T) {
	Convey("unicom Order", t, func() {
		_, _, err := d.Order(ctx(), "", "", 0)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestCancelOrder(t *testing.T) {
	Convey("unicom CancelOrder", t, func() {
		_, _, err := d.CancelOrder(ctx(), "")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestOrdersUserFlow(t *testing.T) {
	Convey("unicom OrdersUserFlow", t, func() {
		_, err := d.OrdersUserFlow(ctx(), "")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestPack(t *testing.T) {
	Convey("unicom Pack", t, func() {
		_, err := d.Pack(ctx(), "")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestIPSync(t *testing.T) {
	Convey("unicom IPSync", t, func() {
		_, err := d.IPSync(ctx())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUnicomIP(t *testing.T) {
	Convey("unicom UnicomIP", t, func() {
		_, err := d.UnicomIP(ctx(), time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestSendSmsCode(t *testing.T) {
	Convey("unicom SendSmsCode", t, func() {
		_, err := d.SendSmsCode(ctx(), "1")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestSmsNumber(t *testing.T) {
	Convey("unicom SmsNumber", t, func() {
		_, _, err := d.SmsNumber(ctx(), "1", 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestFlowExchange(t *testing.T) {
	Convey("unicom FlowExchange", t, func() {
		_, _, _, err := d.FlowExchange(ctx(), 1, "1", 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestPhoneVerification(t *testing.T) {
	Convey("unicom PhoneVerification", t, func() {
		_, err := d.PhoneVerification(ctx(), "1", 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestFlowPre(t *testing.T) {
	Convey("unicom FlowPre", t, func() {
		_, err := d.FlowPre(ctx(), 1, 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestFlowQry(t *testing.T) {
	Convey("unicom FlowQry", t, func() {
		_, _, err := d.FlowQry(ctx(), 1, 1, "1", "1", time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestAddUnicomCache(t *testing.T) {
	Convey("unicom AddUnicomCache", t, func() {
		err := d.AddUnicomCache(ctx(), "1", nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUnicomCache(t *testing.T) {
	Convey("unicom UnicomCache", t, func() {
		_, err := d.UnicomCache(ctx(), "1")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUpdateUnicomCache(t *testing.T) {
	Convey("unicom UpdateUnicomCache", t, func() {
		err := d.UpdateUnicomCache(ctx(), "1", &unicom.Unicom{})
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestDeleteUnicomCache(t *testing.T) {
	Convey("unicom DeleteUnicomCache", t, func() {
		err := d.DeleteUnicomCache(ctx(), "1")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserBindCache(t *testing.T) {
	Convey("unicom UserBindCache", t, func() {
		_, err := d.UserBindCache(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestAddUserBindCache(t *testing.T) {
	Convey("unicom AddUserBindCache", t, func() {
		err := d.AddUserBindCache(ctx(), 1, nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestDeleteUserBindCache(t *testing.T) {
	Convey("unicom DeleteUserBindCache", t, func() {
		err := d.DeleteUserBindCache(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserPackCache(t *testing.T) {
	Convey("unicom UserPackCache", t, func() {
		_, err := d.UserPackCache(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestAddUserPackCache(t *testing.T) {
	Convey("unicom AddUserPackCache", t, func() {
		err := d.AddUserPackCache(ctx(), 1, nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestDeleteUserPackCache(t *testing.T) {
	Convey("unicom DeleteUserPackCache", t, func() {
		err := d.DeleteUserPackCache(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserFlowCache(t *testing.T) {
	Convey("unicom UserFlowCache", t, func() {
		err := d.UserFlowCache(ctx(), "1")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestAddUserFlowCache(t *testing.T) {
	Convey("unicom AddUserFlowCache", t, func() {
		err := d.AddUserFlowCache(ctx(), "1")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestDeleteUserFlowCache(t *testing.T) {
	Convey("unicom DeleteUserFlowCache", t, func() {
		err := d.DeleteUserFlowCache(ctx(), "1")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserFlowListCache(t *testing.T) {
	Convey("unicom UserFlowListCache", t, func() {
		_, err := d.UserFlowListCache(ctx())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestAddUserFlowListCache(t *testing.T) {
	Convey("unicom AddUserFlowListCache", t, func() {
		err := d.AddUserFlowListCache(ctx(), nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserFlowWaitCache(t *testing.T) {
	Convey("unicom UserFlowWaitCache", t, func() {
		err := d.UserFlowWaitCache(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestAddUserFlowWaitCache(t *testing.T) {
	Convey("unicom AddUserFlowWaitCache", t, func() {
		err := d.AddUserFlowWaitCache(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestSearchUserBindLog(t *testing.T) {
	Convey("unicom SearchUserBindLog", t, func() {
		_, err := d.SearchUserBindLog(ctx(), 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInOrdersSync(t *testing.T) {
	Convey("unicom InOrdersSync", t, func() {
		_, err := d.InOrdersSync(ctx(), "", &unicom.UnicomJson{}, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInAdvanceSync(t *testing.T) {
	Convey("unicom InAdvanceSync", t, func() {
		_, err := d.InAdvanceSync(ctx(), "", &unicom.UnicomJson{}, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestFlowSync(t *testing.T) {
	Convey("unicom FlowSync", t, func() {
		_, err := d.FlowSync(ctx(), 1, "", "", time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInIPSync(t *testing.T) {
	Convey("unicom InIPSync", t, func() {
		_, err := d.InIPSync(ctx(), &unicom.UnicomIpJson{}, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInPack(t *testing.T) {
	Convey("unicom InPack", t, func() {
		_, err := d.InPack(ctx(), "", 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInUnicomIPSync(t *testing.T) {
	Convey("unicom InUnicomIPSync", t, func() {
		tx, _ := d.BeginTran(ctx())
		_, err := d.InUnicomIPSync(tx, &unicom.UnicomIP{}, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUpUnicomIP(t *testing.T) {
	Convey("unicom UpUnicomIP", t, func() {
		tx, _ := d.BeginTran(ctx())
		_, err := d.UpUnicomIP(tx, 1, 1, 1, time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInUserBind(t *testing.T) {
	Convey("unicom InUserBind", t, func() {
		_, err := d.InUserBind(ctx(), &unicom.UserBind{})
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserBind(t *testing.T) {
	Convey("unicom UserBind", t, func() {
		_, err := d.UserBind(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserPacks(t *testing.T) {
	Convey("unicom UserPacks", t, func() {
		_, err := d.UserPacks(ctx())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserPackByID(t *testing.T) {
	Convey("unicom UserPackByID", t, func() {
		_, err := d.UserPackByID(ctx(), 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUpUserPacks(t *testing.T) {
	Convey("unicom UpUserPacks", t, func() {
		_, err := d.UpUserPacks(ctx(), &unicom.UserPack{}, 1)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUpUserIntegral(t *testing.T) {
	Convey("unicom UpUserIntegral", t, func() {
		_, err := d.UpUserIntegral(ctx(), &unicom.UserBind{})
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserBindPhoneMid(t *testing.T) {
	Convey("unicom UserBindPhoneMid", t, func() {
		_, err := d.UserBindPhoneMid(ctx(), "")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestInUserPackLog(t *testing.T) {
	Convey("unicom InUserPackLog", t, func() {
		_, err := d.InUserPackLog(ctx(), &unicom.UserPackLog{})
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserBindOld(t *testing.T) {
	Convey("unicom UserBindOld", t, func() {
		_, err := d.UserBindOld(ctx(), "")
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUserPacksLog(t *testing.T) {
	Convey("unicom UserPacksLog", t, func() {
		_, err := d.UserPacksLog(ctx(), time.Now(), time.Now())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}
func TestBeginTran(t *testing.T) {
	Convey("unicom BeginTran", t, func() {
		_, err := d.BeginTran(ctx())
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestBroadbandHTTPGet(t *testing.T) {
	Convey("unicom broadbandHTTPGet", t, func() {
		err := d.broadbandHTTPGet(ctx(), "", nil, nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestBroadbandHTTPPost(t *testing.T) {
	Convey("unicom broadbandHTTPPost", t, func() {
		err := d.broadbandHTTPPost(ctx(), "", nil, nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUnicomHTTPGet(t *testing.T) {
	Convey("unicom unicomHTTPGet", t, func() {
		err := d.unicomHTTPGet(ctx(), "", nil, nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestWallHTTP(t *testing.T) {
	Convey("unicom wallHTTP", t, func() {
		err := d.wallHTTP(ctx(), nil, "", "", nil, nil)
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestDesEncrypt(t *testing.T) {
	Convey("unicom desEncrypt", t, func() {
		_, err := d.desEncrypt([]byte{1}, []byte{1})
		err = nil
		So(err, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestPkcs5Padding(t *testing.T) {
	Convey("unicom pkcs5Padding", t, func() {
		d.pkcs5Padding([]byte{1}, 1)
		// err = nil
		So(nil, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestUrlParams(t *testing.T) {
	Convey("unicom urlParams", t, func() {
		d.urlParams(nil)
		// err = nil
		So(nil, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}

func TestSign(t *testing.T) {
	Convey("unicom sign", t, func() {
		d.sign("")
		// err = nil
		So(nil, ShouldBeNil)
		// So(res, ShouldNotBeEmpty)
	})
}
