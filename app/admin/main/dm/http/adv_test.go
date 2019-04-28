package http

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestAdvList(t *testing.T) {
	params := url.Values{}
	params.Set("oid", "10108682")
	params.Set("bType", "all")
	params.Set("mode", "all")

	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := managerHost + "/x/admin/dm/adv/list?" + params.Encode()
	body, err := oget(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	t.Logf("=========%v\n", body)
	fmt.Println(string(body))
}
