package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/dao"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Retry(t *testing.T) {
	Convey("retry", t, WithoutProcService(func(s *Service) {

		var (
			err error
			c   = context.TODO()
		)

		SkipConvey("push retry", func() {
			var (
				favCnt    = int64(1)
				replyCnt  = int64(2)
				statRetry = &dao.StatRetry{
					Data: &artmdl.StatMsg{
						Aid:      888,
						Favorite: &favCnt,
						Reply:    &replyCnt,
					},
					Action: dao.RetryUpdateStatCache,
					Count:  3,
				}
			)
			err = s.dao.PushStat(c, statRetry)
			So(err, ShouldBeNil)
		})

		SkipConvey("pop retry", func() {
			bs, err := s.dao.PopStat(c)
			if err != redis.ErrNil {
				So(err, ShouldBeNil)
			}
			So(bs, ShouldNotBeEmpty)
			msg := &dao.StatRetry{}
			err = json.Unmarshal(bs, msg)
			So(err, ShouldBeNil)
			So(msg.Action, ShouldEqual, dao.RetryUpdateStatCache)
			So(msg.Count, ShouldEqual, 3)
		})

		Convey("retry reply", func() {
			err := s.dao.PushReply(context.TODO(), 8, 9)
			So(err, ShouldBeNil)
			aid, mid, err := s.dao.PopReply(c)
			So(err, ShouldBeNil)
			fmt.Println(aid)
			fmt.Println(mid)
		})
	}))
}
