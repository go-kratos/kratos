package huawei

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Push(t *testing.T) {
	Convey("push huawei", t, func() {
		// ac, err := NewAccess("10125085", "iejq6hn3ds3d4neq1m21v443lmbm31gs")
		// if err != nil {
		// 	t.Fatal(err)
		// } else {
		// 	t.Log(ac)
		// }
		// return
		ac := &Access{
			AppID:  "10125085",
			Token:  "CFrF0b079efz2JUoDNBs1lwk9wtL4LfxExYqZvM3lAuDAeZcytQS3CPjYO6qMv9h+6FJoKrGIsQEwcKOmODdeg==",
			Expire: 1522913725,
		}
		palyod := NewMessage().SetContent("huawei-content").SetTitle("huawei-title").SetCustomize("task_id", "123").SetCustomize("scheme", "bilibili://search/你好").SetIcon("http://pic.qiantucdn.com/58pic/12/38/18/13758PIC4GV.jpg")
		c := NewClient("tv.danmaku.bili", ac, time.Minute)
		// tokens := []string{"0866090037077934300001050400CN01"}
		tokens := []string{"1", "2", ""}
		res, err := c.Push(palyod, tokens, time.Now().Add(time.Hour))
		So(err, ShouldBeNil)
		t.Logf("huawei push res(%+v)", res)
	})
}
