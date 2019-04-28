package auth

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/net/netutil/breaker"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"

	"github.com/stretchr/testify/assert"
)

const (
	_testUID = "2231365"
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

func create() *Auth {
	return New(&Config{
		Identify:    &warden.ClientConfig{},
		DisableCSRF: false,
	})
}

func engine() *bm.Engine {
	e := bm.NewServer(nil)
	authn := create()
	e.GET("/user", authn.User, func(ctx *bm.Context) {
		mid, _ := ctx.Get("mid")
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	})
	e.GET("/metadata/user", authn.User, func(ctx *bm.Context) {
		mid := metadata.Value(ctx, metadata.Mid)
		ctx.JSON(fmt.Sprintf("%d", mid.(int64)), nil)
	})
	e.GET("/mobile", authn.UserMobile, func(ctx *bm.Context) {
		mid, _ := ctx.Get("mid")
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	})
	e.GET("/metadata/mobile", authn.UserMobile, func(ctx *bm.Context) {
		mid := metadata.Value(ctx, metadata.Mid)
		ctx.JSON(fmt.Sprintf("%d", mid.(int64)), nil)
	})
	e.GET("/web", authn.UserWeb, func(ctx *bm.Context) {
		mid, _ := ctx.Get("mid")
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	})
	e.GET("/guest", authn.Guest, func(ctx *bm.Context) {
		var (
			mid int64
		)
		if _mid, ok := ctx.Get("mid"); ok {
			mid, _ = _mid.(int64)
		}
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	})
	e.GET("/guest/web", authn.GuestWeb, func(ctx *bm.Context) {
		var (
			mid int64
		)
		if _mid, ok := ctx.Get("mid"); ok {
			mid, _ = _mid.(int64)
		}
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	})
	e.GET("/guest/mobile", authn.GuestMobile, func(ctx *bm.Context) {
		var (
			mid int64
		)
		if _mid, ok := ctx.Get("mid"); ok {
			mid, _ = _mid.(int64)
		}
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	})
	e.POST("/guest/csrf", authn.Guest, func(ctx *bm.Context) {
		var (
			mid int64
		)
		if _mid, ok := ctx.Get("mid"); ok {
			mid, _ = _mid.(int64)
		}
		ctx.JSON(fmt.Sprintf("%d", mid), nil)
	})
	return e
}

func TestFromNilConfig(t *testing.T) {
	New(nil)
}

func TestIdentifyHandler(t *testing.T) {
	e := engine()
	go e.Run(":18080")

	time.Sleep(time.Second)

	// test cases
	testWebUser(t, "/user")
	testWebUser(t, "/metadata/user")
	testWebUser(t, "/web")
	testWebUser(t, "/guest")
	testWebUser(t, "/guest/web")
	testWebUserFailed(t, "/user")
	testWebUserFailed(t, "/web")

	testMobileUser(t, "/user")
	testMobileUser(t, "/mobile")
	testMobileUser(t, "/metadata/mobile")
	testMobileUser(t, "/guest")
	testMobileUser(t, "/guest/mobile")
	testMobileUserFailed(t, "/user")
	testMobileUserFailed(t, "/mobile")

	testGuest(t, "/guest")
	testGuestCSRF(t, "/guest/csrf")
	testGuestCSRFFailed(t, "/guest/csrf")
	testMultipartCSRF(t, "/guest/csrf")

	if err := e.Server().Shutdown(context.TODO()); err != nil {
		t.Logf("Failed to shutdown bm engine: %v", err)
	}
}

func testWebUser(t *testing.T, path string) {
	res := Response{}
	query := url.Values{}
	cli := client()
	req, err := cli.NewRequest(http.MethodGet, "http://127.0.0.1:18080/"+path, "", query)
	assert.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID",
		Value: _testUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID__ckMd5",
		Value: "36976f7a5cb6e4a6",
	})
	req.AddCookie(&http.Cookie{
		Name:  "SESSDATA",
		Value: "7bf20cf0%2C1540627371%2C8ec39f0c",
	})

	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)

	assert.Equal(t, 0, res.Code)
	assert.Equal(t, _testUID, res.Data)
}

