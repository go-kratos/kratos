package pgc

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SearPgcCon(t *testing.T) {
	Convey("search season content count", t, WithService(func(s *Service) {
		s.seaPgcCont()
	}))
}
