package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdvances(t *testing.T) {
	Convey("test adv", t, func() {
		res, _, err := svr.Advances(context.TODO(), 27515260, "all", "all", 1, 20)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
