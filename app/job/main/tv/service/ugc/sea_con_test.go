package ugc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SearUgcCon(t *testing.T) {
	Convey("search season content count", t, WithService(func(s *Service) {
		s.seaUgcCont()
	}))
}
