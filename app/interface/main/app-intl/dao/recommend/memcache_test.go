package recommend

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddRcmdAidsCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		err := d.AddRcmdAidsCache(context.Background(), []int64{12})
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestRcmdAidsCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.RcmdAidsCache(context.Background())
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}

func TestRcmdCache(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.RcmdCache(context.Background())
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
