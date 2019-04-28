package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_LoginLogs(t *testing.T) {
	once.Do(startService)
	Convey("get login logs", t, func() {
		mid := int64(88888970)
		res, err := s.LoginLogs(context.Background(), mid, _maxLimit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
		for _, m := range res {
			str, _ := json.Marshal(m)
			t.Logf("m: %s", str)
		}
	})
}

func TestService_FormattedLoginLogs(t *testing.T) {
	once.Do(startService)
	Convey("get formatted login logs", t, func() {
		mid := int64(88888970)
		res, err := s.FormattedLoginLogs(context.Background(), mid, _maxLimit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
		for _, m := range res {
			str, _ := json.Marshal(m)
			t.Logf("m: %s", str)
		}
	})
}
