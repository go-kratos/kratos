package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCidInfo(t *testing.T) {
	Convey("test cid info", t, func() {
		cidInfo, err := testDao.CidInfo(context.TODO(), 10109082)
		So(err, ShouldBeNil)
		So(cidInfo, ShouldNotBeNil)
	})
}
