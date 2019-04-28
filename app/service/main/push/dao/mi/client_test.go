package mi

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Push(t *testing.T) {
	Convey("push mi", t, func() {
		xmm := &XMMessage{
			Payload:               "bili:///?type=bililive&roomid=33886",
			RestrictedPackageName: "tv.danmaku.bili",
			PassThrough:           0, // 0 表示通知栏消息1 表示透传消息
			Title:                 model.DefaultMessageTitle,
			Description:           "直播推荐",
			NotifyType:            NotifyTypeDefaultAll,
			TaskID:                "vdsfdfs", // 每次不能相同，相同的只会推一次
		}

		// 设置是否被覆盖，不同的数字，可显示多行
		xmm.SetNotifyID(xmm.TaskID)
		xmm.SetCallbackParam("1")

		xmm.SetRegID("device token")
		// xmm.SetRegID("qlRyXrBPQ8ZkTg3x46hvTz3g8Oe/Fyz93XnE5U2NxRk=")
		// xmm.SetUserAccount("15678567,25668444")

		client := NewClient("tv.danmaku.bili", "QlcVxtNh6j7BXBPXjcbGoQ==", time.Hour)
		// client.SetProductionURL(AccountURL)
		client.SetVipURL(RegURL)
		resp, err := client.Push(xmm)
		So(err, ShouldBeNil)
		So(resp.Code, ShouldEqual, ResultCodeNoValidTargets)
		if resp.Result == ResultOk {
			tt := strings.Split(resp.Info, " ")
			if len(tt) == 6 {
				m, _ := strconv.Atoi(tt[4])
				fmt.Println(m + 1)
			}
		}
		t.Logf("push xiaomi res(%+v)", resp)
		// success: &{Result:ok Reason: Code:0 Data:{ID:scm01b20510561935064bK List:[]} Description:成功 Info:Received push messages for 1 REGID}
		// failed: &{Result:error Reason:No valid targets! Code:20301 Data:{ID: List:[]} Description:发送消息失败 Info:}
	})
}

// 需要测的时候再打开，因为失效token获取完了就没了
// func Test_InvalidTokens(t *testing.T) {
// 	client := NewClient("tv.danmaku.bili", "QlcVxtNh6j7BXBPXjcbGoQ==", time.Hour)
// 	client.SetFeedbackURL()
// 	resp, err := client.InvalidTokens()
// 	if err != nil {
// 		t.Log(err)
// 		t.FailNow()
// 	}
// 	t.Log(resp)
// }

// 需要测的时候再打开，因为卸载token获取完了就没了
// func Test_UninstalledTokens(t *testing.T) {
// 	client := NewClient("tv.danmaku.bili", "QlcVxtNh6j7BXBPXjcbGoQ==", time.Hour)
// 	resp, err := client.UninstalledTokens()
// 	if err != nil {
// 		t.Log(err)
// 		t.FailNow()
// 	}
// 	t.Log(resp)
// }
