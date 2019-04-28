package service

import (
	"context"
	"testing"

	"go-common/app/interface/openplatform/monitor-end/model"
	"go-common/app/interface/openplatform/monitor-end/model/monitor"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testGroup = &model.Group{
		Name:      "test",
		Receivers: "lalalal",
		Interval:  30,
	}
	tmpGroup, tmpGroup2 *model.Group
)

// TestAddGroup .
func TestAddGroup(t *testing.T) {
	c := context.Background()
	Convey("add group ", t, func() {
		id, err := svr.AddGroup(c, testGroup)
		So(err, ShouldBeNil)
		tmpGroup, err = svr.dao.Group(c, id)
		So(err, ShouldBeNil)
		testGroup.Name = "hahaha"
		id, err = svr.AddGroup(c, testGroup)
		So(err, ShouldBeNil)
		tmpGroup2, err = svr.dao.Group(c, id)
		So(err, ShouldBeNil)
	})
	Convey("add group again with duplicate name", t, func() {
		_, err := svr.AddGroup(c, testGroup)
		So(err, ShouldNotBeNil)
	})
}

// TestUpdateGroup .
func TestUpdateGroup(t *testing.T) {
	c := context.Background()
	Convey("update group ", t, func() {
		tmpGroup.Name = "huhuhu"
		err := svr.UpdateGroup(c, tmpGroup)
		So(err, ShouldBeNil)
	})
	Convey("update group again with duplicate name", t, func() {
		tmpGroup2.Name = "huhuhu"
		err := svr.UpdateGroup(c, tmpGroup2)
		So(err, ShouldNotBeNil)
	})
}

// TestGroupList .
func TestGroupList(t *testing.T) {
	c := context.Background()
	Convey("group list", t, func() {
		_, err := svr.GroupList(c, nil)
		So(err, ShouldBeNil)
	})
}

// TestDeleteGroup .
func TestDeleteGroup(t *testing.T) {
	c := context.Background()
	Convey("delete group", t, func() {
		err := svr.DeleteGroup(c, tmpGroup.ID)
		So(err, ShouldBeNil)
		err = svr.DeleteGroup(c, tmpGroup2.ID)
		So(err, ShouldBeNil)
	})
}

var (
	testTarget = &model.Target{
		SubEvent:  "mall.bilibili.com/home",
		Event:     "test",
		Product:   "test",
		Source:    "test",
		GroupIDs:  "1,2,3",
		Threshold: 4,
		Duration:  5,
	}
	tmpTarget, tmpTarget2 *model.Target
)

// TestAddTarget .
func TestAddTarget(t *testing.T) {
	c := context.Background()
	Convey("add target ", t, func() {
		id, err := svr.AddTarget(c, testTarget)
		So(err, ShouldBeNil)
		tmpTarget, err = svr.dao.Target(c, id)
		So(err, ShouldBeNil)
		testTarget.Product = "tttt"
		id, err = svr.AddTarget(c, testTarget)
		So(err, ShouldBeNil)
		tmpTarget2, err = svr.dao.Target(c, id)
		So(err, ShouldBeNil)
	})
	Convey("add Target again with duplicate name", t, func() {
		_, err := svr.AddTarget(c, testTarget)
		So(err, ShouldNotBeNil)
	})
}

// TestUpdateTarget .
func TestUpdateTarget(t *testing.T) {
	c := context.Background()
	Convey("update Target ", t, func() {
		tmpTarget.Product = "xxxx"
		err := svr.UpdateTarget(c, tmpTarget)
		So(err, ShouldBeNil)
	})
	Convey("update Target again with duplicate name", t, func() {
		tmpTarget2.Product = "xxxx"
		err := svr.UpdateTarget(c, tmpTarget2)
		So(err, ShouldNotBeNil)
	})
}

// TestTargetList .
func TestTargetList(t *testing.T) {
	c := context.Background()
	Convey("Target list", t, func() {
		t := &model.Target{}
		_, err := svr.TargetList(c, t, 1, 2, "mtime,0")
		So(err, ShouldBeNil)
	})
}

// TestTargetSync .
func TestTargetSync(t *testing.T) {
	c := context.Background()
	Convey("sync Target", t, func() {
		err := svr.TargetSync(c, tmpTarget.ID, 1)
		So(err, ShouldBeNil)
		err = svr.TargetSync(c, tmpTarget2.ID, 1)
		So(err, ShouldBeNil)
	})
}

// TestDeleteTarget .
func TestDeleteTarget(t *testing.T) {
	c := context.Background()
	Convey("delete Target", t, func() {
		err := svr.DeleteTarget(c, tmpTarget.ID)
		So(err, ShouldBeNil)
		err = svr.DeleteTarget(c, tmpTarget2.ID)
		So(err, ShouldBeNil)
	})
}

var (
	testProduct = &model.Product{
		Name:     "test",
		GroupIDs: "1,2",
		State:    1,
	}
	tmpProduct, tmpProduct2 *model.Product
)

// TestAddProduct .
func TestAddProduct(t *testing.T) {
	c := context.Background()
	Convey("add product ", t, func() {
		id, err := svr.AddProduct(c, testProduct)
		So(err, ShouldBeNil)
		tmpProduct, err = svr.dao.Product(c, id)
		So(err, ShouldBeNil)
		testProduct.Name = "hahaha"
		id, err = svr.AddProduct(c, testProduct)
		So(err, ShouldBeNil)
		tmpProduct2, err = svr.dao.Product(c, id)
		So(err, ShouldBeNil)
	})
	Convey("add product again with duplicate name", t, func() {
		_, err := svr.AddProduct(c, testProduct)
		So(err, ShouldNotBeNil)
	})
}

// TestUpdateProduct .
func TestUpdateProduct(t *testing.T) {
	c := context.Background()
	Convey("update product ", t, func() {
		tmpProduct.Name = "huhuhu"
		err := svr.UpdateProduct(c, tmpProduct)
		So(err, ShouldBeNil)
	})
	Convey("update Product again with duplicate name", t, func() {
		tmpProduct2.Name = "huhuhu"
		err := svr.UpdateProduct(c, tmpProduct2)
		So(err, ShouldNotBeNil)
	})
}

// TestAllProducts .
func TestAllProducts(t *testing.T) {
	c := context.Background()
	Convey("all products", t, func() {
		_, err := svr.AllProducts(c)
		So(err, ShouldBeNil)
	})
}

// TestDeleteProduct .
func TestDeleteProduct(t *testing.T) {
	c := context.Background()
	Convey("delete product", t, func() {
		err := svr.DeleteProduct(c, tmpProduct.ID)
		So(err, ShouldBeNil)
		err = svr.DeleteProduct(c, tmpProduct2.ID)
		So(err, ShouldBeNil)
	})
}

var (
	testCollect = &monitor.Log{
		SubEvent: "aaa",
		Event:    "bbb",
		Product:  "bbb",
		ExtJSON:  "xxx",
		HTTPCode: "404",
	}
)

// TestCollect .
func TestCollect(t *testing.T) {
	c := context.Background()
	Convey("test collect", t, func() {
		svr.Collect(c, testCollect)
	})
}
