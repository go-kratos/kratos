package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_POI(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("Test_POI", t, WithDao(func(d *Dao) {
		_, err = d.POI(c, 23333)
		So(err, ShouldBeNil)
	}))
}

func Test_Vote(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("Test_Vote", t, WithDao(func(d *Dao) {
		_, err = d.Vote(c, 23333)
		So(err, ShouldBeNil)
	}))
}
