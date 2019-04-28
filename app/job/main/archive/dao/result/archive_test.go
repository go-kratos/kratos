package result

import (
	"context"
	"testing"

	"go-common/app/job/main/archive/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpPassed(t *testing.T) {
	Convey("UpPassed", t, func() {
		a, err := d.UpPassed(context.TODO(), 1684013)
		So(err, ShouldBeNil)
		Println(a)
	})
}

func Test_Archive(t *testing.T) {
	Convey("Archive", t, func() {
		a, err := d.Archive(context.TODO(), 1684013)
		So(err, ShouldBeNil)
		Println(a)
	})
}

func Test_TxAddArchive(t *testing.T) {
	Convey("TxAddArchive", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		_, err = d.TxAddArchive(context.TODO(), tx, &archive.Archive{}, &archive.Addit{}, 0, 0, "")
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_TxUpArchive(t *testing.T) {
	Convey("TxUpArchive", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		_, err = d.TxUpArchive(context.TODO(), tx, &archive.Archive{ID: 0}, &archive.Addit{}, 0, 0, "")
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}

func Test_TxDelArchive(t *testing.T) {
	Convey("TxDelArchive", t, func() {
		tx, err := d.BeginTran(context.TODO())
		So(err, ShouldBeNil)
		_, err = d.TxDelArchive(context.TODO(), tx, 0)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
	})
}
