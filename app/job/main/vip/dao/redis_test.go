package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx       = context.TODO()
	mid int64 = 7593623
)

func Test_FrozenQueue(t *testing.T) {
	Convey("Test_FrozenQueue", t, func() {
		var (
			err   error
			score = time.Now().Add(-50 * time.Second).Unix()
		)
		err = d.Enqueue(ctx, mid, score)
		So(err, ShouldBeNil)
		res, err2 := d.Dequeue(ctx)
		So(err2, ShouldBeNil)
		So(res[0], ShouldEqual, mid)
		err3 := d.RemQueue(ctx, mid)
		So(err3, ShouldBeNil)
	})
}

func TestDao_Enqueue(t *testing.T) {
	Convey("test enqueue", t, func() {
		duration := time.Duration(d.c.Property.FrozenDate)
		t.Logf("duration %+v", duration)
		err := d.Enqueue(ctx, mid, time.Now().Add(duration).Unix())
		So(err, ShouldBeNil)
	})
}

func TestDao_Dequeue(t *testing.T) {
	Convey("dequeue", t, func() {
		res, err := d.Dequeue(ctx)
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}

func Test_AddLogginIP(t *testing.T) {
	Convey("Test_FrozenQueue", t, func() {
		err := d.AddLogginIP(ctx, mid, 234231)
		So(err, ShouldBeNil)
	})
}
