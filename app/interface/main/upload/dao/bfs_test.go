package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/upload/conf"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBfs(t *testing.T) {
	Convey("new bfs instance", t, func() {
		b := NewBfs(&conf.Config{
			Bfs: &conf.Bfs{
				BfsURL:       "uat-bfs.bilibili.co",
				WaterMarkURL: "http://i0.hdslb.com/imageserver/watermark/gen",
				TimeOut:      xtime.Duration(time.Second * 5),
				WmTimeOut:    xtime.Duration(time.Second * 5),
			},
		})
		So(b, ShouldNotBeNil)
	})
}

func TestGenImage(t *testing.T) {
	Convey("create watermark image", t, func() {
		image, height, width, hasher, err := b.GenImage(context.TODO(), "comic", "hello world", 2, true)
		So(err, ShouldBeNil)
		So(image, ShouldNotBeEmpty)
		So(height, ShouldNotEqual, 0)
		So(width, ShouldNotEqual, 0)
		So(hasher, ShouldNotEqual, "")
	})
}

func TestWatermark(t *testing.T) {
	Convey("do watermark action", t, func() {
		image, err := b.Watermark(context.TODO(), testData, "image/png", "comic", "hello", 0, 0, 0)
		So(err, ShouldBeNil)
		So(image, ShouldNotBeEmpty)
	})
}

func TestUpload(t *testing.T) {
	Convey("upload", t, func() {
		var (
			dir      = "dir1/"
			filename = "1111.jpg"
		)
		location, _, err := b.Upload(context.Background(), "1b24a3d8560d2213", "415aaa6ff53659dabf8a2de394025a", "image/jpg", "static", dir, filename, testData)
		So(err, ShouldBeNil)
		So(location, ShouldNotBeEmpty)
	})
	Convey("upload", t, func() {
		var (
			dir      = "dir1/"
			filename = ""
		)
		location, _, err := b.Upload(context.Background(), "1b24a3d8560d2213", "415aaa6ff53659dabf8a2de394025a", "image/jpg", "static", dir, filename, testData)
		So(err, ShouldBeNil)
		So(location, ShouldNotBeEmpty)
	})
	Convey("upload", t, func() {
		var (
			dir      = ""
			filename = "1111.jpg"
		)
		location, _, err := b.Upload(context.Background(), "1b24a3d8560d2213", "415aaa6ff53659dabf8a2de394025a", "image/jpg", "static", dir, filename, testData)
		So(err, ShouldBeNil)
		So(location, ShouldNotBeEmpty)
	})
	Convey("upload", t, func() {
		var (
			dir      = ""
			filename = ""
		)
		location, _, err := b.Upload(context.Background(), "1b24a3d8560d2213", "415aaa6ff53659dabf8a2de394025a", "image/jpg", "static", dir, filename, testData)
		So(err, ShouldBeNil)
		So(location, ShouldNotBeEmpty)
	})
}
