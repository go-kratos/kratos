package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpBusinessExtra(t *testing.T) {
	convey.Convey("Test UpBusinessExtra", t, func() {
		err := s.UpBusinessExtra(context.Background(), 1, 1, 1, "", "")
		convey.ShouldBeNil(err)
	})
}

func TestBusinessExtra(t *testing.T) {
	convey.Convey("Test TestBusinessExtra", t, func() {
		_, err := s.BusinessExtra(context.Background(), 1, 1, 1)
		convey.ShouldBeNil(err)
	})
}
