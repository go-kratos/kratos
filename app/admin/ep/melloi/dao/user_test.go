package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_User(t *testing.T) {
	Convey("test QueryUser", t, func() {
		_, err := d.QueryUserByUserName("hujianping")
		So(err, ShouldBeNil)
	})
}
