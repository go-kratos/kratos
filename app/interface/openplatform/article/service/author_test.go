package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpdateAuthorCache(t *testing.T) {
	mid := int64(99999999)
	Convey("add author", t, WithCleanCache(func() {
		//load data
		err := s.UpdateAuthorCache(context.TODO(), mid)
		So(err, ShouldBeNil)
		time.Sleep(time.Millisecond * 100)
		ok, _, _ := s.IsAuthor(context.TODO(), mid)
		So(ok, ShouldBeTrue)
	}))
}
