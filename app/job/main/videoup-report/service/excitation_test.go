package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
	"testing"
)

func TestService_HdlExcitation(t *testing.T) {
	Convey("HdlExcitation", t, func() {
		n := &archive.Archive{
			ID:  17191032,
			Mid: 27515256,
		}
		err := s.hdlExcitation(n, nil)
		So(err, ShouldBeNil)
	})
}
