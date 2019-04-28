package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestBase(t *testing.T) {
	convey.Convey("Base", t, func() {
		base, err := d.Base(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(base, convey.ShouldNotBeNil)
	})
}

func TestExp(t *testing.T) {
	convey.Convey("Exp", t, func() {
		exp, err := d.Exp(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(exp, convey.ShouldNotBeNil)
	})
}

func TestMoral(t *testing.T) {
	convey.Convey("Moral", t, func() {
		moral, err := d.Moral(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(moral, convey.ShouldNotBeNil)
	})
}

func TestDao_BatchUserAddit(t *testing.T) {
	convey.Convey("BatchUserAddit", t, func() {
		userAddits, err := d.BatchUserAddit(context.Background(), []int64{5, 13, 3})
		convey.So(err, convey.ShouldBeNil)
		convey.So(userAddits, convey.ShouldNotBeNil)
	})
}

func TestDao_UpName(t *testing.T) {
	convey.Convey("UpName", t, func() {
		name := fmt.Sprintf("100_Bili_%v", time.Now().Unix())
		err := d.UpName(context.Background(), 100, name)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDao_UpSign(t *testing.T) {
	convey.Convey("UpSign", t, func() {
		sign := fmt.Sprintf("签名%v", time.Now().Unix())
		err := d.UpSign(context.Background(), 100, sign)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDao_UpFace(t *testing.T) {
	convey.Convey("UpFace", t, func() {
		face := fmt.Sprintf("testFace%v.jpg", time.Now().Unix())
		err := d.UpFace(context.Background(), 100, face)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDao_PubExpMsg(t *testing.T) {
	convey.Convey("PubExpMsg", t, func() {
		msg := &model.AddExpMsg{
			Mid:   1,
			IP:    "127.0.0.1",
			Ts:    time.Now().Unix(),
			Event: "test",
		}
		err := d.PubExpMsg(context.Background(), msg)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDao_UserAddit(t *testing.T) {
	convey.Convey("UserAddit", t, func() {
		userAddit, err := d.UserAddit(context.Background(), 5)
		convey.So(err, convey.ShouldBeNil)
		convey.So(userAddit, convey.ShouldNotBeNil)
	})
}

func TestDao_UpAdditRemark(t *testing.T) {
	convey.Convey("UpAdditRemark", t, func() {
		remark := fmt.Sprintf("remark%v", time.Now().Unix())
		err := d.UpAdditRemark(context.Background(), 27515431, remark)
		convey.So(err, convey.ShouldBeNil)
	})
}
