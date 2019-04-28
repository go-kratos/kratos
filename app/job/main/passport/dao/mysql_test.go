package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/passport/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_AddLoginLog(t *testing.T) {
	vs := make([]*model.LoginLog, 0)
	v := &model.LoginLog{
		Mid:       10,
		LoginIP:   InetAtoN("127.0.0.1"),
		Timestamp: time.Now().Unix(),
		Type:      1,
		Server:    "server",
	}
	for i := 0; i < 10; i++ {
		vs = append(vs, v)
	}
	if err := d.AddLoginLog(vs); err != nil {
		t.Errorf("dao.AddLoginLog(%v) error(%v)", vs, err)
		t.FailNow()
	}
}

func TestDao_QueryTelBindLog(t *testing.T) {
	convey.Convey("", t, func() {
		res, err := d.QueryTelBindLog(1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(res.ID, convey.ShouldEqual, 1)
	})
}

func TestDao_QueryEmailBindLog(t *testing.T) {
	convey.Convey("", t, func() {
		res, err := d.QueryEmailBindLog(1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(res.ID, convey.ShouldEqual, 1)
	})
}

func TestDao_BatchGetPwdLog(t *testing.T) {
	convey.Convey("BatchGetPwdLog", t, func() {
		var (
			c = context.Background()
		)
		res, err := d.BatchGetPwdLog(c, 100000000)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestDao_GetPwdLog(t *testing.T) {
	convey.Convey("GetPwdLog", t, func() {
		var (
			c = context.Background()
		)
		res, err := d.GetPwdLog(c, 1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
