package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubject(t *testing.T) {
	Convey("", t, func() {
		sub, err := svr.subject(context.TODO(), 1, 1221)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeNil)
		t.Logf("subject:%+v", sub)
	})
}
