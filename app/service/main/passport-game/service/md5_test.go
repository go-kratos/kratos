package service

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"testing"
)

func TestMD52IntStr(t *testing.T) {
	md5d := md5.Sum([]byte("foo"))
	intStr := MD52IntStr(md5d[:])
	md5dAfter := IntStr2Md5(intStr)
	if !bytes.Equal(md5d[:], md5dAfter) {
		t.FailNow()
	}
}

func TestClout(t *testing.T) {
	fmt.Println(md5Hex(fmt.Sprintf("%s>>BiLiSaLt<<%s", md5Hex("prFP8dF2z1"), "bBTXzlsT")))
}
