package email

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/email"
	"go-common/library/cache/redis"
)

var member = email.Retry{
	Action: email.RetryActionReply,
	AID:    11,
	Flag:   archive.ReplyOn,
	FlagA:  archive.ReplyDefault,
}
var key = email.RetryListKey

func TestEmailPushRedis(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("PushRedis", t, func(ctx convey.C) {
		err := d.PushRedis(c, member, key)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestEmailPopRedis(t *testing.T) {
	var (
		c = context.TODO()
	)
	TestEmailPushRedis(t)
	convey.Convey("PopRedis", t, func(ctx convey.C) {
		bs, err := d.PopRedis(c, key)
		ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(bs, convey.ShouldNotBeNil)
		})
	})
}

func TestEmailRemoveRedis1(t *testing.T) {
	var (
		c = context.TODO()
	)

	TestEmailPushRedis(t)
	convey.Convey("RemoveRedis a member", t, func(ctx convey.C) {
		bs, _ := json.Marshal(member)
		reply, err := d.RemoveRedis(c, key, string(bs))
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(reply, convey.ShouldBeGreaterThan, 0)
		})
	})
}

func TestEmailRemoveRedis2(t *testing.T) {
	var (
		c     = context.TODO()
		err   error
		reply int
	)

	key1 := "list"
	bsList := []interface{}{"ah1", "ah3", "ah5"}
	reply, err = d.RemoveRedis(c, key1, bsList...)
	t.Logf("function reply(%v) err(%v)", reply, err)

	//多个元素删除
	var bs []byte
	list := []interface{}{}
	old := []int64{-1, 0, 1}
	nw := []int64{0, 1}
	for _, v := range old {
		for _, j := range nw {
			m := email.Retry{
				AID:    member.AID,
				Action: member.Action,
				Flag:   j,
				FlagA:  v,
			}
			if bs, err = json.Marshal(m); err != nil {
				continue
			}
			list = append(list, string(bs))
		}
	}
	TestEmailPushRedis(t)
	member.Flag = 1
	TestEmailPushRedis(t)

	convey.Convey("RemoveRedis many member", t, func(ctx convey.C) {
		reply, err = d.RemoveRedis(c, key, list...)
		t.Logf("member reply(%d) error(%v)", reply, err)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})

		bs, err := d.PopRedis(c, key)
		convey.So(err, convey.ShouldEqual, redis.ErrNil)
		convey.So(bs, convey.ShouldBeNil)
	})

}
