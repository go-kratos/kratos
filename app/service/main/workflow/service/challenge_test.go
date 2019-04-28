package service

import (
	"context"
	"testing"

	"go-common/app/service/main/workflow/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestAddChallenge(t *testing.T) {
	convey.Convey("Test AddChallenge", t, func() {
		_, err := s.AddChallenge(context.Background(), &model.ChallengeParam{})
		convey.ShouldBeNil(err)
	})
}

func TestCancelChallenge(t *testing.T) {
	convey.Convey("Test CancelChallenge", t, func() {
		err := s.CloseChallenge(context.Background(), 1, 1, 1, 1, "abc")
		convey.ShouldBeNil(err)
	})
}
