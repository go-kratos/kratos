package service

import (
	"testing"

	"go-common/app/service/openplatform/ticket-item/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersionSearch(t *testing.T) {
	Convey("VersionSearch", t, func() {
		params := []*model.VersionSearchParam{
			{
				TargetItem: 1,
				PageParam: model.PageParam{
					Pn: 1,
					Ps: 10,
				},
			},
			{
				ItemName: "测试",
				PageParam: model.PageParam{
					Pn: 1,
					Ps: 10,
				},
			},
		}

		for _, param := range params {
			res, err := s.VersionSearch(ctx, param)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		}
	})
}
