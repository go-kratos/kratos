package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/model/archive"
	"testing"
)

func TestArchive(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		aid := int64(10098814)
		a, err := d.Archive(context.Background(), aid)
		So(err, ShouldBeNil)
		So(a, ShouldNotBeNil)
		t.Logf("resp: %v", a)
	}))
}

func TestArchives(t *testing.T) {
	Convey("Archives", t, WithDao(func(d *Dao) {
		_, err := d.Archives(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArcNote(t *testing.T) {
	Convey("TxUpArcNote", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcNote(tx, 111, "2")
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArcMtime(t *testing.T) {
	Convey("TxUpArcMtime", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcMtime(tx, 111)
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArcAuthor(t *testing.T) {
	Convey("TxUpArcAuthor", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcAuthor(tx, 111, 222, "222")
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArcState(t *testing.T) {
	Convey("TxUpArcState", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcState(tx, 111, 0)
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArcAccess(t *testing.T) {
	Convey("TxUpArcAccess", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcAccess(tx, 111, 0)
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArcReason(t *testing.T) {
	Convey("TxUpArcReason", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcReason(tx, 111, 0, "")
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArcAttr(t *testing.T) {
	Convey("TxUpArcAttr", t, WithDao(func(d *Dao) {
		var c = context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcAttr(tx, 111, 0, 1)
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpTag(t *testing.T) {
	Convey("TxUpTag", t, WithDao(func(d *Dao) {
		c := context.TODO()
		aid := int64(2880441)
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpTag(tx, aid, "haha1,haha2,haha3")
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpInnerAttr(t *testing.T) {
	Convey("TxUpInnerAttr", t, WithDao(func(d *Dao) {
		c := context.TODO()
		addit := &archive.Addit{
			Aid: 3,
		}
		addit.InnerAttrSet(1, archive.InnerAttrChannelReview)
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpInnerAttr(tx, addit.Aid, addit.InnerAttr)
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpArchive(t *testing.T) {
	Convey("TxUpArchive", t, WithDao(func(d *Dao) {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		a, err := d.Archive(c, 10098217)
		t.Logf("archive(%+v)", a)
		So(err, ShouldBeNil)
		if err == nil {
			_, err = d.TxUpArchive(tx, a.Aid, a.Title, "随便一个内容啦", a.Cover, "随意一个note", a.Copyright, a.PTime)
			So(err, ShouldBeNil)
		}
		tx.Commit()
	}))
}

func TestDao_TxUpArcTypeID(t *testing.T) {
	Convey("TxUpArcTypeID", t, WithDao(func(d *Dao) {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcTypeID(tx, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func TestDao_TxUpArcRound(t *testing.T) {
	Convey("TxUpArcRound", t, WithDao(func(d *Dao) {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcRound(tx, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func TestDao_TxUpArcPTime(t *testing.T) {
	Convey("TxUpArcPTime", t, WithDao(func(d *Dao) {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcPTime(tx, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func TestDao_TxUpArcCopyRight(t *testing.T) {
	Convey("TxUpArcCopyRight", t, WithDao(func(d *Dao) {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpArcCopyRight(tx, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func TestDao_ArcStateMap(t *testing.T) {
	Convey("ArcStateMap", t, WithDao(func(d *Dao) {
		c := context.TODO()
		_, err := d.ArcStateMap(c, []int64{1, 2, 3})
		So(err, ShouldBeNil)
	}))
}
