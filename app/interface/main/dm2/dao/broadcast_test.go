package dao

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBroadcast(t *testing.T) {
	Convey("err should be nil", t, func() {
		var err error
		info := []byte(fmt.Sprintf(`["%.2f,%d,%d,%d,%d,%d,%d,%s,%d","%s"]`, 1.01, 0, 25, 0xFF, 1123, 23, 0, "abc", 1, `德玛西亚万岁\\(≧▽≦)/`))
		if err = testDao.BroadcastInGoim(c, 123, 456, info); err != nil {
			t.Logf("testDao.Broadcast(%s) error(%v)", info, err)
		}
		So(err, ShouldBeNil)
	})
}
