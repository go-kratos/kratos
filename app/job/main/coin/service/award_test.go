package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAward(t *testing.T) {
	Convey("award", t, func() {
		err := s.award(context.TODO(), 1, time.Now().Unix(), "127.0.0.1")
		So(err, ShouldBeNil)
	})
}
