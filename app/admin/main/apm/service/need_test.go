package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/apm/dao/mock"
	"go-common/app/admin/main/apm/model/need"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

//TestServiceNeedList is
func TestServiceNeedList(t *testing.T) {
	convey.Convey("TestServiceNeedList", t, func() {
		arg := &need.NListReq{
			Status: 1,
			Ps:     10,
			Pn:     1,
		}
		guard := mock.MockDaoNeedInfoCount(svr.dao, 0, nil)
		defer guard.Unpatch()
		res, _, err := svr.NeedInfoList(context.Background(), arg, "fengshanshan")
		for _, v := range res {
			t.Logf("res:%+v", v)
		}
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)

	})
}

//TestServiceNeedEdit is
func TestServiceNeedEdit(t *testing.T) {
	convey.Convey("TestServiceNeedEdit", t, func() {
		arg := &need.NEditReq{
			ID:    147,
			Title: "22222",
		}
		err := svr.NeedInfoEdit(context.Background(), arg, "fengshanshan")
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("TestServiceNeedEdit status err", t, func() {
		arg := &need.NEditReq{
			ID:    99,
			Title: "22222",
		}
		err := svr.NeedInfoEdit(context.Background(), arg, "fengshanshan")
		convey.So(err, convey.ShouldEqual, 70018)
	})
	convey.Convey("TestServiceNeedEdit no access", t, func() {
		arg := &need.NEditReq{
			ID:    147,
			Title: "22222",
		}
		err := svr.NeedInfoEdit(context.Background(), arg, "fss")
		convey.So(err, convey.ShouldEqual, -403)
	})

}

//TestServiceNeedVerify is
func TestServiceNeedVerify(t *testing.T) {
	convey.Convey("TestServiceNeedVerify status err", t, func() {
		arg := &need.NVerifyReq{
			ID:     117,
			Status: 2,
		}
		_, err := svr.NeedInfoVerify(context.Background(), arg)
		convey.So(err, convey.ShouldEqual, 70019)
	})
	convey.Convey("TestServiceNeedVerify not exist", t, func() {
		arg := &need.NVerifyReq{
			ID:     10000,
			Status: 1,
		}
		_, err := svr.NeedInfoVerify(context.Background(), arg)
		convey.So(err, convey.ShouldEqual, 70017)

	})

}

func TestServiceNeedVote(t *testing.T) {
	convey.Convey("TestServiceNeedVote", t, func() {
		arg := &need.Likereq{
			ReqID:    148,
			LikeType: 1,
		}
		err := svr.NeedInfoVote(context.Background(), arg, "fengshanshan")
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("TestServiceNeedVote not exist", t, func() {
		arg := &need.Likereq{
			ReqID:    789,
			LikeType: 1,
		}
		err := svr.NeedInfoVote(context.Background(), arg, "fengshanshan")
		convey.So(err, convey.ShouldEqual, 70017)
	})
}

func TestServiceNeedVoteList(t *testing.T) {
	convey.Convey("NeedVoteList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &need.Likereq{
				ReqID:    11,
				LikeType: 1,
			}
			resp = []*need.UserLikes{
				{ID: 1, User: "fengshanshan", LikeType: 1},
				{ID: 2, User: "fengshanshan", LikeType: 2},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			svr.dao.MockVoteInfoCounts(6, nil)
			svr.dao.MockVoteInfoList(resp, nil)
			res, count, err := svr.NeedVoteList(c, arg)
			ctx.Convey("Then err should be nil.res,count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
			ctx.Reset(func() {
				monkey.UnpatchAll()
			})
		})
	})
}
func TestServiceNeedInfoAdd(t *testing.T) {
	convey.Convey("NeedInfoAdd", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &need.NAddReq{
				Title:   "32323",
				Content: "ewewerw",
			}
			username = "fengshanshan"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			guard := svr.dao.MockNeedInfoAdd(nil)
			defer guard.Unpatch()
			err := svr.NeedInfoAdd(c, req, username)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When return err", func(ctx convey.C) {
			guard := svr.dao.MockNeedInfoAdd(fmt.Errorf("aaa"))
			defer guard.Unpatch()
			err := svr.NeedInfoAdd(c, req, username)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestServiceSendWeMessage(t *testing.T) {
	convey.Convey("SendWeMessage", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			title    = "有个xxxxxxx"
			task     = need.VerifyType[need.NeedReview]
			result   = need.VerifyType[need.VerifyAccept]
			sender   = "fengshanshan"
			receiver = []string{"fengshanshan"}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := svr.SendWeMessage(c, title, task, result, sender, receiver)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
