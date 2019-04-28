package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_ZAddPush(t *testing.T) {
	Convey("TestDao_ZAddPush", t, WithDao(func(d *Dao) {
		err := d.ZAddPush(context.Background(), 381)
		So(err, ShouldBeNil)
	}))
}
