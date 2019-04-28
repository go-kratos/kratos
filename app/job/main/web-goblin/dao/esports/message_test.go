package esports

import (
	"testing"

	mdlesp "go-common/app/job/main/web-goblin/model/esports"

	"github.com/smartystreets/goconvey/convey"
)

func TestEsportsSendMessage(t *testing.T) {
	var (
		mids    = []int64{}
		msg     = "2018 LPL 春季赛 中，你订阅的赛程“2018-01-25 19:00 RNG VS SNG”即将开播，快前去观看比赛吧！ 点击前往直播间"
		contest = &mdlesp.Contest{}
	)
	mids = append(mids, 111)
	mids = append(mids, 222)
	convey.Convey("SendMessage", t, func(ctx convey.C) {
		err := d.SendMessage(mids, msg, contest)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestEsportsBatch(t *testing.T) {
	var (
		list      = []int64{}
		msg       = ""
		contest   = &mdlesp.Contest{}
		batchSize = int(0)
		f         func(users []int64, msg string, contest *mdlesp.Contest) error
	)
	convey.Convey("Batch", t, func(ctx convey.C) {
		d.Batch(list, msg, contest, batchSize, f)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}
