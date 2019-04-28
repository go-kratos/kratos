package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceExpiredEquip(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := s.expiredEquip(context.TODO())
		So(err, ShouldBeNil)
	})
}
