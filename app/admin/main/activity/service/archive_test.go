package service

import (
	"context"
	"fmt"
	"testing"

	lmdl "go-common/app/admin/main/activity/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Archives(t *testing.T) {
	Convey("service test", t, WithService(func(s *Service) {
		res, err := s.Archives(context.Background(), &lmdl.ArchiveParam{Aids: []int64{10110582, 10110581}})
		So(err, ShouldBeNil)
		for _, v := range res {
			fmt.Printf("%+v", v)
		}
	}))
}
