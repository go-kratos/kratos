package service

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_MD52IntStr(t *testing.T) {
	Convey("MD52IntStr", t, func() {
		md5d := md5.Sum([]byte("foo"))
		intStr := MD52IntStr(md5d[:])
		md5dAfter := IntStr2Md5(intStr)
		if !bytes.Equal(md5d[:], md5dAfter) {
			t.FailNow()
		}
	})
}

func TestService_Clout(t *testing.T) {
	Convey("decompress", t, func() {
		fmt.Println(md5Hex(fmt.Sprintf("%s>>BiLiSaLt<<%s", md5Hex("prFP8dF2z1"), "bBTXzlsT")))
	})
}
