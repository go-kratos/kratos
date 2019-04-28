package bfs

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpload(t *testing.T) {
	Convey("internal upload", t, func() {
		b := New(nil)
		req := &Request{
			Bucket:      "b",
			ContentType: "application/json",
			File:        []byte("hello world"),
		}
		res, err := b.Upload(context.TODO(), req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		t.Logf("%+v", res)
	})
}

func TestFilenameUpload(t *testing.T) {
	Convey("internal upload by specify filename", t, func() {
		b := New(nil)
		req := &Request{
			Bucket:      "b",
			Dir:         "/test",
			Filename:    "test.plain",
			ContentType: "application/json",
			File:        []byte("........"),
		}
		res, err := b.Upload(context.TODO(), req)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		t.Logf("%+v", res)
	})
}

func TestGenWatermark(t *testing.T) {
	Convey("create watermark by key and text", t, func() {
		b := New(nil)
		location, err := b.GenWatermark(context.TODO(), "c605dd5324f91ea1", "comic", "hello world", true, 2)
		So(err, ShouldBeNil)
		t.Logf("location:%s", location)
	})
}
