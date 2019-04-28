package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func init() {
	log.Init(&log.Config{
		Stdout: true,
	})
}

func verify() *Verify {
	return New(&Config{
		OpenServiceHost: "http://uat-open.bilibili.co",
		HTTPClient: &bm.ClientConfig{
			App: &bm.App{
				Key:    "53e2fa226f5ad348",
				Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
			},
			Dial:      xtime.Duration(time.Second),
			Timeout:   xtime.Duration(time.Second),
			KeepAlive: xtime.Duration(time.Second * 10),
			Breaker: &breaker.Config{
				Window:  xtime.Duration(time.Second),
				Sleep:   xtime.Duration(time.Millisecond * 100),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		},
	})
}

func client() *bm.Client {
	return bm.NewClient(&bm.ClientConfig{
		App: &bm.App{
			Key:    "53e2fa226f5ad348",
			Secret: "3cf6bd1b0ff671021da5f424fea4b04a",
		},
		Dial:      xtime.Duration(time.Second),
		Timeout:   xtime.Duration(time.Second),
		KeepAlive: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(time.Second),
			Sleep:   xtime.Duration(time.Millisecond * 100),
			Bucket:  10,
			Ratio:   0.5,
			Request: 100,
		},
	})
}

func engine() *bm.Engine {
	e := bm.New()
	idt := verify()
	e.GET("/verify", idt.Verify, func(c *bm.Context) {
		c.JSON("pass", nil)
	})
	e.GET("/verifyUser", idt.VerifyUser, func(c *bm.Context) {
		mid := metadata.Int64(c, metadata.Mid)
		fmt.Println(mid)
		c.JSON(fmt.Sprintf("%d", mid), nil)
	})
	return e
}

func TestNewWithNilConfig(t *testing.T) {
	New(nil)
}

func TestVerifyIdentifyHandler(t *testing.T) {
	e := engine()
	go e.Run(":18080")

	time.Sleep(time.Second)

	// test cases
	testVerifyFailed(t)
	testVerifySuccess(t)
	testVerifyUser(t)
	testVerifyUserFailed(t)
	testVerifyUserInvalid(t)

	if err := e.Server().Shutdown(context.TODO()); err != nil {
		t.Logf("Failed to shutdown bm engine: %v", err)
	}
}

func httpGet(url string) (code int, content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	code = resp.StatusCode
	return
}

func testVerifyFailed(t *testing.T) {
	res := Response{}
	code, content, err := httpGet("http://127.0.0.1:18080/verify?ts=1&appkey=53e2fa226f5ad348")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	err = json.Unmarshal(content, &res)
	assert.NoError(t, err)
	assert.Equal(t, -3, res.Code)
}

func testVerifySuccess(t *testing.T) {
	res := Response{}
	uv := url.Values{}
	uv.Set("appkey", "53e2fa226f5ad348")
	err := client().Get(context.TODO(), "http://127.0.0.1:18080/verify", "", uv, &res)
	assert.NoError(t, err)
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, "pass", res.Data)
}

func testVerifyUser(t *testing.T) {
	res := Response{}
	query := url.Values{}
	query.Set("mid", "1")
	query.Set("appkey", "53e2fa226f5ad348")
	err := client().Get(context.TODO(), "http://127.0.0.1:18080/verifyUser", "", query, &res)
	assert.NoError(t, err)
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, "1", res.Data)
}

func testVerifyUserFailed(t *testing.T) {
	res := Response{}
	code, content, err := httpGet("http://127.0.0.1:18080/verifyUser?ts=1&appkey=53e2fa226f5ad348")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	err = json.Unmarshal(content, &res)
	assert.NoError(t, err)
	assert.Equal(t, -3, res.Code)
}

func testVerifyUserInvalid(t *testing.T) {
	res := Response{}
	query := url.Values{}
	query.Set("mid", "aaaa/")
	query.Set("appkey", "53e2fa226f5ad348")
	err := client().Get(context.TODO(), "http://127.0.0.1:18080/verifyUser", "", query, &res)
	assert.NoError(t, err)
	assert.Equal(t, -400, res.Code)
	assert.Equal(t, "", res.Data)
}
