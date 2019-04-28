package rpc

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRpcFansCount(t *testing.T) {
	Convey("FansCount", t, WithMockRelation(t, func(d *Dao) {
		res, err := d.FansCount(context.TODO(), _Mids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestRpcUserInfos(t *testing.T) {
	Convey("UserInfos", t, WithDao(t, func(d *Dao) {
		res, err := d.UserInfos(context.TODO(), _Mids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestRpcProfile(t *testing.T) {
	Convey("Profile", t, WithMockAccount(t, func(d *Dao) {
		res, err := d.Profile(context.TODO(), _Mid)
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	}))
}
func TestRpcInfo3(t *testing.T) {
	Convey("Info3", t, WithMockAccount(t, func(d *Dao) {
		res, err := d.Info3(context.TODO(), _Mid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func TestRpcUpSpecial(t *testing.T) {
	Convey("UpSpecial", t, WithMockUp(t, func(d *Dao) {
		res1, err := d.UpSpecial(context.TODO(), _Mid)
		So(err, ShouldBeNil)
		So(res1, ShouldNotBeNil)
		res2, err := d.UpsSpecial(context.TODO(), _Mids)
		So(err, ShouldBeNil)
		So(res2, ShouldNotBeNil)
		res3, err := d.UpGroups(context.TODO())
		So(err, ShouldBeNil)
		So(res3, ShouldNotBeNil)
	}))
}
