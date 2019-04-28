package tag

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_TagIncome(t *testing.T) {
	Convey("TagIncome", t, WithService(func(s *Service) {
		err := s.TagIncome(context.Background(), "2018-01-01", _video)
		So(err, ShouldBeNil)
	}))
}
