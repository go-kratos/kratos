package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func TestPing(t *testing.T) {
	Convey("TestPing: ", t, func() {
		err := d.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}
