package http

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"go-common/app/admin/main/dm/conf"
)

const (
	localHost   = "http://127.0.0.1:6311"
	host        = "http://uat-api.bilibili.co"
	managerHost = "http://uat-manager.bilibili.co"
)

var (
	_dmList         = host + "/x/internal/dmadmin/content/list"
	_refresh        = host + "/x/internal/dmadmin/content/refresh"
	_editState      = host + "/x/internal/dmadmin/content/edit/state"
	_editPool       = host + "/x/internal/dmadmin/content/edit/pool"
	_editAttr       = host + "/x/internal/dmadmin/content/edit/attr"
	_monitorList    = host + "/x/internal/dmadmin/monitor/list"
	_editMonitor    = host + "/x/internal/dmadmin/monitor/edit"
	_dmReportStat   = host + "/x/internal/dmadmin/report/stat/change"
	_dmReportList   = host + "/x/internal/dmadmin/report/list"
	_dmReportdetail = host + "/x/internal/dmadmin/report"
	_archiveList    = localHost + "/x/admin/dm/subject/archive"
	_subjectLog     = localHost + "/x/admin/dm/subject/log"
	_uptSubState    = localHost + "/x/admin/dm/subject/state/edit"
	_dmSendJudge    = host + "/x/internal/dmadmin/report/judge"
	_dmJudgeResult  = host + "/x/internal/dmadmin/report/judge/result"
	_dmIndexInfoURL = host + "/x/internal/dmadmin/content/index/info"
	_subjectLimit   = host + "/x/admin/dm/subject/maxlimit"
	_maskUps        = localHost + "/x/admin/dm/mask/up"
	_maskUpOpen     = localHost + "/x/admin/dm/mask/up/open"
	_taskList       = managerHost + "/x/admin/dm/task/list"
	_taskNew        = managerHost + "/x/admin/dm/task/new"
	_taskReview     = localHost + "/x/admin/dm/task/review"
	_taskView       = localHost + "/x/admin/dm/task/view"
)

