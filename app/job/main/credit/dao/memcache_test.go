package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DelOpinionCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelOpinionCache(context.TODO(), 27517340)
		So(err, ShouldBeNil)
	})
}

func TestDaoDelJuryInfoCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelJuryInfoCache(context.TODO(), 27517340)
		So(err, ShouldBeNil)
	})
}
