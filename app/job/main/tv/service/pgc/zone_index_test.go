package pgc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ZoneIdx(t *testing.T) {
	Convey("ZoneIdx", t, WithService(func(s *Service) {
		s.ZoneIdx()
	}))
}
