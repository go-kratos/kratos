package block

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/member/model/block"
	"testing"
)

func TestUserDetails(t *testing.T) {
	var (
		mids    = []int64{1, 2}
		details map[int64]*block.BlockUserDetail
		err     error
	)
	convey.Convey("get block user details by mids", t, func() {
		details, err = s.UserDetails(context.Background(), mids)
		convey.So(err, convey.ShouldBeNil)
		convey.So(details, convey.ShouldNotBeNil)
	})
}

func TestInfos(t *testing.T) {
	var (
		mids  = []int64{1, 2}
		infos []*block.BlockInfo
		err   error
	)
	convey.Convey("get block user details by mids", t, func() {
		infos, err = s.Infos(context.Background(), mids)
		convey.So(err, convey.ShouldBeNil)
		convey.So(infos, convey.ShouldNotBeNil)
	})
}
