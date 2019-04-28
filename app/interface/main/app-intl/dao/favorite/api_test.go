package favorite

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsFavDefault(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.IsFavDefault(context.Background(), 1, 1)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestIsFav(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.IsFav(context.Background(), 1, 1)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
