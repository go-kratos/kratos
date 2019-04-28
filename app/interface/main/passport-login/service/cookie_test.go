package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_AddCookie(t *testing.T) {
	once.Do(startService)
	convey.Convey("Encrypt  param ", t, func() {
		res, _ := s.AddCookie(context.Background(), 1)

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}

func TestService_AddOldCookie(t *testing.T) {
	once.Do(startService)
	convey.Convey("Encrypt  param ", t, func() {
		res, _ := s.AddOldCookie(context.Background(), 1)

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}

func TestService_oldSession(t *testing.T) {
	once.Do(startService)
	convey.Convey("Encrypt  param ", t, func() {
		session := s.oldSession(1, time.Now().Unix(), 10)
		fmt.Println(session)
	})
}
