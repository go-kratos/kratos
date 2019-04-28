package dm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	model "go-common/app/job/main/activity/model/dm"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDmBroadcast(t *testing.T) {
	convey.Convey("Broadcast", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			ds, _ = json.Marshal("ttt")
			dm    = &model.Broadcast{RoomID: 10344, CMD: "act", Info: ds}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			url := fmt.Sprintf("%s?cids=%d", d.broadcastURL, 10344)
			httpMock("POST", url).Reply(200).SetHeaders(map[string]string{
				"Code": "0",
			})
			err := d.Broadcast(c, dm)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
