package push

import (
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_CallPush(t *testing.T) {
	Convey("Http error", t, WithDao(func(d *Dao) {
		httpMock("POST", d.c.Cfg.Push.URL).Reply(-400).JSON("")
		err := d.CallPush(ctx, "ios", "", "")
		So(err, ShouldNotBeNil)
		fmt.Println(err)
	}))
	Convey("Business Code error", t, WithDao(func(d *Dao) {
		httpMock("POST", d.c.Cfg.Push.URL).Reply(200).JSON(`{"code":-400}"`)
		err := d.CallPush(ctx, "ios", "", "")
		So(err, ShouldNotBeNil)
		fmt.Println(err)
	}))
	Convey("Everything is fine", t, WithDao(func(d *Dao) {
		httpMock("POST", d.c.Cfg.Push.URL).Reply(200).JSON(`{"code" : 0}`)
		err := d.CallPush(ctx, "ios", "", "")
		So(err, ShouldBeNil)
	}))
}

func TestDao_DiffFinish(t *testing.T) {
	Convey("TestDao_DiffFinish", t, WithDao(func(d *Dao) {
		data, err := d.DiffFinish(ctx, "57")
		So(err, ShouldBeNil)
		So(data, ShouldBeTrue)
	}))
}

func TestDao_PushMsg(t *testing.T) {
	Convey("TestDao_DiffFinish", t, WithDao(func(d *Dao) {
		data, err := d.PushMsg(ctx, "57")
		fmt.Println(data)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	}))
}

func TestDao_Platform(t *testing.T) {
	Convey("TestDao_Platform", t, WithDao(func(d *Dao) {
		data, err := d.Platform(ctx, "77")
		fmt.Println(data)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	}))
}
