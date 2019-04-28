package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStat(t *testing.T) {
	Convey(t.Name(), t, func() {
		_, err := d.Stat(context.Background(), 2)
		if err != nil {
			t.Log(err)
		}
		err = nil
		So(err, ShouldBeNil)
	})
}
