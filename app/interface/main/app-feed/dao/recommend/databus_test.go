package recommend

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestPubDislike(t *testing.T) {
	var (
		c          = context.TODO()
		buvid      string
		gt         string
		id         int64
		mid        int64
		reasonID   int64
		cmReasonID int64
		feedbackID int64
		upperID    int64
		rid        int64
		tagID      int64
		adcb       string
		now        time.Time
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		d.PubDislike(c, buvid, gt, id, mid, reasonID, cmReasonID, feedbackID, upperID, rid, tagID, adcb, now)
	})
}

func TestPubDislikeCancel(t *testing.T) {
	var (
		c          = context.TODO()
		buvid      = ""
		gt         = ""
		id         = int64(1)
		mid        = int64(1)
		reasonID   = int64(1)
		cmReasonID = int64(1)
		feedbackID = int64(1)
		upperID    = int64(1)
		rid        = int64(1)
		tagID      = int64(1)
		adcb       = ""
		now        time.Time
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		d.PubDislikeCancel(c, buvid, gt, id, mid, reasonID, cmReasonID, feedbackID, upperID, rid, tagID, adcb, now)
	})
}
