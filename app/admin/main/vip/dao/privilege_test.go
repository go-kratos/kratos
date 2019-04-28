package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/vip/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPrivilegeResourcesList(t *testing.T) {
	convey.Convey("PrivilegeResourcesList", t, func() {
		res, err := d.PrivilegeResourcesList(context.TODO())
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

// go test  -test.v -test.run TestDaoAddPrivilege
func TestDaoAddPrivilege(t *testing.T) {
	tx := d.BeginGormTran(context.Background())
	p := &model.Privilege{
		Name:        "超高清",
		Title:       "超高清标题",
		Explain:     "超高清描述描述超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字",
		Type:        1,
		Operator:    "admin",
		State:       0,
		Deleted:     0,
		IconURL:     "https://activity.hdslb.com/blackboard/activity9757/static/img/title_screen.c529058.jpg",
		IconGrayURL: "https://activity.hdslb.com/blackboard/activity9757/static/img/title_screen.c529058.jpg",
		Order:       1,
		LangType:    1,
	}
	convey.Convey("get max id", t, func() {
		ep := new(model.Privilege)
		db := d.vip.Table(_vipPrivileges).Order("order_num DESC").First(&ep)
		convey.So(db.Error, convey.ShouldBeNil)
		d.vip.Table(_vipPrivilegesResources).Where("pid > ?", ep.ID).Delete(model.PrivilegeResources{})
		p.ID = ep.ID + 1
	})
	convey.Convey("AddPrivilege", t, func() {
		id, err := d.AddPrivilege(tx, p)
		convey.So(err, convey.ShouldBeNil)
		convey.So(id, convey.ShouldNotBeNil)
	})
	convey.Convey("AddPrivilegeResources", t, func() {
		a, err := d.AddPrivilegeResources(tx, &model.PrivilegeResources{
			PID:  p.ID,
			Link: "web",
			Type: model.WebResources,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(a, convey.ShouldNotBeNil)
		a, err = d.AddPrivilegeResources(tx, &model.PrivilegeResources{
			PID:  p.ID,
			Link: "app",
			Type: model.AppResources,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(a, convey.ShouldNotBeNil)
	})
	convey.Convey("AddPrivilege Commit", t, func() {
		err := tx.Commit().Error
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("PrivilegeList", t, func() {
		res, err := d.PrivilegeList(context.TODO(), 1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
	convey.Convey("clean data", t, func() {
		d.vip.Delete(p)
		d.vip.Table(_vipPrivilegesResources).Where("pid >= ?", p.ID).Delete(model.PrivilegeResources{})
	})
}

func TestDaoMaxOrder(t *testing.T) {
	convey.Convey("MaxOrder", t, func() {
		order, err := d.MaxOrder(context.TODO())
		convey.So(err, convey.ShouldBeNil)
		convey.So(order, convey.ShouldNotBeNil)
	})
}

// go test  -test.v -test.run TestDaoUpdateStatePrivilege

func TestDaoUpdateStatePrivilege(t *testing.T) {
	convey.Convey("UpdateStatePrivilege", t, func() {
		a, err := d.UpdateStatePrivilege(context.TODO(), &model.Privilege{
			ID:    3,
			State: 1,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(a, convey.ShouldNotBeNil)
	})
}

// go test  -test.v -test.run TestDaoDeletePrivilege
func TestDaoDeletePrivilege(t *testing.T) {
	convey.Convey("DeletePrivilege", t, func() {
		a, err := d.DeletePrivilege(context.TODO(), 3)
		convey.So(err, convey.ShouldBeNil)
		convey.So(a, convey.ShouldNotBeNil)
	})
}

// go test  -test.v -test.run TestDaoUpdatePrivilege
func TestDaoUpdatePrivilege(t *testing.T) {
	convey.Convey("UpdatePrivilege", t, func() {
		tx := d.BeginGormTran(context.TODO())
		a, err := d.UpdatePrivilege(tx, &model.Privilege{
			ID:          3,
			Name:        "超高清",
			Title:       "超高清标题",
			Explain:     "超高清描述描述超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字超多字",
			Type:        1,
			Operator:    "admin",
			State:       0,
			Deleted:     0,
			IconURL:     "https://activity.hdslb.com/blackboard/activity9757/static/img/title_screen.c529058.jpg",
			IconGrayURL: "https://activity.hdslb.com/blackboard/activity9757/static/img/title_screen.c529058.jpg",
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(a, convey.ShouldNotBeNil)
		err = tx.Commit().Error
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDaoUpdatePrivilegeResources

func TestDaoUpdatePrivilegeResources(t *testing.T) {
	convey.Convey("UpdatePrivilegeResources", t, func() {
		tx := d.BeginGormTran(context.TODO())
		aff, err := d.UpdatePrivilegeResources(tx, &model.PrivilegeResources{
			PID:  3,
			Link: "app2",
			Type: model.AppResources,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(aff, convey.ShouldNotBeNil)
		err = tx.Commit().Error
		convey.So(err, convey.ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDaoUpdateOrder
func TestDaoUpdateOrder(t *testing.T) {
	convey.Convey("UpdateOrder", t, func() {
		a, err := d.UpdateOrder(context.TODO(), 4, 5)
		convey.So(err, convey.ShouldBeNil)
		convey.So(a, convey.ShouldNotBeNil)
	})
}
