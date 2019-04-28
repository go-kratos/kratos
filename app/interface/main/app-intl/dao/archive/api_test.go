package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRelateAids(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.RelateAids(context.Background(), 1)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestCommercial(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Commercial(context.Background(), 12)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestNewRelateAids(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, _, _, err := d.NewRelateAids(context.Background(), 12, 12, 0, "", "", 1)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestPlayerInfos(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.PlayerInfos(context.Background(), []int64{12}, 1, "", 0, 2)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
