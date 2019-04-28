package dao

import (
	"testing"
	"time"

	"go-common/app/admin/ep/merlin/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	machineLog = model.MachineLog{
		Username:      "fengyifeng",
		MachineID:     1,
		OperateType:   "DelTest",
		OperateResult: "",
		OperateTime:   time.Now(),
	}
)

func TestInsertMachineLog(t *testing.T) {
	Convey("Everything goes well when names is slice with value", t, func() {
		err := d.InsertMachineLog(&machineLog)
		So(err, ShouldBeNil)
		d.db.Where("OperateType=DelTest").Delete(machineLog)
	})
}

func TestFindMachineLogsByMachineID(t *testing.T) {
	Convey("Everything goes well when names is slice with value", t, func() {
		_, _, err := d.FindMachineLogsByMachineID(0, 1, 5)
		So(err, ShouldBeNil)
		d.db.Where("OperateType=DelTest").Delete(machineLog)
	})
}
