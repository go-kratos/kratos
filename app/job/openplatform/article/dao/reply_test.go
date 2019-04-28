package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_OpenReply(t *testing.T) {
	Convey("work", t, WithDao(func(d *Dao) {
		err := d.OpenReply(context.Background(), 1, 2)
		So(err, ShouldBeNil)
	}))
}
