package like

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_HotDot(t *testing.T) {
	Convey("test hot dot", t, WithService(func(s *Service) {
		mid := int64(908085)
		data, err := s.RedDot(context.Background(), mid)
		So(err, ShouldBeNil)
		Println(data)
	}))
}
