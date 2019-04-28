package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_accNotify(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := s.accNotify(context.Background(), 1, "updateMedal")
		So(err, ShouldBeNil)
	})
}
