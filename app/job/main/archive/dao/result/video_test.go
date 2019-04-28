package result

import (
	"context"
	"testing"

	"go-common/app/job/main/archive/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_TxAddVideo(t *testing.T) {
	Convey("Archive", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		_, err = d.TxAddVideo(context.TODO(), tx, &archive.Video{Aid: 1, Cid: 1})
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_TxDelVideoByCid(t *testing.T) {
	Convey("TxDelVideoByCid", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		_, err = d.TxDelVideoByCid(context.TODO(), tx, 1, 1)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_TxDelVideos(t *testing.T) {
	Convey("TxDelVideos", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		_, err = d.TxDelVideos(context.TODO(), tx, 0)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}
