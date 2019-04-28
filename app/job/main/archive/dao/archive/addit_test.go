package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Addit(t *testing.T) {
	Convey("Addit", t, func() {
		addit, err := d.Addit(context.TODO(), 1)
		So(err, ShouldBeNil)
		Println(addit)
	})
}
