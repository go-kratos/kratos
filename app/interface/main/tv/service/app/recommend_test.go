package service

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_RecomFilter(t *testing.T) {
	Convey("RecomFilter test", t, WithService(func(s *Service) {
		res, err := s.RecomFilter("23874", "3")
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}
