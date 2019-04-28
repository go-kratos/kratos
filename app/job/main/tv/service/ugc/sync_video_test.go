package ugc

import (
	"testing"

	"go-common/app/job/main/tv/model/ugc"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SyncLic(t *testing.T) {
	Convey("TestService_SyncLic", t, WithService(func(s *Service) {
		err := s.syncLic(10099174, &ugc.SimpleArc{
			AID:      10099174,
			Title:    "test",
			Duration: 400,
			Cover:    "testCover",
		})
		So(err, ShouldBeNil)
	}))
}
