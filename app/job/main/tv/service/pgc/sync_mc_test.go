package pgc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_FullRefresh(t *testing.T) {
	Convey("No redundant data", t, WithService(func(s *Service) {
		s.refreshCache()
	}))
}
