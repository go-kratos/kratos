package dao

import (
	"context"
	"go-common/app/interface/live/push-live/model"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_task(t *testing.T) {
	initd()
	Convey("Parse Json To Struct", t, func() {
		task := &model.ApPushTask{
			Type:       rand.Intn(9999) + 1,
			TargetID:   rand.Int63n(9999) + 1,
			AlertTitle: "title",
			AlertBody:  "body",
			MidSource:  rand.Intn(15),
			LinkType:   rand.Intn(10),
			LinkValue:  "link_value",
			Total:      rand.Intn(9999),
		}
		affected, err := d.CreateNewTask(context.TODO(), task)
		t.Logf("the result included(%v) err(%v)", affected, err)

		So(err, ShouldEqual, nil)
	})
}
