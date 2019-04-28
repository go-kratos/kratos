package ugc

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_PassedArcs(t *testing.T) {
	Convey("TestDao_PassedArcs", t, WithDao(func(d *Dao) {
		res, err := d.PassedArcs(ctx, []int64{76, 75})
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		data, _ := json.Marshal(res)
		fmt.Println(string(data))
	}))
}
