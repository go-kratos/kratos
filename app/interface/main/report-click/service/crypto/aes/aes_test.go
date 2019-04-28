package aes

import (
	"testing"

	"go-common/app/interface/main/report-click/service/crypto/padding"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAES(t *testing.T) {
	var (
		aesKey = ""
		aesIv  = ""
		bs     = []byte("this is  test massage ")
	)
	Convey("aes test ", t, func() {
		bs, _ = CBCDecrypt(bs, []byte(aesKey), []byte(aesIv), padding.PKCS5)
		ECBDecrypt(bs, []byte(aesKey), padding.PKCS5)
		bs = []byte("this is  test massage ")
		bs, _ = CBCEncrypt(bs, []byte(aesKey), []byte(aesIv), padding.PKCS5)
		CBCDecrypt(bs, []byte(aesKey), []byte(aesIv), padding.PKCS5)
	})
}
