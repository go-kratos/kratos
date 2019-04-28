package pendant

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DelPKGCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelPKGCache(context.TODO(), 11)
		So(err, ShouldBeNil)
	})
}

func Test_DelEquipCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelEquipCache(context.TODO(), 11)
		So(err, ShouldBeNil)
	})
}
