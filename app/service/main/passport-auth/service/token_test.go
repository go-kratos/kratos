package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("test tmp", t, func() {

		Convey("get month int", func() {
			So(int(time.Unix(1526980273, 0).Month()), ShouldEqual, 5)
		})
	})
}

func TestService_TokenInfo(t *testing.T) {
	once.Do(startService)
	Convey("Test Query Token", t, func() {
		token := "2ee11df3a6ba1b7f04a5a15336c2a051"
		res, err := s.tokenInfo(context.TODO(), token)
		So(err, ShouldBeNil)
		So(res.Token, ShouldEqual, token)

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}
