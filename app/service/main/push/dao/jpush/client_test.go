package jpush

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPush(t *testing.T) {
	Convey("test jpush", t, func() {
		var (
			ad      Audience
			notice  Notice
			plat    = NewPlatform(PlatformAndroid)
			payload = NewPayload()
			cbr     = NewCallbackReq()
			an      = &AndroidNotice{
				Title:     "test title",
				Alert:     "test alert",
				AlertType: AndroidAlertTypeLight | AndroidAlertTypeSound, // 通知提醒类型
				Extras: map[string]interface{}{
					"task_id": "tid",
					"scheme":  "bili:///?type=bilivideo&avid=123",
				},
			}
		)
		// ad.SetID([]string{"190e35f7e068f4a19d1"})
		ad.SetID([]string{""})
		notice.SetAndroidNotice(an)
		payload.SetPlatform(plat)
		payload.SetAudience(&ad)
		payload.SetNotice(&notice)
		payload.Options.SetTimelive(1000)
		payload.Options.SetReturnInvalidToken(true)
		cbr.SetParam(map[string]string{"task": "tid"})
		payload.SetCallbackReq(cbr)

		// bs, err := payload.ToBytes()
		// fmt.Printf("payload(%s) error(%v)", bs, err)

		cli := NewClient("62396b3e57f0b2b4b2c7bf48", "588f56e1bedd3c6b46db4863", time.Second)
		res, err := cli.Push(payload)
		So(err, ShouldBeNil)
		t.Logf("push result(%+v)", res)
	})
}
