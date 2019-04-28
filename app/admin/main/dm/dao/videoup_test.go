package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeInfo(t *testing.T) {
	Convey("", t, func() {
		types, err := testDao.TypeInfo(context.TODO())
		So(types, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestArchiveVideos(t *testing.T) {
	Convey("", t, func() {
		avm, err := testDao.ArchiveVideos(context.TODO(), []int64{10110320, 123, 3})
		So(avm, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
