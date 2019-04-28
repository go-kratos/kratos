package manager

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_GetUIDByName(t *testing.T) {
	var (
		err error
	)
	Convey("GetUIDByName", t, WithDao(func(d *Dao) {
		_, err = d.GetUIDByName(context.Background(), "1111")
		So(err, ShouldBeNil)
	}))
}

func Test_GetNameByUID(t *testing.T) {
	var (
		err error
	)
	Convey("GetUIDByName", t, WithDao(func(d *Dao) {
		_, err = d.GetNameByUID(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}
