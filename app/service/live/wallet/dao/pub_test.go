package dao

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testMsg struct {
	Name string
	Age  int64
}

func TestDao_Pub(t *testing.T) {
	Convey("pub", t, func() {
		once.Do(startService)
		//startSubDataBus()
		//msgs := subDataBus.Messages()

		var uid int64 = 1
		u := testMsg{
			Name: "test",
			Age:  23,
		}
		err := d.Pub(ctx, uid, &u)
		So(err, ShouldBeNil)
		t.Logf("pub")

	})
}
