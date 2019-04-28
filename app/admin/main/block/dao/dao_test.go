package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/block/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	dao *Dao
	ctx = context.TODO()
)

func TestMain(m *testing.M) {
	defer os.Exit(0)
	flag.Set("conf", "../cmd/block-admin-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	dao = New()
	defer dao.Close()
	m.Run()
}

func TestDB(t *testing.T) {
	Convey("db", t, func() {
		err := dao.SendSysMsg(ctx, "123", []int64{1, 2, 3}, "test title", "test content", "")
		So(err, ShouldBeNil)
		_, err = dao.HistoryCount(ctx, 46333)
		So(err, ShouldBeNil)
		_, err = dao.History(ctx, 46333, 0, 100)
		So(err, ShouldBeNil)
		_, err = dao.User(ctx, 46333)
		So(err, ShouldBeNil)
		_, err = dao.Users(ctx, []int64{46333})
		So(err, ShouldBeNil)
		_, err = dao.UserDetails(ctx, []int64{46333})
		So(err, ShouldBeNil)
	})
}

func TestRPC(t *testing.T) {
	Convey("rpc", t, func() {
		mid := int64(46333)
		_, err := dao.SpyScore(ctx, mid)
		So(err, ShouldBeNil)
		_, err = dao.FigureRank(ctx, mid)
		So(err, ShouldBeNil)
		_, _, _, _, err = dao.AccountInfo(ctx, mid)
		So(err, ShouldBeNil)
	})
}

func TestTool(t *testing.T) {
	Convey("tool", t, func() {
		var (
			mids = []int64{1, 2, 3, 46333, 35858}
		)
		str := midsToParam(mids)
		So(str, ShouldEqual, "1,2,3,46333,35858")
	})
}

func TestHTTP(t *testing.T) {
	Convey("http", t, func() {
		var (
			mid      int64 = 46333
			telInfo  string
			mailInfo string
			err      error
		)
		telInfo, err = dao.TelInfo(ctx, mid)
		So(err, ShouldBeNil)
		mailInfo, err = dao.MailInfo(ctx, mid)
		So(err, ShouldBeNil)
		t.Logf("telinfo : %s , mailinfo : %s", telInfo, mailInfo)
	})
}
