package http

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

const (
	_host      = "http://127.0.0.1:6214"
	_innerHost = "http://127.0.0.1:6214"
)

var (
	apply             = _host + "/x/dm/protect/apply"
	recallURL         = _host + "/x/dm/recall"
	my                = _host + "/x/dm/my/listxml.so"
	applyList         = _innerHost + "/x/dm/up/protect/apply/list"
	videoList         = _innerHost + "/x/dm/up/protect/apply/video/list"
	applyChangeStatus = _innerHost + "/x/dm/up/protect/apply/status"
	applyNoticeSwitch = _innerHost + "/x/dm/up/protect/apply/notice/switch"
)

func TestUpFilter(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "1")
	params.Set("oid", "0")

	s := _host + "/x/dm/filter/up?" + params.Encode()
	body, err := cget(s)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
}

func TestAddUpFilter(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "1")
	params.Set("oid", "2")
	params.Set("filter", "aabbccdd1234567")
	params.Set("type", "2")
	s := _host + "/x/dm/filter/user/add?" + params.Encode()
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

func TestDelUpFilter1(t *testing.T) {
	params := url.Values{}
	params.Set("ids", "1554")
	s := _host + "/x/dm/filter/user/del?" + params.Encode()
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

func TestAddFilter2(t *testing.T) {
	params := url.Values{}
	params.Set("filters", "1233432,78439214")
	params.Set("type", "2")
	params.Set("comment", "reason")
	s := _host + "/x/dm/filter/user/add2?" + params.Encode()
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

func TestAddReport(t *testing.T) {
	params := url.Values{}
	params.Set("dmid", "719923250")
	params.Set("cid", "8937277")
	params.Set("reason", "1")
	s := _host + "/x/dm/report/add?" + params.Encode()
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

func TestAddReport2(t *testing.T) {
	params := url.Values{}
	params.Set("dmids", "719189659,719189660")
	params.Set("cid", "8937185")
	params.Set("reason", "1")
	s := _host + "/x/dm/report/add2?" + params.Encode()
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

// TestList list
func TestApply(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937277")
	params.Set("dmids", "719214016,123")
	s := apply + "?" + params.Encode()
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

func TestRecall(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937185")
	params.Set("dmid", "719189659")
	s := recallURL + "?" + params.Encode()
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

// TestList list
func TestMy(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937277")
	s := my + "?" + params.Encode()
	body, err := cget(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	fmt.Println(string(body))
}

// TestList list
func TestApplyList(t *testing.T) {
	var (
		s    string
		body []byte
		err  error
		mh   [16]byte
	)
	params := url.Values{}
	params.Set("uid", "15555180")
	params.Set("page", "1")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh = md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s = applyList + "?" + params.Encode()
	body, err = oget(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))

	params.Del("sign")
	params.Set("sort", "playtime")
	mh = md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s = applyList + "?" + params.Encode()
	body, err = oget(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))

	params.Del("sign")
	params.Set("aid", "4052732")
	mh = md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s = applyList + "?" + params.Encode()
	body, err = oget(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))

	params.Del("sign")
	params.Set("sort", "playtime")
	mh = md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s = applyList + "?" + params.Encode()
	body, err = oget(s)
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

// TesVideoList list
func TestVideoList(t *testing.T) {
	params := url.Values{}
	params.Set("uid", "15555180")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := videoList + "?" + params.Encode()
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

// TestApplyChangeStatus TestApplyChangeStatus
func TestApplyChangeStatus(t *testing.T) {
	params := url.Values{}
	params.Set("uid", "15555180")
	params.Set("status", "1")
	params.Set("ids", "1,2")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := applyChangeStatus + "?" + params.Encode()
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

// TestApplyNoticeSwitch TestApplyChangeStatus
func TestApplyNoticeSwitch(t *testing.T) {
	var (
		mh   [16]byte
		s    string
		body []byte
		err  error
	)
	params := url.Values{}
	params.Set("uid", "15555180")
	params.Set("status", "0")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh = md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s = applyNoticeSwitch + "?" + params.Encode()
	body, err = opost(s)
	if err != nil {
		t.Errorf("url(%s) error(%s)", s, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))

	params.Del("sign")
	params.Set("status", "1")
	mh = md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s = applyNoticeSwitch + "?" + params.Encode()
	body, err = opost(s)
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

// TestUptStat list
func TestUptStat(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("dmid", "719189706")
	params.Set("stat", "0")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/stat/upt?" + params.Encode()
	fmt.Println(s)
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

// TestUptStat list
func TestDMTransfer(t *testing.T) {
	params := url.Values{}
	params.Set("from", "10108810")
	params.Set("to", "10108809")
	params.Set("mid", "452156")
	params.Set("offset", "0")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/up/transfer2?" + params.Encode()
	fmt.Println(s)
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

// TestAssistBanned list
func TestAssistBanned(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "10108163")
	params.Set("dmids", "719925837,719925859")
	s := _host + "/x/dm/assist/banned?" + params.Encode()
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

// TestAssistDelBanned2 list
func TestAssistDelBanned2(t *testing.T) {
	params := url.Values{}
	params.Set("aid", "4052732")
	params.Set("hashes", "131141213,mdzz")
	s := _host + "/x/dm/up/banned/del?" + params.Encode()
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

// TestUser user
func TestUser(t *testing.T) {
	s := _host + "/x/dm/user"
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

func TestBannedUsers(t *testing.T) {
	s := _host + "/x/dm/up/banned/users"
	params := url.Values{}
	params.Set("aid", "10097375")
	body, err := cget(s + "?" + params.Encode())
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

// TestAssistDel list
func TestAssistDel(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("dmids", "719189706,719189707")
	s := "http://api.bilibili.com/x/dm/assist/del?" + params.Encode()
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

// TestAssistProtect list
func TestAssistProtect(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("dmids", "719189706,719189707")
	s := _host + "/x/dm/assist/protect?" + params.Encode()
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

// TestAssistProtect list
func TestAssistMove(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("dmids", "719189706,719189707")
	s := _host + "/x/dm/assist/move?" + params.Encode()
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

// TestAssistIgnore list
func TestAssistBannedUpt(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "15555180")
	params.Set("hash", "123")
	params.Set("stat", "0")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/assist/banned/upt?" + params.Encode()
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

// TestUptPool
func TestUptPool(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("dmid", "719189708")
	params.Set("pool", "1")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _innerHost + "/x/internal/dm/pool/upt?" + params.Encode()
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

// TestAssistIgnore list
func TestAssistIgnore(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("dmids", "719189706,719189707")
	s := _host + "/x/dm/assist/ignore?" + params.Encode()
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

func TestDMList(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("page", "1")
	s := _host + "/x/dm/list?" + params.Encode()
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

func TestDMedit(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937278")
	params.Set("dmids", "1,2,3,4")
	params.Set("state", "1")

	s := _host + "/x/dm/edit?" + params.Encode()
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

func TestFilterList(t *testing.T) {
	params := url.Values{}
	params.Set("mid", "150781")
	params.Set("cid", "0")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _host + "/x/internal/dm/filter/index/list?" + params.Encode()
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
	fmt.Println(string(body))
}

func TestTransferList(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "10109082")
	params.Set("mid", "1")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _host + "/x/dm/transfer/list?" + params.Encode()
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
	fmt.Println(string(body))
}

func TestTransferRetry(t *testing.T) {
	params := url.Values{}
	params.Set("id", "265")
	params.Set("mid", "1")

	s := _host + "/x/dm/transfer/retry?" + params.Encode()
	fmt.Println(s)
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

func TestEditFilter(t *testing.T) {
	fStr := `[
    {
        "type": 0,
        "filter": "傻逼",
        "activate": 1
    },
    {
        "type": 1,
        "filter": "*",
        "activate": 1
    },
    {
        "type": 4,
        "filter": "",
        "activate":1
    },
    {
        "type": 2,
        "filter": "150781",
        "activate": 0
    }
]`
	params := url.Values{}
	params.Set("mid", "150781")
	params.Set("cid", "0")
	params.Set("filters", fStr)
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))

	s := _host + "/x/internal/dm/filter/index/edit?" + params.Encode()
	fmt.Println(s)
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

func TestUserRules(t *testing.T) {
	s := _host + "/x/dm/filter/user"
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

// oget http get request
func oget(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// ccode check code == 0
// ccode check code == 0
func ccode(body []byte) (err error) {
	var d interface{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		return
	}
	ma, ok := d.(map[string]interface{})
	if !ok {
		return
	}
	code := ma["code"].(float64)
	if code != 0 {
		err = errors.New("code != 0")
		return
	}
	return
}

// opost http post request
func opost(url string) (body []byte, err error) {
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func cpost(url string) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", `buvid3=7B833293-671D-47D9-9BAB-464DE756DFBF14331infoc; UM_distinctid=15ba97a9afed0-0b51074d03d7c4-396b7805-1fa400-15ba97a9aff25b; pgv_pvi=6667228160; pgv_si=s603312128; fts=1493197036; sid=chto560j; rpdid=kwmqqoqxomdoplppwwppw; PHPSESSID=nepo3f6lhgci20gbdet1vvfig1; HTML5PlayerCRC32=3819803418; biliMzIsnew=1; biliMzTs=0; html5_player_gray=false; player_gray=false; LIVE_BUVID=ecbc29c14b41a4c5ff883313e734a336; LIVE_BUVID__ckMd5=f2922ceec4b33165; user_face=http%3A%2F%2Fi2.hdslb.com%2Fbfs%2Fface%2F9c843ff4cdcf38e6031e493ff17b1b8610e36e7a.jpg; purl_token=bilibili_1500257956; finger=0e029071; DedeUserID=27515260; DedeUserID__ckMd5=ae01241c690f95b5; SESSDATA=8bb74c97%2C1502854413%2Ccffcbc21; bili_jct=2b4c5f02e710a7f676acecb80c03271f; _cnt_pm=0; _cnt_notify=0; CNZZDATA2724999=cnzz_eid%3D714476842-1493196394-null%26ntime%3D1500262370; _dfcaptcha=a1ca150412337140d3d763aea1da2d41`)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func cget(url string) (body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", `sid=88dol9eo; fts=1508407962; UM_distinctid=15f3436b684ab-08804fe751374a-31657c00-1fa400-15f3436b685298; pgv_pvi=3618563072; rpdid=kmilkmximpdoswqploqxw; buvid3=CBB97852-6CF6-4D3A-B97F-5A9AD1D5827F26561infoc; biliMzIsnew=1; biliMzTs=null; LIVE_BUVID=aabb755e95239c26d7f2dbeba748ce27; LIVE_BUVID__ckMd5=8b3ae56b45b2626a; member_v2=1; finger=14bc3c4e; BANGUMI_SS_21715_REC=192513; im_notify_type_3078992=0; DedeUserID=3078992; DedeUserID__ckMd5=55845496fd6119b5; SESSDATA=f7d955fd%2C1519879629%2Ca8074050; bili_jct=553a35adf1c94efb8d125d000798d1ca; pgv_si=s6786742272; purl_token=bilibili_1517461841; _dfcaptcha=30e17d94b2de231c8862ff9d16fff5dc`)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}
