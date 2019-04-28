package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Dao(t *testing.T) {
	var err error
	d := New(nil)
	Convey("set cache", t, func() {
		PromError("")
		PromInfo("")
		err = d.SendSMS("")
		So(err, ShouldBeNil)
	})
}
