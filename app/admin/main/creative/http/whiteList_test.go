package http

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"go-common/app/admin/main/creative/conf"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	_host = "http://0.0.0.0:6344"
)

var (
	_view = _host + "/x/admin/creative/whitelist/view"
)

func init() {
	dir, _ := filepath.Abs("../cmd/creative-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
}

// Sign fn
func Sign(params url.Values) (query string, err error) {
	if len(params) == 0 {
		return
	}
	if params.Get("appkey") == "" {
		err = fmt.Errorf("utils http get must have parameter appkey")
		return
	}
	if params.Get("appsecret") == "" {
		err = fmt.Errorf("utils http get must have parameter appsecret")
		return
	}
	if params.Get("sign") != "" {
		err = fmt.Errorf("utils http get must have not parameter sign")
		return
	}
	// sign
	secret := params.Get("appsecret")
	params.Del("appsecret")
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp + secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	query = params.Encode()
	return
}

func Test_View(t *testing.T) {
	Convey("View", t, func() {
		params := url.Values{}
		params.Set("id", "7")
		params.Set("appkey", conf.Conf.App.Key)
		params.Set("appsecret", conf.Conf.App.Secret)
		params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
		var (
			query, _ = Sign(params)
			url      string
		)
		url = _view + "?" + query
		body, err := oget(url)
		So(err, ShouldBeNil)
		So(body, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
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
