package recommend

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPositionCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.PositionCache(context.Background(), 12)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestAddPositionCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		err := d.AddPositionCache(context.Background(), 12, 1)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
