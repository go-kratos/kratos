package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoOldFrozenChange(t *testing.T) {
	convey.Convey("OldFrozenChange", t, func() {
		err := d.OldFrozenChange(7593623, 0)
		convey.So(err, convey.ShouldBeNil)
	})
}
