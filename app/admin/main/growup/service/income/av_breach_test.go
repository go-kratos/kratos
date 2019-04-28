package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ArchiveBreach(t *testing.T) {
	Convey("ArchiveBreach", t, WithService(func(s *Service) {
		avID := []int64{}
		mid := int64(1)
		reason := "test"
		operator := "szy"
		err := s.ArchiveBreach(context.Background(), 0, avID, mid, reason, operator)
		So(err, ShouldBeNil)
	}))
}
