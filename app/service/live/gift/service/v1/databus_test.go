package v1

import (
	"context"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/smartystreets/goconvey/convey"
)

func TestV1SendAddGiftMsg(t *testing.T) {
	convey.Convey("SendAddGiftMsg", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(88895029)
			giftID   = int64(1)
			giftNum  = int64(3)
			expireAt = int64(0)
			source   = "test"
			msgID    = uuid.NewV4().String()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			s.SendAddGiftMsg(ctx, uid, giftID, giftNum, expireAt, source, msgID)
			c.Convey("No return values", func(c convey.C) {
			})
		})
	})
}
