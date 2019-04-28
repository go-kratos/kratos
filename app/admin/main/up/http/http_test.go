package http

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-common/app/admin/main/up/conf"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	dir, _ := filepath.Abs("../cmd/up-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	// Init(conf.Conf)
	time.Sleep(time.Second)
}

// Sign fn
func Sign(params url.Values) (sign string) {
	secret := params.Get("appsecret")
	params.Del("appsecret")
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp + secret))
	sign = hex.EncodeToString(mh[:])
	return
}

var (
	err     error
	req     *http.Request
	resp    *http.Response
	HOST    = "http://localhost:7441"
	URI     = "/x/internal/up/register"
	infoURI = "/x/internal/up/info"
	c       = context.Background()
	client  = &http.Client{
		Timeout: time.Duration(time.Second * 2),
	}
)

func Test_Up(t *testing.T) {
	Convey("register", t, func() {
		params := url.Values{}
		params.Set("mid", strconv.FormatInt(2089809, 10))
		params.Set("from", strconv.FormatInt(0, 10))
		params.Set("is_author", strconv.FormatInt(0, 10))
		params.Set("appkey", conf.Conf.App.Key)
		params.Set("appsecret", conf.Conf.App.Secret)
		params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
		sign := Sign(params)
		params.Set("sign", sign)
		u, _ := url.ParseRequestURI(HOST)
		u.Path = URI
		url := u.String()

		req, err = http.NewRequest("POST", url, strings.NewReader(params.Encode()))
		So(err, ShouldBeNil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// timeout
		ctx, cancel := context.WithTimeout(c, time.Second*2)
		req = req.WithContext(ctx)
		defer cancel()
		resp, err = client.Do(req)
		So(err, ShouldBeNil)

		body, err1 := ioutil.ReadAll(resp.Body)
		err = err1
		So(err, ShouldBeNil)
		defer resp.Body.Close()
		var result struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				Result bool `json:"result"`
			} `json:"data"`
		}
		spew.Dump(string(body))
		json.Unmarshal(body, &result)
		So(result, ShouldNotBeNil)
		So(result.Data.Result, ShouldBeTrue)
	})
	Convey("info", t, func() {
		params := url.Values{}
		params.Set("mid", strconv.FormatInt(2089809, 10))
		params.Set("from", strconv.FormatInt(1, 10))
		params.Set("appkey", conf.Conf.App.Key)
		params.Set("appsecret", conf.Conf.App.Secret)
		params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
		sign := Sign(params)
		params.Set("sign", sign)
		u, _ := url.ParseRequestURI(HOST)
		u.Path = infoURI
		url := u.String()
		reqURL := url + "?" + params.Encode()
		req, err = http.NewRequest("GET", reqURL, nil)
		So(err, ShouldBeNil)
		// timeout
		ctx, cancel := context.WithTimeout(c, time.Second*2)
		req = req.WithContext(ctx)
		defer cancel()
		resp, err = client.Do(req)
		So(err, ShouldBeNil)

		body, err := ioutil.ReadAll(resp.Body)
		So(err, ShouldBeNil)
		defer resp.Body.Close()
		var result struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				IsAuthor bool `json:"is_author"`
			} `json:"data"`
		}
		spew.Dump(string(body))
		json.Unmarshal(body, &result)
		So(result, ShouldNotBeNil)
	})
}
