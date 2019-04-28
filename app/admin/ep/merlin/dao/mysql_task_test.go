package dao

import (
	"testing"
	"time"

	"go-common/app/admin/ep/merlin/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testMachine = model.Machine{
		EndTime: time.Now().Add((-1) * time.Hour),
		ID:      123123,
	}
)

func Test_Task(t *testing.T) {
	Convey("test InsertDeleteMachinesTasks", t, func() {
		err := d.InsertDeleteMachinesTasks([]*model.Machine{&testMachine})
		So(err, ShouldBeNil)
	})

	Convey("test FindDeleteMachineTasks", t, func() {
		tasks, err := d.FindDeleteMachineTasks()
		So(err, ShouldBeNil)
		So(len(tasks), ShouldBeGreaterThan, 0)
	})

	Convey("test UpdateTaskStatusByMachines", t, func() {
		err := d.UpdateTaskStatusByMachines([]int64{testMachine.ID}, 2)
		So(err, ShouldBeNil)
	})

}
