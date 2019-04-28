package dao

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/satori/go.uuid"

	"go-common/app/admin/main/usersuit/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_RangeInvites(t *testing.T) {
	mid := int64(2)
	Convey("add invite", t, func() {
		now := time.Now().Unix()
		inv := &model.Invite{
			Mid:     mid,
			Code:    uuid.NewV4().String(),
			IPng:    net.ParseIP("127.0.0.1"),
			Expires: now + int64(time.Hour*72/time.Second),
			Ctime:   xtime.Time(now),
		}
		affected, err := d.AddIgnoreInvite(context.Background(), inv)
		So(err, ShouldBeNil)
		So(affected, ShouldEqual, 1)
	})

	Convey("range a test account's current month invite codes", t, func() {
		now := time.Now()
		start, end := rangeMonth(now)
		res, err := d.RangeInvites(context.Background(), mid, start, end)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func rangeMonth(now time.Time) (start, end time.Time) {
	year := now.Year()
	month := now.Month()
	loc := now.Location()
	start = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end = time.Date(year, month+1, 0, 23, 59, 59, 0, loc)
	return
}
