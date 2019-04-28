package service

import (
	"context"
	"math/rand"
	"testing"

	"go-common/app/service/openplatform/abtest/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testVer int64
	testID  int
	testAb  *model.AB
)

func TestSyncVersionID(t *testing.T) {
	Convey("TestSyncVersionID: ", t, func() {
		svr.SyncVersionID(context.TODO())
	})
}

func TestVersionID(t *testing.T) {
	var err error
	Convey("TestVersionID: ", t, func() {
		testVer, err = svr.VersionID(nil, 1)
		So(err, ShouldBeNil)
	})
}

func TestVersion(t *testing.T) {
	Convey("TestVersion: ", t, func() {
		over := &model.Version{}
		_, err := svr.Version(context.TODO(), 1, "key", over, "1")
		So(err, ShouldBeNil)
	})
}

func TestListAb(t *testing.T) {
	Convey("TestListAb: ", t, func() {
		_, _, err := svr.ListAb(context.TODO(), 1, 10, "0,1,2", 1)
		So(err, ShouldBeNil)
	})
}

func TestAddAb(t *testing.T) {
	var nab = &model.AB{
		Name:   "test",
		Desc:   "desc",
		Stra:   model.Stra{Precision: 100, Ratio: []int{80, 20}},
		Seed:   rand.Intn(1000000000),
		Result: 0,
		Group:  1,
		Author: "test",
	}
	Convey("TestAddAb: ", t, func() {
		res, err := svr.AddAb(context.TODO(), nab)
		So(err, ShouldBeNil)
		testID = int(res["newid"].(int64))
	})
}

func TestAb(t *testing.T) {
	var err error
	Convey("TestAb: ", t, func() {
		testAb, err = svr.Ab(context.TODO(), testID, 1)
		So(err, ShouldBeNil)
	})
}

func TestUpdateAb(t *testing.T) {
	Convey("TestUpdateAb: ", t, func() {
		testAb.Desc = "testUpdate"
		res, err := svr.UpdateAb(context.TODO(), testID, testAb)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, true)
	})
}

func TestUpdateStatus(t *testing.T) {
	Convey("TestUpdateStatus: ", t, func() {
		res, err := svr.UpdateStatus(context.TODO(), testID, 1, "update", 1)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, true)
		svr.UpdateStatus(context.TODO(), testID, 0, "update", 1)
	})
}

func TestDeleteAb(t *testing.T) {
	Convey("TestDeleteAb: ", t, func() {
		res, err := svr.DeleteAb(context.TODO(), testID)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, true)
		res, _ = svr.DeleteAb(context.TODO(), 1111)
		So(res, ShouldEqual, false)
	})
}
