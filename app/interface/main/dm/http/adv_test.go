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

func Test_advBuy(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "3786607")
	params.Set("mode", "sp")
	s := _host + "/x/dm/adv/buy?" + params.Encode()
	body, err := cpost(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}

func Test_advStat(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "3786607")
	params.Set("mode", "sp")
	s := _host + "/x/dm/adv/state?" + params.Encode()
	body, err := cget(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}

// Test_advPass list
func Test_advList(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "27515260")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/adv/list?" + params.Encode()
	body, err := oget(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}

// Test_advPass list
func Test_advPass(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "13444")
	params.Set("id", "8")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/adv/pass?" + params.Encode()
	body, err := opost(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}

// Test_advDeny list
func Test_advDeny(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "13444")
	params.Set("id", "8")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/adv/deny?" + params.Encode()
	body, err := opost(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}

// Test_advCancel list
func Test_advCancel(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "13444")
	params.Set("id", "8")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/adv/cancel?" + params.Encode()
	body, err := opost(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}
