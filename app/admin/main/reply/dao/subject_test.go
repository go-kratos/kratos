package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func deleteSubject(d *Dao, oid int64, typ int32) error {
	_delSQL := "Delete from reply_subject_%d where oid=? and type=?"
	_, err := d.db.Exec(context.Background(), fmt.Sprintf(_delSQL, subHit(oid)), oid, typ)
	if err != nil {
		return err
	}
	return nil
}

func TestSubjectAttr(t *testing.T) {
	var (
		sub = &model.Subject{
			Oid:  1,
			Type: 1,
		}
		now = time.Now()
		c   = context.Background()
	)
	Convey("get and update a subject", t, WithDao(func(d *Dao) {
		_, err := d.AddSubject(c, sub.Mid, sub.Oid, sub.Type, sub.State, now)
		So(err, ShouldBeNil)
		sub.AttrSet(model.AttrYes, model.SubAttrMonitor)
		_, err = d.UpSubjectAttr(c, sub.Oid, sub.Type, sub.Attr, now)
		So(err, ShouldBeNil)
		sub, err = d.Subject(c, sub.Oid, sub.Type)
		So(err, ShouldBeNil)
		So(sub.AttrVal(model.SubAttrMonitor), ShouldEqual, model.AttrYes)
	}))
	Convey("get and update a subject", t, WithDao(func(d *Dao) {
		_, err := d.AddSubject(c, sub.Mid, sub.Oid, sub.Type, sub.State, now)
		So(err, ShouldBeNil)
		sub.AttrSet(model.AttrNo, model.SubAttrMonitor)
		_, err = d.UpSubjectAttr(c, sub.Oid, sub.Type, sub.Attr, now)
		So(err, ShouldBeNil)
		sub, err = d.Subject(c, sub.Oid, sub.Type)
		So(err, ShouldBeNil)
		So(sub.AttrVal(model.SubAttrMonitor), ShouldEqual, model.AttrNo)
	}))
}

func TestSubjectState(t *testing.T) {
	var (
		sub = &model.Subject{
			Oid:  1,
			Type: 1,
		}
		now = time.Now()
		c   = context.Background()
	)
	Convey("get and update a subject", t, WithDao(func(d *Dao) {
		_, err := d.AddSubject(c, sub.Mid, sub.Oid, sub.Type, sub.State, now)
		So(err, ShouldBeNil)
		_, err = d.UpSubjectState(c, sub.Oid, sub.Type, model.SubStateForbid, now)
		So(err, ShouldBeNil)
		sub, err = d.Subject(c, sub.Oid, sub.Type)
		So(sub.State, ShouldEqual, model.SubStateForbid)
	}))
}

func TestSubjectCount(t *testing.T) {
	d := _d
	var (
		sub = &model.Subject{
			Oid:  2,
			Type: 3,
		}
		now = time.Now()
		c   = context.Background()
	)
	_, err := d.AddSubject(c, sub.Mid, sub.Oid, sub.Type, sub.State, now)
	if err != nil {
		t.Logf("add subject failed!%v", err)
		t.FailNow()
	}
	defer deleteSubject(d, 2, 3)
	Convey("increase and decrease subject", t, WithDao(func(d *Dao) {
		Convey("increase subject rcount", WithDao(func(d *Dao) {
			t, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			rows, err := d.TxIncrSubRCount(t, 2, 3, now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
			err = t.Commit()
			So(err, ShouldBeNil)
			sub, err = d.Subject(c, sub.Oid, sub.Type)
			So(err, ShouldBeNil)
			So(sub.RCount, ShouldEqual, 1)
		}))
		Convey("decrease subject rcount", WithDao(func(d *Dao) {
			t, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			rows, err := d.TxDecrSubRCount(t, 2, 3, now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
			err = t.Commit()
			So(err, ShouldBeNil)
			sub, err = d.Subject(c, sub.Oid, sub.Type)
			So(err, ShouldBeNil)
			So(sub.RCount, ShouldEqual, 0)
		}))
		Convey("increase subject acount", WithDao(func(d *Dao) {
			t, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			rows, err := d.TxIncrSubACount(t, 2, 3, 123, now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
			err = t.Commit()
			So(err, ShouldBeNil)
			sub, err = d.Subject(c, sub.Oid, sub.Type)
			So(err, ShouldBeNil)
			So(sub.ACount, ShouldEqual, 123)
		}))
		Convey("decrease subject acount", WithDao(func(d *Dao) {
			t, err := d.BeginTran(c)
			So(err, ShouldBeNil)
			rows, err := d.TxSubDecrACount(t, 2, 3, 23, now)
			So(err, ShouldBeNil)
			So(rows, ShouldBeGreaterThan, 0)
			err = t.Commit()
			So(err, ShouldBeNil)
			sub, err = d.Subject(c, sub.Oid, sub.Type)
			So(err, ShouldBeNil)
			So(sub.ACount, ShouldEqual, 100)
		}))
	}))
}
