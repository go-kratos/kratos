package fcm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode"

	. "github.com/smartystreets/goconvey/convey"
)

const apiKey = "AIzaSyBtMplqJkuTIDyIx-CM74MoPHbxHCBcYYQ"

func TestPush(t *testing.T) {
	Convey("test jpush", t, func() {

		data := map[string]string{
			"task_id": "123456",
			// "scheme":  model.Scheme(model.LinkTypeVideo, "123", model.PlatformAndroid, 390000),
			"scheme": "bilibili://video/123",
		}
		client := NewClient(apiKey, 5*time.Second)
		message := &Message{
			// DryRun:          true, // 如果是 true，消息不会下发给用户，用于测试
			Data:            data,
			RegistrationIDs: []string{"fpICefK-jfE:APA91bHjZTxe503tpFoFMmXXX9LAiMmg7OwgTPYmTb8Ox-yF88umTQnmTQUGbALplxqre7R6v3d0-vSK5MyT4jFtSqklbY1GIaM4d8uZ0wJlwWrRWdBDeOJ4rlpvamd3aGyBlHKAH18N"},
			Priority:        PriorityHigh,
			DelayWhileIdle:  true,
			Notification: Notification{
				Title:       "Hello",
				Body:        "World",
				ClickAction: "com.bilibili.app.in.com.bilibili.push.FCM_MESSAGE",
			},
			CollapseKey: strings.TrimFunc("t123456", func(r rune) bool {
				return !unicode.IsNumber(r)
			}), // 值转成 int 传到客户端
			TimeToLive: int(time.Hour.Seconds()),
			Android:    Android{Priority: PriorityHigh},
		}
		response, err := client.Send(message)
		msgb, _ := json.Marshal(message)
		fmt.Printf("msg(%s)", msgb)
		So(err, ShouldNotBeNil)
		if err != nil {
			t.Errorf("fcm send response(%+v) error(%v)", response, err)
		} else {
			fmt.Println("Status Code   :", response.StatusCode)
			fmt.Println("Success       :", response.Success)
			fmt.Println("Fail          :", response.Fail)
			fmt.Println("Canonical_ids :", response.CanonicalIDs)
			fmt.Println("Topic MsgId   :", response.MsgID)
		}
	})
}

func Test_ClientFaild(t *testing.T) {
	Convey("test jpush", t, func() {
		client := NewClient(apiKey, 5*time.Second)
		err := client.Failed(&Response{})
		So(err, ShouldBeNil)
		r := &Response{RetryAfter: "3m"}
		_, err = r.GetRetryAfterTime()
		So(err, ShouldBeNil)

	})
}
