package http

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"
)

func TestUpFilters(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "27515615")
	params.Set("type", "3")
	params.Set("pn", "1")
	params.Set("ps", "5")
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := managerHost + "/x/admin/dm/upfilter/list?" + params.Encode()
	fmt.Println(s)
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
}
