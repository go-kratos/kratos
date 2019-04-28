package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/job/main/answer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddQue(t *testing.T) {
	Convey("TestAddQue add question data", t, func() {
		for i := 0; i < 10; i++ {
			que := &model.LabourQs{
				Question: fmt.Sprintf("测试d=====(￣▽￣*)b厉害 %d", i),
				Ans:      int8(i%2 + 1),
				AvID:     int64(i),
				Status:   int8(i%2 + 1),
				Source:   1,
				State:    model.HadCreateImg,
				ID:       int64(i),
			}
			err := d.AddQs(context.TODO(), que)
			So(err, ShouldBeNil)
		}
	})
}

func TestUpdateState(t *testing.T) {
	Convey("TestUpdateState", t, func() {
		for i := 0; i < 10; i++ {
			que := &model.LabourQs{
				Question: fmt.Sprintf("测试d=====(￣▽￣*)b厉害 %d", i),
				Ans:      int8(i%2 + 1),
				AvID:     int64(i),
				Status:   int8(i%2 + 1),
				Source:   1,
				State:    model.HadCreateImg,
				ID:       int64(i),
			}
			err := d.UpdateState(context.TODO(), que)
			So(err, ShouldBeNil)
		}
	})
}

func TestByID(t *testing.T) {
	Convey("TestByID", t, func() {
		res, err := d.ByID(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestBeFormal(t *testing.T) {
	Convey("TestBeFormal", t, func() {
		err := d.BeFormal(context.TODO(), 7593623, "127.0.0.1")
		So(err, ShouldNotBeNil)
	})
}
