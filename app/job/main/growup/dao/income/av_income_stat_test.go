package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AvIncomeStat(t *testing.T) {
	Convey("AvIncomeStat", t, func() {
		_, _, err := d.AvIncomeStat(context.Background(), 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_InsertAvIncomeStat(t *testing.T) {
	Convey("InsertAvIncomeStat", t, func() {
		_, err := d.InsertAvIncomeStat(context.Background(), "(123,2,6,1,'2018-06-01',100)")
		So(err, ShouldBeNil)
	})
}
