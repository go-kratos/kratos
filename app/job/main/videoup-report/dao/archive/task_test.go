package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_TaskByUntreated(t *testing.T) {
	Convey("TaskByUntreated", t, func() {
		configs, err := d.TaskByUntreated(context.TODO())
		So(err, ShouldBeNil)
		Println(configs)
	})
}

func Test_TaskTookByHalfHour(t *testing.T) {
	Convey("TaskTookByHalfHour", t, func() {
		configs, err := d.TaskTookByHalfHour(context.TODO())
		So(err, ShouldBeNil)
		Println(configs)
	})
}
