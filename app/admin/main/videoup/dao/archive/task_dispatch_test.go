package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UserUndoneSpecTask(t *testing.T) {
	Convey("UserUndoneSpecTask", t, WithDao(func(d *Dao) {
		_, err := d.UserUndoneSpecTask(context.Background(), 0)
		So(err, ShouldBeNil)
	}))
}

func Test_GetDispatchTask(t *testing.T) {
	Convey("GetDispatchTask", t, WithDao(func(d *Dao) {
		_, err := d.GetDispatchTask(context.Background(), 0)
		So(err, ShouldBeNil)
	}))
}

func Test_UpDispatchTask(t *testing.T) {
	Convey("UpDispatchTask", t, WithDao(func(d *Dao) {
		_, err := d.UpDispatchTask(context.Background(), 0, []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func Test_TaskByID(t *testing.T) {
	Convey("TaskByID", t, WithDao(func(d *Dao) {
		_, err := d.TaskByID(context.Background(), 0)
		So(err, ShouldBeNil)
	}))
}

func Test_ListByCondition(t *testing.T) {
	var c = context.Background()
	Convey("ListByCondition", t, WithDao(func(d *Dao) {
		_, err := d.ListByCondition(c, 0, 0, 0, 0, 0)
		So(err, ShouldBeNil)
	}))
}

func Test_GetWeightDB(t *testing.T) {
	Convey("GetWeightDB", t, WithDao(func(d *Dao) {
		_, err := d.GetWeightDB(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func Test_TaskDispatchByID(t *testing.T) {
	Convey("TaskDispatchByID", t, WithDao(func(d *Dao) {
		_, err := d.TaskDispatchByID(context.Background(), 0)
		So(err, ShouldBeNil)
	}))
}
