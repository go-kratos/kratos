package dao

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go-common/app/admin/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestDaoAddJointly
func TestDaoAddJointly(t *testing.T) {
	Convey("TestDaoAddJointly", t, func() {
		a, err := d.AddJointly(context.TODO(), &model.Jointly{
			Title:     "这条是有效的,并且不 hot",
			Content:   "这叫副标题?",
			Operator:  "admin",
			StartTime: 1533202903,
			EndTime:   1543202904,
			Link:      fmt.Sprintf("https://t.cn/%d", rand.Int63()),
			IsHot:     1,
		})
		So(err, ShouldBeNil)
		So(a, ShouldEqual, 1)
	})
}

// go test  -test.v -test.run TestDaoUpdateJointly
func TestDaoUpdateJointly(t *testing.T) {
	Convey("TestDaoUpdateJointly", t, func() {
		a, err := d.UpdateJointly(context.TODO(), &model.Jointly{
			Title:    "这条是有效的,并且no hot",
			Content:  "这叫副标题??",
			Operator: "admin2",
			Link:     fmt.Sprintf("https://t.cn/%d", rand.Int63()),
			IsHot:    0,
			ID:       1,
		})
		So(err, ShouldBeNil)
		So(a, ShouldEqual, 1)
	})
}

// go test  -test.v -test.run TestDaoJointlysByState
func TestDaoJointlysByState(t *testing.T) {
	Convey("TestDaoJointlysByState", t, func() {
		res, err := d.JointlysByState(context.TODO(), 1, time.Now().Unix())
		t.Logf("count %+v", len(res))
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestDaoDeleteJointly
func TestDaoDeleteJointly(t *testing.T) {
	Convey("TestDaoDeleteJointly", t, func() {
		res, err := d.DeleteJointly(context.TODO(), 1)
		t.Logf("count %+v", res)
		So(res, ShouldBeGreaterThanOrEqualTo, 0)
		So(err, ShouldBeNil)
	})
}
