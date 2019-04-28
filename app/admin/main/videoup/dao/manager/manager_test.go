package manager

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_User(t *testing.T) {
	convey.Convey("User", t, WithDao(func(d *Dao) {
		_, err := d.User(context.TODO(), 1)
		convey.So(err, convey.ShouldBeNil)
	}))
}
