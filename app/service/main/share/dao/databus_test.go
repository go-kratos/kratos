package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/share/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPubShare(t *testing.T) {
	convey.Convey("PubShare", t, func(ctx convey.C) {
		p := &model.ShareParams{
			OID: int64(1),
			MID: int64(1),
			TP:  int(3),
		}
		err := d.PubShare(context.Background(), p)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoPubStatShare(t *testing.T) {
	convey.Convey("PubStatShare", t, func(ctx convey.C) {
		oid := int64(1)
		count := int64(666)
		err := d.PubStatShare(context.Background(), model.ArchiveMsgTyp, oid, count)
		convey.So(err, convey.ShouldBeNil)
	})
}
