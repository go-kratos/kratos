package space

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ClipList(t *testing.T) {
	Convey("ClipList", t, func() {
		_, _, _, err := d.ClipList(context.TODO(), 1, 1, 10, "")
		So(err, ShouldBeNil)
	})
}

func Test_AlbumList(t *testing.T) {
	Convey("AlbumList", t, func() {
		d.AlbumList(context.TODO(), 1, 1, 10, "")
	})
}

func Test_AudioList(t *testing.T) {
	Convey("AudioList", t, func() {
		d.AudioList(context.TODO(), 1, 1, 1, "")
	})
}
