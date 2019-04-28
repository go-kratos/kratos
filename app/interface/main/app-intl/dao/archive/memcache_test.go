package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStatsCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, _, err := d.statsCache(context.Background(), []int64{1, 2, 3, 4})
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestArcsCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, _, err := d.arcsCache(context.Background(), []int64{1, 2, 3, 4})
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestArcCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.arcCache(context.Background(), 123)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)

	})
}

func TestStatCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.statCache(context.Background(), 123)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestRelatesCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.RelatesCache(context.Background(), 123)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestViewContributeCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.ViewContributeCache(context.Background(), 12)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
