package ugc

import (
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_DeletedVideos(t *testing.T) {
	Convey("TestDao_DeletedVideos", t, WithDao(func(d *Dao) {
		res, err := d.DeletedVideos(ctx)
		if err == nil && len(res) == 0 {
			fmt.Println("No Delete Data")
			d.DB.Exec(ctx, "UPDATE ugc_video SET deleted = 1, submit =1 WHERE deleted = 0 LIMIT 1")
		}
		res, err = d.DeletedVideos(ctx)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		fmt.Println(res)
	}))
}
