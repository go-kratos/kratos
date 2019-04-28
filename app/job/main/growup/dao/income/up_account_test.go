package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpAccounts(t *testing.T) {
	Convey("UpAccounts", t, func() {
		_, _, err := d.UpAccounts(context.Background(), 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_UpAccount(t *testing.T) {
	Convey("UpAccount", t, func() {
		_, err := d.UpAccount(context.Background(), 123)
		So(err, ShouldBeNil)
	})
}

func Test_InsertUpAccount(t *testing.T) {
	Convey("InsertUpAccount", t, func() {
		_, err := d.InsertUpAccount(context.Background(), "(123,1,100,300,2018-05,22)")
		So(err, ShouldBeNil)
	})
}

func Test_UpdateUpAccount(t *testing.T) {
	Convey("UpdateUpAccount", t, func() {
		_, err := d.UpdateUpAccount(context.Background(), int64(100), int64(22), int64(123), int64(100))
		So(err, ShouldBeNil)
	})
}
