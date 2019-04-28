package redis

import (
	"context"
	"testing"

	tmod "go-common/app/job/main/videoup-report/model/task"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SetVideoJam(t *testing.T) {
	Convey("SetVideoJam", t, func() {
		_ = d.SetVideoJam(context.TODO(), 1)
	})
}

func Test_TrackAddRedis(t *testing.T) {
	Convey("TrackAddRedis", t, func() {
		err := d.TrackAddRedis(context.TODO(), "", "")
		So(err, ShouldNotBeNil)
	})
}

func Test_SetWeight(t *testing.T) {
	Convey("SetWeight", t, func() {
		err := d.SetWeight(context.TODO(), map[int64]*tmod.WeightParams{
			1: &tmod.WeightParams{
				TaskID: 1,
				Weight: 1,
			},
			2: &tmod.WeightParams{
				TaskID: 2,
				Weight: 2,
			},
		})
		So(err, ShouldBeNil)
	})
}

func Test_GetWeight(t *testing.T) {
	Convey("GetWeight", t, func() {
		mcases, err := d.GetWeight(context.TODO(), []int64{1, 2})
		So(err, ShouldBeNil)
		So(mcases, ShouldNotBeNil)
	})
}
