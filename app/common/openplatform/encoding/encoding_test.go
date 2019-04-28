package encoding

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var c = &EncryptConfig{Key: "bilibili_key_GYl", IV: "biliBiliIv123456"}

const (
	plat = "13291831554"
	enc  = "dWLVOucFkBBmmXmTm6so5A=="
)

func TestEncrypt(t *testing.T) {
	convey.Convey("Encrypt", t, func() {
		s, _ := Encrypt(plat, c)
		convey.So(s, convey.ShouldEqual, enc)
	})
}

func TestDecrypt(t *testing.T) {
	convey.Convey("Decrypt", t, func() {
		s, _ := Decrypt(enc, c)
		convey.So(s, convey.ShouldEqual, plat)
	})
}
