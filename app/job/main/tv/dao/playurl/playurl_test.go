package playurl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Playurl(t *testing.T) {
	Convey("TestDao_Playurl", t, WithDao(func(d *Dao) {
		url, hitDead, err := d.Playurl(ctx, 41057587)
		So(err, ShouldBeNil)
		So(hitDead, ShouldBeTrue)
		So(url, ShouldNotBeEmpty)
	}))
}
