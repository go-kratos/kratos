package archive

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_Watermark(t *testing.T) {
	convey.Convey("水印", t, WithDao(func(d *Dao) {
		m, err := d.Watermark(context.TODO(), 1)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("watermark(%+v)", m)
	}))
}
