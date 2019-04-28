package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_AudioCard(t *testing.T) {
	convey.Convey("test audio card", t, func(ctx convey.C) {
		mid := int64(28272030)
		data, err := d.AudioCard(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}

func TestDao_AudioUpperCert(t *testing.T) {
	convey.Convey("test audio upper cert", t, func(ctx convey.C) {
		mid := int64(28272030)
		data, err := d.AudioUpperCert(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}

func TestDao_AudioCnt(t *testing.T) {
	convey.Convey("test audio cnt", t, func(ctx convey.C) {
		mid := int64(28272030)
		data, err := d.AudioCnt(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}
