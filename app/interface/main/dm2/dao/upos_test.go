package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpos(t *testing.T) {
	Convey("test upos", t, func() {
		saveTo, err := testDao.Upos(context.TODO(), 1111)
		So(err, ShouldBeNil)
		t.Logf("%+v", saveTo)
	})
}
