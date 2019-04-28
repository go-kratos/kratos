package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetLogin(t *testing.T) {
	Convey("TestSetLogin", t, func() {
		err := d.SetLogin(ctx, 1, 2)
		if err != nil {
			t.Errorf("dedeCoins err(%v)", err)
		}
		b, _ := d.Logined(ctx, 1, 2)
		if !b {
			t.Errorf("Logined should be true but get %v", b)
		}
		b, _ = d.Logined(ctx, 1, 3)
		So(b, ShouldNotBeNil)
	})
}
