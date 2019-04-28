package service

import (
	"testing"

	"go-common/app/admin/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

// go test  -test.v -test.run TestServiceAddJointly
func TestServiceAddJointly(t *testing.T) {
	Convey("TestServiceAddJointly", t, func() {
		err := s.AddJointly(c, &model.ArgAddJointly{
			Title:     "这是一条被修改的",
			Content:   "内容",
			StartTime: 1433202904,
			EndTime:   1433202905,
			Link:      "http://www.baidu.com",
			IsHot:     1,
			Operator:  "admin",
		})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceModifyJointly
func TestServiceModifyJointly(t *testing.T) {
	Convey("TestServiceModifyJointly", t, func() {
		err := s.ModifyJointly(c, &model.ArgModifyJointly{
			ID:       2,
			Title:    "无效的记录，修改",
			Content:  "修改内容",
			Link:     "http://www.baidu.com",
			IsHot:    1,
			Operator: "admin",
		})
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestServiceJointlysByState
func TestServiceJointlysByState(t *testing.T) {
	Convey("TestServiceJointlysByState", t, func() {
		res, err := s.JointlysByState(c, 2)
		t.Logf("count %+v", len(res))
		So(err, ShouldBeNil)
	})
}
