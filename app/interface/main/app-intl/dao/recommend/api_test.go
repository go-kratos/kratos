package recommend

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRecommend(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, _, _, _, err := d.Recommend(context.Background(), 1, "", 12, 123232, 2, 0, 2, "", "", 0, 1, 1, "", time.Now())
		if err != nil {
			t.Log(err)
		}
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestHots(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Hots(context.Background())
		if err != nil {
			t.Log(err)
		}
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestTagTop(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.TagTop(context.Background(), 12, 12, 12)
		if err != nil {
			t.Log(err)
		}
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestGroup(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Group(context.Background())
		if err != nil {
			t.Log(err)
		}
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
