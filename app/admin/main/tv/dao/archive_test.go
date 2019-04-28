package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_ArchiveAdd(t *testing.T) {
	Convey("TestDao_ArchiveAdd", t, WithDao(func(d *Dao) {
		err := d.NeedImport(123)
		So(err, ShouldBeNil)
	}))
}
