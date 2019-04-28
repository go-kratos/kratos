package databus

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/message"
	"testing"
)

func TestPopMsgCache(t *testing.T) {
	var (
		err error
	)
	Convey("PopMsgCache", t, WithDao(func(d *Dao) {
		_, err = d.PopMsgCache(context.Background())
		So(err, ShouldBeNil)
	}))
}
func TestDao_PushMultSync(t *testing.T) {
	Convey("PushMultSync", t, WithDao(func(d *Dao) {
		c := context.TODO()
		sync := &archive.MultSyncParam{}
		_, err := d.PushMultSync(c, sync)
		So(err, ShouldBeNil)
	}))
}
func TestDao_PopMultSync(t *testing.T) {
	Convey("PopMultSync", t, WithDao(func(d *Dao) {
		c := context.TODO()
		_, err := d.PopMultSync(c)
		So(err, ShouldBeNil)
	}))
}
func TestDao_PopMsgCache(t *testing.T) {
	Convey("FlowGroupPools", t, WithDao(func(d *Dao) {
		c := context.TODO()
		_, err := d.PopMsgCache(c)
		So(err, ShouldBeNil)
	}))
}

func TestDao_PushMsgCache(t *testing.T) {
	Convey("FlowGroupPools", t, WithDao(func(d *Dao) {
		c := context.TODO()
		msg := &message.Videoup{}
		err := d.PushMsgCache(c, msg)
		So(err, ShouldBeNil)
	}))
}
