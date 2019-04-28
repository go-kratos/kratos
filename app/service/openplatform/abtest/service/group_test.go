package service

import (
	"context"
	"testing"

	"go-common/app/service/openplatform/abtest/model"

	. "github.com/smartystreets/goconvey/convey"
)

var g = model.Group{
	Name: "test",
	Desc: "test add",
}
var testDesc = "test update"

func TestAddGroup(t *testing.T) {
	Convey("TestAddGroup: ", t, func() {
		var err error
		testID, err = svr.AddGroup(context.TODO(), g)
		So(err, ShouldBeNil)
	})
}

func TestUpdateGroup(t *testing.T) {
	g.Desc = testDesc
	g.ID = testID
	Convey("TestUpdateGroup: ", t, func() {
		res, err := svr.UpdateGroup(context.TODO(), g)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 1)
	})
}

func TestListGroup(t *testing.T) {
	Convey("TestListGroup: ", t, func() {
		res, err := svr.ListGroup(context.TODO())
		So(err, ShouldBeNil)
		x := res[len(res)-1]
		So(x.Name, ShouldEqual, g.Name)
	})
}

func TestDeleteGroup(t *testing.T) {
	Convey("TestDeleteGroup: ", t, func() {
		res, err := svr.DeleteGroup(context.TODO(), testID)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 1)
	})
}
