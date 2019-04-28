package service

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceBase64Decode(t *testing.T) {
	convey.Convey("", t, func() {
		bytes, err := base64.StdEncoding.DecodeString("igIDgs/yFxaFI+oiu2HoDw==")
		convey.So(bytes, convey.ShouldNotBeEmpty)
		convey.So(err, convey.ShouldBeNil)
		convey.So(hex.EncodeToString(bytes), convey.ShouldEqual, "8a020382cff217168523ea22bb61e80f")
	})
}

func TestService_cleanTokenCache(t *testing.T) {
	once.Do(startService)
	convey.Convey("cleanTokenCache", t, func() {
		err := s.cleanTokenCache("igIDgs/yFxaFI+oiu2HoDw==", 0)
		convey.So(err, convey.ShouldNotBeNil)
	})
}
