package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/interface/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSendAction(t *testing.T) {
	var (
		c   = context.TODO()
		act = &model.ReportAction{
			Cid:      int64(10106598),
			Did:      int64(719918177),
			HideTime: time.Now().Unix() + 10,
		}
	)
	Convey("", t, func() {
		err := testDao.SendAction(c, fmt.Sprint(act.Cid), act)
		So(err, ShouldBeNil)
	})
}
