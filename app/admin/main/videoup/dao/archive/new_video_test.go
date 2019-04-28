package archive

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_videoSrcTypeByIDs(t *testing.T) {
	convey.Convey("根据cid获取最新的上传类型src_type", t, WithDao(func(d *Dao) {
		ids := []int64{385, 386, 387, 388}
		m, err := d.VideoSrcTypeByIDs(context.TODO(), ids)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(m), convey.ShouldBeLessThanOrEqualTo, len(ids))
	}))
}

func Test_vIDByAIDFilename(t *testing.T) {
	convey.Convey("根据aid+filename获取分p的vid", t, WithDao(func(d *Dao) {
		aid := int64(161)
		filename := "d74b1c1cda32e5740658a2517fd82965"
		_, err := d.VIDByAIDFilename(context.TODO(), aid, filename)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func TestDao_VideoInfo(t *testing.T) {
	convey.Convey("VideoInfo", t, WithDao(func(d *Dao) {
		_, err := d.VideoInfo(context.Background(), 10098493, 10109201)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func TestDao_VideoByCID(t *testing.T) {
	convey.Convey("VideoByCID", t, WithDao(func(d *Dao) {
		info, err := d.VideoByCID(context.Background(), 10109201)
		convey.So(err, convey.ShouldBeNil)
		convey.So(info, convey.ShouldNotBeNil)
	}))
}

func TestDao_VideoRelated(t *testing.T) {
	convey.Convey("VideoRelated", t, WithDao(func(d *Dao) {
		v, err := d.VideoRelated(context.Background(), 10098493)
		t.Logf("VideoRelated(%+v)\r\n", v)
		convey.So(err, convey.ShouldBeNil)
		convey.So(v, convey.ShouldNotBeNil)
	}))
}

func Test_TxUpRelation(t *testing.T) {
	var c = context.Background()
	convey.Convey("TxUpRelation", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpRelation(tx, 0, "", "")
		convey.So(err, convey.ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpRelationOrder(t *testing.T) {
	var c = context.Background()
	convey.Convey("TxUpRelationOrder", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpRelationOrder(tx, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpRelationState(t *testing.T) {
	var c = context.Background()
	convey.Convey("TxUpRelationState", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpRelationState(tx, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpWebLink(t *testing.T) {
	var c = context.Background()
	convey.Convey("TxUpWebLink", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpWebLink(tx, 0, "")
		convey.So(err, convey.ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpStatus(t *testing.T) {
	var c = context.Background()
	convey.Convey("TxUpStatus", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpStatus(tx, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		tx.Commit()
	}))
}

func Test_TxUpAttr(t *testing.T) {
	var c = context.Background()
	convey.Convey("TxUpAttr", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpAttr(tx, 0, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		tx.Commit()
	}))
}

func Test_VideoPlayurl(t *testing.T) {
	var c = context.Background()
	convey.Convey("VideoPlayurl", t, WithDao(func(d *Dao) {
		_, err := d.VideoPlayurl(c, 0)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func Test_NewVideoByID(t *testing.T) {
	var c = context.Background()
	convey.Convey("NewVideoByID", t, WithDao(func(d *Dao) {
		_, err := d.NewVideoByID(c, 0)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func Test_NewVideoByIDs(t *testing.T) {
	var c = context.Background()
	convey.Convey("NewVideoByIDs", t, WithDao(func(d *Dao) {
		_, err := d.NewVideoByIDs(c, []int64{1, 2, 3})
		convey.So(err, convey.ShouldBeNil)
	}))
}

func Test_NewVideosByAid(t *testing.T) {
	var c = context.Background()
	convey.Convey("NewVideosByAid", t, WithDao(func(d *Dao) {
		_, err := d.NewVideosByAid(c, 0)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func Test_NewVideoCount(t *testing.T) {
	var c = context.Background()
	convey.Convey("NewVideoCount", t, WithDao(func(d *Dao) {
		_, err := d.NewVideoCount(c, 0)
		convey.So(err, convey.ShouldBeNil)
	}))
}
