package dao

import (
	"strconv"
	"testing"
	"time"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tmpIssueNoStr3 = strconv.FormatInt(time.Now().Unix(), 10)

	scheduleTask = &model.ScheduleTask{
		Name:   tmpIssueNoStr3,
		Status: model.TaskStatusRunning,
	}
)

func Test_Schedule_task(t *testing.T) {
	Convey("test insert schedule task", t, func() {
		err := d.InsertScheduleTask(scheduleTask)
		So(err, ShouldBeNil)
	})

	Convey("test update schedule task", t, func() {
		scheduleTask.Status = model.TaskStatusDone
		err := d.UpdateScheduleTask(scheduleTask)
		So(err, ShouldBeNil)
	})
}
