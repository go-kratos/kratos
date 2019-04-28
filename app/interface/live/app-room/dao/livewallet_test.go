package dao

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_LiveWallet(t *testing.T) {
	Convey("normal", t, func() {
		mid := int64(1)
		w, e := testDao.LiveWallet(context.Background(), mid, "android")
		So(e, ShouldBeNil)
		t.Logf("%+v", w)

	})
}