func TestMain(m *testing.M) {
	var err error
	conf.ConfPath = "../cmd/dm-admin-test.toml"
	if err = conf.Init(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// TestDMList list
func TestDMList(t *testing.T) {
	params := url.Values{}
	params.Set("type", "1")
	params.Set("oid", "10108994")
	params.Set("keyword", "哈哈")
	params.Set("page_size", "1")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := _dmList + "?" + params.Encode()
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

// TestRefrash refrash
func TestRefresh(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "8937277")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	url := _refresh + "?" + params.Encode()
	t.Log(url)
	body, err := opost(url)
	if err != nil {
		t.Errorf("url(%s) error(%s)", url, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

// TestEditState batch
func TestEditState(t *testing.T) {
	params := url.Values{}
	params.Set("type", "1")
	params.Set("oid", "1222")
	params.Set("state", "0")
	params.Set("dmids", "719181990,719181991")
	params.Set("moral", "1")
	params.Set("uname", "test")
	params.Set("remark", "测试下")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	url := _editState + "?" + params.Encode()
	body, err := opost(url)
	if err != nil {
		t.Errorf("url(%s) error(%s)", url, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

// TestEditPool batch
func TestEditPool(t *testing.T) {
	params := url.Values{}
	params.Set("type", "1")
	params.Set("oid", "1222")
	params.Set("pool", "0")
	params.Set("dmids", "719181990,719181991")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	url := _editPool + "?" + params.Encode()
	fmt.Println(url)
	body, err := opost(url)
	if err != nil {
		t.Errorf("url(%s) error(%s)", url, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

// TestEditAttr batch
func TestEditAttr(t *testing.T) {
	params := url.Values{}
	params.Set("type", "1")
	params.Set("oid", "1222")
	params.Set("attr", "1")
	params.Set("dmids", "719181990,719181991")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	url := _editAttr + "?" + params.Encode()
	fmt.Println(url)
	body, err := opost(url)
	if err != nil {
		t.Errorf("url(%s) error(%s)", url, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestArchiveList(t *testing.T) {
	p := url.Values{}
	//	p.Set("type", "oid")
	//	p.Set("id", "1345")
	p.Set("sort", "desc")
	p.Set("order", "mtime")
	requestURL := _archiveList + "?" + p.Encode()
	fmt.Println(requestURL)
	body, err := oget(requestURL)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestSubjectLog(t *testing.T) {
	p := url.Values{}
	p.Set("type", "1")
	p.Set("oid", "1221")
	requestURL := _subjectLog + "?" + p.Encode()
	fmt.Println(requestURL)
	body, err := oget(requestURL)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestUptSubjectsState(t *testing.T) {
	p := url.Values{}
	p.Set("oids", "1223")
	p.Set("state", "1")
	requestURL := _uptSubState + "?" + p.Encode()
	fmt.Println(requestURL)
	body, err := opost(requestURL)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestReportStatChange(t *testing.T) {
	params := url.Values{}
	params.Set("ids", "9968618:719923090,719923092,719923093")
	params.Set("state", "2")
	params.Set("uid", "150781")
	params.Set("reason", "1")
	params.Set("state", "3")
	params.Set("remark", "test")
	params.Set("notice", "3")
	params.Set("block", "-1")
	params.Set("moral", "10")
	params.Set("block_reason", "5")
	params.Set("uname", "zhanghongwen")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _dmReportStat + "?" + params.Encode()
	body, err := opost(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestReportList(t *testing.T) {
	params := url.Values{}
	params.Set("state", "0,1")
	params.Set("up_op", "0")
	params.Set("page", "1")
	params.Set("tid", "24")
	params.Set("aid", "")
	params.Set("uid", "")
	params.Set("rp_user", "")
	params.Set("rp_type", "")
	params.Set("start", url.QueryEscape("2017-03-21 00:00:00"))
	params.Set("end", url.QueryEscape("2017-05-29 00:00:00"))
	params.Set("order", "desc")
	params.Set("sort", "rp_time")
	params.Set("page_size", "2")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _dmReportList + "?" + params.Encode()
	t.Log(requestURL)
	body, err := oget(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestReportDetail(t *testing.T) {
	params := url.Values{}
	params.Set("dmid", "719254592")
	params.Set("cid", "10106598")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _dmReportdetail + "?" + params.Encode()
	body, err := oget(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestSendJudge(t *testing.T) {
	params := url.Values{}
	params.Set("ids", "10108441:719925897")
	params.Set("uname", "luoxiaofan")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _dmSendJudge + "?" + params.Encode()
	t.Log(requestURL)
	fmt.Println(requestURL)
	body, err := opost(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestJudgeResult(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "1")
	params.Set("dmid", "2")
	params.Set("result", "1")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _dmJudgeResult + "?" + params.Encode()
	fmt.Println(requestURL)
	body, err := opost(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestSubjectLimit(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "1221")
	params.Set("limit", "99")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _subjectLimit + "?" + params.Encode()
	fmt.Println(requestURL)
	body, err := opost(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestDMIndexInfo(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "9967205")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _dmIndexInfoURL + "?" + params.Encode()
	fmt.Println(requestURL)
	body, err := oget(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestAddTrJob(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "9967205")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
}

func TestMonitorList(t *testing.T) {
	params := url.Values{}
	params.Set("page", "1")
	params.Set("page_size", "5")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _monitorList + "?" + params.Encode()
	t.Log(requestURL)
	body, err := oget(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestEditMonitor(t *testing.T) {
	params := url.Values{}
	params.Set("type", "1")
	params.Set("oids", "1221,1222")
	params.Set("state", "1")
	params.Set("appkey", "f6433799dbd88751")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "36f8ddb1806207fe07013ab6a77a3935"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	requestURL := _editMonitor + "?" + params.Encode()
	t.Log(requestURL)
	body, err := opost(requestURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", requestURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestTransferList(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "10108682")
	params.Set("state", "1")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := managerHost + "/x/admin/dm/transfer/list?" + params.Encode()
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
	//t.Logf("===========%v", body)
}

func TestReTransferJob(t *testing.T) {
	params := url.Values{}
	params.Set("id", "265")
	params.Set("mid", "1")
	s := managerHost + "/x/admin/dm/transfer/retry?" + params.Encode()
	t.Log(s)
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

func TestMaskUps(t *testing.T) {
	body, err := oget(_maskUps)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", _maskUps, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestMaskUpOpen(t *testing.T) {
	params := url.Values{}
	params.Set("mids", "11111")
	params.Set("state", "1")
	//	params.Set("comment", "test")
	requestURL := _maskUpOpen + "?" + params.Encode()
	fmt.Println(requestURL)
	body, err := opost(requestURL)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", _maskUps, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestTaskList(t *testing.T) {
	body, err := oget(_taskList)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", _taskList, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestTaskNew(t *testing.T) {
	params := url.Values{}
	params.Set("title", "test")
	params.Set("mids", "1111")
	params.Set("start", "2018-11-11 00:00:00")
	params.Set("end", "2018-11-12 00:00:00")
	requestURL := _taskNew + "?" + params.Encode()
	body, err := cpost(requestURL)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", _taskNew, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestTaskReview(t *testing.T) {
	params := url.Values{}
	params.Set("id", "15")
	params.Set("state", "1")
	requestURL := _taskReview + "?" + params.Encode()
	body, err := opost(requestURL)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", _taskNew, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
}

func TestTaskView(t *testing.T) {
	params := url.Values{}
	params.Set("id", "14")
	requestURL := _taskView + "?" + params.Encode()
	body, err := oget(requestURL)
	fmt.Println(string(body))
	if err != nil {
		t.Errorf("url(%s) error(%s)", _taskNew, err)
		t.FailNow()
	}
	var out bytes.Buffer
	if err = json.Indent(&out, body, "", " "); err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	fmt.Println(out.String())
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
func ccode(body []byte) (err error) {
	var d interface{}
	err = json.Unmarshal(body, d)
	if err != nil {
		return
	}
	ma, ok := d.(map[string]interface{})
	if !ok {
		return
	}
	code := ma["code"].(int)
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
	req.Header.Add("Cookie", `uid=1130; username=fengduzhen; mng-bilibili=b4nec4bi6un0l9ftkfsgj59uq1; _AJSESSIONID=0935506167be9a62bd477a9f10bd6021; JSESSIONID=940FB2C9CAE4FD14171A5D632689C6E5; _uuid=FCC1090C-7E26-52E5-2DC6-26D8EE1B830506546infoc; mng-go=d56eabfa9edf138604586d5dee4929273d1f76a00cfcd69278b6264c453bae2e`)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}
