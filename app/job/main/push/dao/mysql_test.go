package dao

import (
	"context"
	"testing"
	"time"

	pushmdl "go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DelCallbacks(t *testing.T) {
	Convey("del callbacks", t, func() {
		loc, _ := time.LoadLocation("Local")
		tm := time.Date(2018, 1, 11, 18, 27, 03, 0, loc)
		rows, err := d.DelCallbacks(context.TODO(), tm, 1000)
		So(err, ShouldBeNil)
		t.Logf("del callback rows:%d", rows)
	})
}

func Test_DelTasks(t *testing.T) {
	Convey("del tasks", t, func() {
		loc, _ := time.LoadLocation("Local")
		tm := time.Date(2018, 4, 2, 16, 00, 00, 0, loc)
		rows, err := d.DelTasks(context.TODO(), tm, 1000)
		So(err, ShouldBeNil)
		t.Logf("del task rows:%d", rows)
	})
}

func Test_ReportLastID(t *testing.T) {
	Convey("get report latest id", t, func() {
		id, err := d.ReportLastID(context.TODO())
		So(err, ShouldBeNil)
		t.Logf("latest report id(%d)", id)
	})
}

func Test_TxTaskByStatus(t *testing.T) {
	Convey("tx task by status", t, func() {
		tx, _ := d.BeginTx(context.Background())
		_, err := d.TxTaskByStatus(tx, pushmdl.TaskStatusPrepared)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_AddTask(t *testing.T) {
	Convey("tx add task", t, func() {
		err := d.AddTask(context.Background(), &pushmdl.Task{})
		So(err, ShouldBeNil)
	})
}
