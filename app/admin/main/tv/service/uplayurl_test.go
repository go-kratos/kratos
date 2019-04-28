package service

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_UPlayurl(t *testing.T) {
	Convey("PlayURL test", t, WithService(func(s *Service) {
		url, err := s.UPlayurl(120094301)
		fmt.Println(url)
		So(err, ShouldBeNil)
	}))
}
