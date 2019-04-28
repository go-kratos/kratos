package tag

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetUpTagIncomeMap(t *testing.T) {
	Convey("GetUpTagIncomeMap", t, WithService(func(s *Service) {
		_, err := s.GetUpTagIncomeMap(context.Background())
		So(err, ShouldBeNil)
	}))
}
