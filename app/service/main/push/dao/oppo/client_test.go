package oppo

import (
	"encoding/json"
	"testing"
	"time"

	"go-common/app/service/main/push/model"

	"github.com/smartystreets/goconvey/convey"
)

var (
	auth = &Auth{Token: "7826d9a0-f192-4402-8a2c-b0dffbdd94da", Expire: 1519336387}
	cli  = NewClient(auth, "com.bilibili.oppo.push.internal", time.Hour)
)

func init() {
	auth, _ = NewAuth("UTlf5g2bAOSQA9aCqYiFQh3X", "Ppq0B4xk73augxMJbEBSyu9m")
}

func Test_AuthExpire(t *testing.T) {
	convey.Convey("auth expire", t, func() {
		a := Auth{Expire: time.Now().Add(-8 * time.Hour).Unix()}
		if !a.IsExpired() {
			t.Errorf("access should be expire")
		}
		if auth.IsExpired() {
			t.Error("access should not be expire")
		}
	})
}

func Test_Message(t *testing.T) {
	convey.Convey("message", t, func() {
		params, _ := json.Marshal(map[string]string{
			"task_id": "123",
			"scheme":  model.Scheme(1, "2", model.PlatformAndroid, model.UnknownBuild),
		})
		m := &Message{
			Title:        "this is title",
			Content:      "this is content",
			ActionType:   ActionTypeInner,
			ActionParams: string(params),
			OfflineTTL:   3600,
		}
		res, err := cli.Message(m)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("message result(%+v)", res)
	})
}

func Test_Pushs(t *testing.T) {
	convey.Convey("pushs", t, func() {
		res, err := cli.Push("5a72f82ba250c94f9f51540d", []string{"token1", "token2"})
		convey.So(err, convey.ShouldBeNil)
		t.Logf("push result(%+v)", res)
	})
}

func Test_PushOne(t *testing.T) {
	convey.Convey("push one", t, func() {
		params, _ := json.Marshal(map[string]string{
			"task_id": "123",
			"scheme":  model.Scheme(1, "2", model.PlatformAndroid, model.UnknownBuild),
		})
		m := &Message{
			Title:        "this is title",
			Content:      "this is content",
			ActionType:   ActionTypeInner,
			ActionParams: string(params),
			OfflineTTL:   3600,
			// CallbackURL: oppo.CallbackURL(1, 123),
		}
		res, err := cli.PushOne(m, "") // baab653406d187af12daa9980c87f4e5
		convey.So(err, convey.ShouldBeNil)
		t.Logf("pushOne result(%+v)", res)
	})
}
