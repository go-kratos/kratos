package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Consumers(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, err := d.Consumers(context.Background())
		So(err, ShouldBeNil)
	}))
}
func Test_IsConsumerOn(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		state := d.IsConsumerOn(context.Background(), 1)
		So(state, ShouldNotBeNil)
	}))
}
func Test_WeightConf(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		r, err := d.WeightConf(context.Background())
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
	}))
}

func Test_TaskUserCheckIn(t *testing.T) {
	Convey("TaskUserCheckIn", t, WithDao(func(d *Dao) {
		_, err := d.TaskUserCheckIn(context.Background(), 0)
		So(err, ShouldBeNil)
	}))
}

func Test_TaskUserCheckOff(t *testing.T) {
	Convey("TaskUserCheckOff", t, WithDao(func(d *Dao) {
		r, err := d.TaskUserCheckOff(context.Background(), 0)
		So(err, ShouldBeNil)
		So(r, ShouldNotBeNil)
	}))
}