func testMobileUser(t *testing.T, path string) {
	res := Response{}
	query := url.Values{}
	query.Set("access_key", "cdbd166be6673a5a4f6fbcdd88569edf")
	cli := client()
	req, err := cli.NewRequest(http.MethodGet, "http://127.0.0.1:18080"+path, "", query)
	assert.NoError(t, err)

	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)

	assert.Equal(t, 0, res.Code)
	assert.Equal(t, _testUID, res.Data)
}

func testWebUserFailed(t *testing.T, path string) {
	res := Response{}
	query := url.Values{}
	cli := client()
	req, err := cli.NewRequest(http.MethodGet, "http://127.0.0.1:18080/"+path, "", query)
	assert.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID",
		Value: _testUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID__ckMd5",
		Value: "53c4b106fb4462f1",
	})
	req.AddCookie(&http.Cookie{
		Name:  "SESSDATA",
		Value: "6eeda532%2C1515837495%2C5a6baa4e",
	})

	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)

	assert.Equal(t, ecode.NoLogin.Code(), res.Code)
	assert.Empty(t, res.Data)
}

func testMobileUserFailed(t *testing.T, path string) {
	res := Response{}
	query := url.Values{}
	query.Set("access_key", "5dce488c2ff8d62d7b131da40ae18729")
	cli := client()
	req, err := cli.NewRequest(http.MethodGet, "http://127.0.0.1:18080"+path, "", query)
	assert.NoError(t, err)

	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)

	assert.Equal(t, ecode.NoLogin.Code(), res.Code)
	assert.Empty(t, res.Data)
}

func testGuest(t *testing.T, path string) {
	res := Response{}
	cli := client()
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:18080"+path, nil)
	assert.NoError(t, err)

	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)

	assert.Equal(t, 0, res.Code)
	assert.Equal(t, "0", res.Data)
}

func testGuestCSRF(t *testing.T, path string) {
	res := Response{}
	param := url.Values{}
	param.Set("csrf", "c1524bbf3aa5a1996ff7b1f29a09e796")
	cli := client()
	req, err := cli.NewRequest(http.MethodPost, "http://127.0.0.1:18080"+path, "", param)
	assert.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID",
		Value: _testUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID__ckMd5",
		Value: "36976f7a5cb6e4a6",
	})
	req.AddCookie(&http.Cookie{
		Name:  "SESSDATA",
		Value: "7bf20cf0%2C1540627371%2C8ec39f0c",
	})
	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)

	assert.Equal(t, 0, res.Code)
	assert.Equal(t, _testUID, res.Data)
}

func testGuestCSRFFailed(t *testing.T, path string) {
	res := Response{}
	param := url.Values{}
	param.Set("csrf", "invalid-csrf-token")
	cli := client()
	req, err := cli.NewRequest(http.MethodPost, "http://127.0.0.1:18080"+path, "", param)
	assert.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID",
		Value: _testUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID__ckMd5",
		Value: "36976f7a5cb6e4a6",
	})
	req.AddCookie(&http.Cookie{
		Name:  "SESSDATA",
		Value: "7bf20cf0%2C1540627371%2C8ec39f0c",
	})
	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)

	assert.Equal(t, ecode.CsrfNotMatchErr.Code(), res.Code)
	assert.Empty(t, res.Data)
}

func testMultipartCSRF(t *testing.T, path string) {
	res := Response{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("csrf", "c1524bbf3aa5a1996ff7b1f29a09e796")
	writer.Close()
	req, err := http.NewRequest("POST", "http://127.0.0.1:18080"+path, body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	cli := client()
	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID",
		Value: _testUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "DedeUserID__ckMd5",
		Value: "36976f7a5cb6e4a6",
	})
	req.AddCookie(&http.Cookie{
		Name:  "SESSDATA",
		Value: "7bf20cf0%2C1540627371%2C8ec39f0c",
	})
	err = cli.Do(context.TODO(), req, &res)
	assert.NoError(t, err)
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, _testUID, res.Data)
}
