package service

import (
	"context"
	"testing"

	"go-common/app/service/main/usersuit/model"

	. "github.com/smartystreets/goconvey/convey"
)

// TestServiceGroupPendantMid .
func TestServiceGroupPendantMid(t *testing.T) {
	arg := &model.ArgGPMID{
		GID: 30,
		MID: 1,
	}
	Convey("need return something", t, func() {
		res, err := s.GroupPendantMid(context.Background(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
