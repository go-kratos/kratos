package history

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestHistorykeyFirst(t *testing.T) {
	convey.Convey("keyFirst", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
			no  = "14771787"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			keyFirst(mid, no)
		})
	})
}

func TestHistoryPushFirstQueue(t *testing.T) {
	convey.Convey("PushFirstQueue", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			aid = int64(1222)
			now = time.Now().Unix()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.PushFirstQueue(c, mid, aid, now)
		})
	})
}
