package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SubSort(t *testing.T) {
	Convey("AddCustomSubTags", t, WithService(func(s *Service) {
		s.AddCustomSubTags(context.Background(), mid, 3, tids, ip)
	}))
	Convey("CustomSubTags", t, WithService(func(s *Service) {
		s.CustomSubTags(context.Background(), mid, 3, pn, ps, order)
	}))
}
