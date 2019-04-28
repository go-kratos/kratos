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
	"strings"
	"testing"
	"time"

	"go-common/app/interface/main/dm2/model"
)

const (
	_host  = "http://127.0.0.1:6701"
	_host2 = "http://api.bilibili.com"
	_host3 = "http://uat-api.bilibili.co"
)

var (
	cookie = `sid=88dol9eo; fts=1508407962; LIVE_BUVID=AUTO4115084079913953; UM_distinctid=15f3436b684ab-08804fe751374a-31657c00-1fa400-15f3436b685298; pgv_pvi=3618563072; rpdid=kmilkmximpdoswqploqxw; biliMzIsnew=1; biliMzTs=null; buvid3=CBB97852-6CF6-4D3A-B97F-5A9AD1D5827F26561infoc; finger=14bc3c4e; pgv_si=s6534896640; purl_token=bilibili_1511162275; DedeUserID=3078992; DedeUserID__ckMd5=55845496fd6119b5; SESSDATA=7bf20cf0%2C1513762568%2C258e3d17; bili_jct=160db0def61324719de104b7bfe1e3ef; _cnt_pm=0; _cnt_notify=0; _dfcaptcha=b5463dc96d22c22164a8c0e5a1ba9e96; member_v2=1`

	dmLikeActURL     = _host2 + "/x/v2/dm/thumbup/add"
	dmLikeListURL    = _host2 + "/x/v2/dm/thumbup/stats"
	dmEditURL        = _host + "/x/internal/v2/dm/edit/state"
	dmPoolURL        = _host + "/x/internal/v2/dm/edit/pool"
	dmUpSearchURL    = _host3 + "/x/v2/dm/search"
	dmRecentURL      = _host + "/x/internal/v2/dm/recent"
	updateMaskURL    = _host3 + "/x/internal/v2/dm/mask/update"
	editUpFiltersURL = _host + "/x/v2/dm/filter/up/edit"
)

func TestEditUpFilter(t *testing.T) {
	fliters := model.UpFilter{
		Filter:  "zzzzzz",
		Comment: "",
	}
	p := url.Values{}
	p.Set("active", "1")
	p.Set("type", "0")
	filter, err := json.Marshal(fliters)
	if err != nil {
		return
	}
	fmt.Print(string(filter))
	p.Set("filters", string(filter))
	mh := md5.Sum([]byte(p.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	p.Set("sign", hex.EncodeToString(mh[:]))
	reqURL := editUpFiltersURL + "?" + p.Encode()
	fmt.Println(reqURL)
	body, err := opost(reqURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", reqURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}
func TestRecent(t *testing.T) {
	params := url.Values{}
	params.Set("uid", "27515256")
	params.Set("page", "5")
	//params.Set("mid", "27515313")
	params.Set("order", "ctime")
	params.Set("size", "1")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := dmRecentURL + "?" + params.Encode()
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

func TestDMEdit(t *testing.T) {
	params := url.Values{}
	params.Set("oid", "5")
	params.Set("type", "1")
	params.Set("mid", "27515313")
	params.Set("dmids", "719149463")
	params.Set("state", "0")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := dmEditURL + "?" + params.Encode()
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

func TestDMPool(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "5")
	params.Set("mid", "27515313")
	params.Set("dmids", "719149463")
	params.Set("pool", "2")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := dmPoolURL + "?" + params.Encode()
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

func TestLikeDM(t *testing.T) {
	params := url.Values{}
	params.Set("oid", "5")
	params.Set("dmid", "719149462")
	params.Set("op", "1")
	reqURL := dmLikeActURL + "?" + params.Encode()
	body, err := cpost(reqURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", reqURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}

func TestUpdateMask(t *testing.T) {
	params := url.Values{}
	params.Set("cid", "10109227")
	params.Set("time", "30")
	params.Set("fps", "20")
	params.Set("plat", "1")
	params.Set("count", "20")
	params.Set("list", "26777486_s0_0_1499")
	params.Set("appkey", "8f62754d8d594e90")
	params.Set("ts", strconv.FormatInt(1527564095, 10))
	data := params.Encode()
	if strings.IndexByte(data, '+') > -1 {
		data = strings.Replace(data, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(data + "test"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	s := updateMaskURL + "?" + params.Encode()
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

func TestLikeList(t *testing.T) {
	params := url.Values{}
	params.Set("oid", "27139273")
	params.Set("ids", "719149462,719149463")
	reqURL := dmLikeListURL + "?" + params.Encode()
	fmt.Println(reqURL)
	body, err := cget(reqURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", reqURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	fmt.Println(string(body))
}

func TestDMUpSearch(t *testing.T) {
	params := url.Values{}
	params.Set("oid", "10131156")
	params.Set("type", "1")
	params.Set("appkey", "53e2fa226f5ad348")
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + "3cf6bd1b0ff671021da5f424fea4b04a"))
	params.Set("sign", hex.EncodeToString(mh[:]))
	reqURL := dmUpSearchURL + "?" + params.Encode()
	fmt.Println(reqURL)
	body, err := oget(reqURL)
	if err != nil {
		t.Errorf("url(%s) error(%s)", reqURL, err)
		t.FailNow()
	}
	if err = ccode(body); err != nil {
		t.Fatal(err, string(body))
		t.FailNow()
	}
	t.Logf("======%+v", body)
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
	req.Header.Add("Cookie", cookie)
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
	req.Header.Add("Cookie", cookie)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}
