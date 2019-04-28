package music

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMusicAudio(t *testing.T) {
	var (
		c     = context.TODO()
		ids   = []int64{1, 2, 3}
		level = int(0)
		ip    = "127.0.0.1"
	)
	convey.Convey("Audio", t, func(ctx convey.C) {
		au, err := d.Audio(c, ids, level, ip)
		ctx.Convey("Then err should be nil.au should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(au, convey.ShouldNotBeNil)
		})
	})
}
