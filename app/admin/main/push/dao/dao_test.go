package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/push/conf"
	"go-common/app/admin/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/push-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		f(d)
	}
}

func Test_Dao(t *testing.T) {
	Convey("dao test", t, WithDao(func(d *Dao) {
		d.Ping(context.TODO())
	}))
}

func Test_AddDPCondition(t *testing.T) {
	Convey("AddDPCondition", t, WithDao(func(d *Dao) {
		cond := &model.DPCondition{
			Task:      123,
			Job:       "456",
			Condition: "cond",
			SQL:       "sql",
			Status:    2,
			StatusURL: "status url",
			File:      "file",
		}
		id, err := d.AddDPCondition(context.Background(), cond)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThan, 0)
	}))
}

func Test_DPCondition(t *testing.T) {
	Convey("DPContion", t, WithDao(func(d *Dao) {
		res, err := d.DPCondition(context.Background(), "456")
		So(err, ShouldBeNil)
		t.Logf("res(%+v)", res)
	}))
}

func Test_AddTask(t *testing.T) {
	Convey("add task", t, WithDao(func(d *Dao) {
		t := &model.Task{Job: "123", AppID: 2}
		_, err := d.AddTask(context.Background(), t)
		So(err, ShouldBeNil)
	}))
}

func Test_TaskInfo(t *testing.T) {
	Convey("task info", t, WithDao(func(d *Dao) {
		task, err := d.TaskInfo(context.Background(), 117)
		So(err, ShouldBeNil)
		t.Logf("task(%+v)", task)
	}))
}

func Test_Partitions(t *testing.T) {
	Convey("partitions", t, WithDao(func(d *Dao) {
		res, err := d.Partitions(context.Background())
		So(err, ShouldBeNil)
		t.Logf("partitions(%v)", res)
	}))
}
