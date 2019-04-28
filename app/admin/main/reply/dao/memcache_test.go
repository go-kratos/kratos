package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubjectCache(t *testing.T) {
	var (
		sub = &model.Subject{
			Oid:  1,
			Type: 1,
		}
		c = context.Background()
	)
	Convey("test subject cache", t, WithDao(func(d *Dao) {
		// add
		err := d.AddSubjectCache(c, sub)
		So(err, ShouldBeNil)
		// get
		cache, err := d.SubjectCache(c, sub.Oid, sub.Type)
		So(err, ShouldBeNil)
		So(cache.Oid, ShouldEqual, sub.Oid)
		// del
		err = d.DelSubjectCache(c, sub.Oid, sub.Type)
		So(err, ShouldBeNil)
		// get
		cache, err = d.SubjectCache(c, sub.Oid, sub.Type)
		So(err, ShouldBeNil)
		So(cache, ShouldBeNil)
	}))
}

func TestReplyCache(t *testing.T) {
	var (
		rp = &model.Reply{
			ID:   1,
			Oid:  1,
			Type: 1,
		}
		c = context.Background()
	)
	Convey("test reply cache", t, WithDao(func(d *Dao) {
		// add
		err := d.AddReplyCache(c, rp)
		So(err, ShouldBeNil)
		// get
		cache, err := d.ReplyCache(c, rp.ID)
		So(err, ShouldBeNil)
		So(cache.ID, ShouldEqual, rp.ID)
		// get
		caches, miss, err := d.RepliesCache(c, []int64{rp.ID})
		So(err, ShouldBeNil)
		So(len(caches), ShouldEqual, 1)
		So(len(miss), ShouldEqual, 0)
		// del
		err = d.DelReplyCache(c, rp.ID)
		So(err, ShouldBeNil)
		// get
		cache, err = d.ReplyCache(c, rp.ID)
		So(err, ShouldBeNil)
		So(cache, ShouldBeNil)
	}))
	Convey("test top reply cache", t, WithDao(func(d *Dao) {
		rp.AttrSet(model.AttrYes, model.AttrTopAdmin)
		// add
		err := d.AddTopCache(c, rp)
		So(err, ShouldBeNil)
		// get
		cache, err := d.TopCache(c, rp.Oid, model.SubAttrTopAdmin)
		So(err, ShouldBeNil)
		So(cache.ID, ShouldEqual, rp.ID)
		// del
		err = d.DelTopCache(c, rp.Oid, model.SubAttrTopAdmin)
		// get
		cache, err = d.TopCache(c, rp.Oid, model.SubAttrTopAdmin)
		So(err, ShouldBeNil)
		So(cache, ShouldBeNil)
	}))
}
