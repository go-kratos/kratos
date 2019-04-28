package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AlbumIndex(t *testing.T) {
	Convey("test album index", t, WithService(func(s *Service) {
		mid := int64(883968)
		ps := 10
		data, err := s.AlbumIndex(context.Background(), mid, ps)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
