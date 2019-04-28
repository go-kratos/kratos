package dao

import (
	"testing"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
	"time"
)

var (
	tapdBugPriorityConf = &model.TapdBugPriorityConf{
		ProjectTemplateID: 1000,
		Urgent:            123,
		High:              321,
		Medium:            111,
		UpdateBy:          "fengyifeng",
		StartTime:         time.Now(),
		EndTime:           time.Now().AddDate(0, 1, 0),
		Status:            model.TapdBugPriorityConfDisable,
	}

	queryTapdBugPriorityConfsRequest = &model.QueryTapdBugPriorityConfsRequest{
		Pagination: model.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
		ProjectTemplateID: 1000,
	}
)

func Test_Tapd_Bug_Priority_Conf(t *testing.T) {
	Convey("test Insert Tapd Bug Priority Conf", t, func() {
		err := d.InsertTapdBugPriorityConf(tapdBugPriorityConf)
		So(err, ShouldBeNil)
	})

	Convey("test Update Tapd Bug Priority Conf", t, func() {
		tapdBugPriorityConf.Urgent = 10010
		err := d.UpdateTapdBugPriorityConf(tapdBugPriorityConf)
		So(err, ShouldBeNil)
	})

	Convey("test Find Tapd Bug Priority Confs", t, func() {
		total, tapdBugPriorityConfs, err := d.FindTapdBugPriorityConfs(queryTapdBugPriorityConfsRequest)
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(tapdBugPriorityConfs), ShouldBeGreaterThan, 0)
	})

}
