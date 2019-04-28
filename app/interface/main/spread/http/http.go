package http

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"

	"go-common/app/interface/main/spread/conf"
	"go-common/app/interface/main/spread/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

var (
	svc      *service.Service
	business = make(map[string]string)
)

// Init init
func Init(c *conf.Config) {
	initService(c)
	// init router
	engine := bm.DefaultServer(c.BM)
	outerRouter(engine)
	if err := engine.Start(); err != nil {
		log.Error("xhttp.Serve error(%v)", err)
		panic(err)
	}
	for _, b := range c.Businesses {
		business[b.Appkey] = b.AppSecret
	}
}

// initService init services.
func initService(c *conf.Config) {
	svc = service.New(c)
}

// outerRouter init outer router api path.
func outerRouter(e *bm.Engine) {
	//init api
	e.Ping(ping)
	group := e.Group("/x/spread", Verify)
	{
		bangumi := group.Group("/bangumi")
		{
			bangumi.GET("/content", bangumiContent)
			bangumi.GET("/offshelve", bangumiOff)
		}
	}
}

// ping check server ok.
func ping(c *bm.Context) {
}

// Verify will inject into handler func as verify required
func Verify(ctx *bm.Context) {
	var (
		secret string
		ok     bool
	)
	req := ctx.Request
	params := req.Form

	if req.Method == "POST" {
		// Give priority to sign in url query, otherwise check sign in post form.
		q := req.URL.Query()
		if q.Get("sign") != "" {
			params = q
		}
	}

	// check timestamp is not empty (TODO : Check if out of some seconds.., like 100s)
	if params.Get("ts") == "" {
		log.Error("ts is empty")
		ctx.JSON(nil, ecode.RequestErr)
		ctx.Abort()
		return
	}

	sign := params.Get("sign")
	params.Del("sign")
	sappkey := params.Get("appkey")
	secret, ok = business[sappkey]
	if !ok {
		ctx.JSON(nil, ecode.AppKeyInvalid)
		ctx.Abort()
		return
	}

	if hsign := paramsSign(params, sappkey, secret, true); hsign != sign {
		if hsign1 := paramsSign(params, sappkey, secret, false); hsign1 != sign {
			log.Error("Get sign: %s, expect %s", sign, hsign)
			ctx.JSON(nil, ecode.SignCheckErr)
			ctx.Abort()
			return
		}
	}
}

// Render .
func Render(c *bm.Context, code int, msg string, data interface{}, total int64, err error) {
	c.Error = err
	bcode := ecode.Cause(err)
	if err != nil {
		c.Render(http.StatusOK, render.JSON{
			Code:    bcode.Code(),
			Message: bcode.Message(),
			Data:    data,
		})
		return
	}
	c.Render(http.StatusOK, render.MapJSON{
		"code":    code,
		"message": msg,
		"data":    data,
		"total":   total,
	})
}

// sign is used to sign form params by given condition.
func paramsSign(params url.Values, appkey string, secret string, lower bool) (hexdigest string) {
	data := params.Encode()
	if strings.IndexByte(data, '+') > -1 {
		data = strings.Replace(data, "+", "%20", -1)
	}
	if lower {
		data = strings.ToLower(data)
	}
	digest := md5.Sum([]byte(data + secret))
	hexdigest = hex.EncodeToString(digest[:])
	return
}
