package http

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strings"

	"go-common/app/interface/openplatform/seo/conf"
	"go-common/app/interface/openplatform/seo/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	srv *service.Service
	vfy *verify.Verify
)

// Init init seo service
func Init(c *conf.Config) {
	srv = service.New(c)
	vfy = verify.New(c.Verify)
	engine := bm.DefaultServer(c.BM)
	router(engine)
	if err := engine.Start(); err != nil {
		log.Error("innerEngine.Start() error(%v)", err)
		panic(err)
	}
}

func router(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)

	g := e.Group("/platform")
	{
		g.GET("/home.html", proList)
		g.GET("/detail.html", proInfo)
	}
	e.GET("/detail.html", itemInfo)
	e.GET("/sitemap.xml", sitemap)
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

// isBot check if user-agent is bot
func isBot(c *bm.Context) bool {
	ua := c.Request.Header.Get("User-Agent")
	ua = strings.ToLower(ua)
	for _, bot := range conf.Conf.Seo.BotList {
		if strings.Contains(ua, strings.ToLower(bot)) {
			return true
		}
	}
	return false
}

// FullUrl get request full url
func FullUrl(c *bm.Context) string {
	return c.Request.Host + c.Request.RequestURI
}

func setCache(c *bm.Context, res []byte) bool {
	h := c.Writer.Header()
	h.Set("Content-Type", "text/html; charset=utf8")
	h.Set("Cache-Control", fmt.Sprintf("max-age=%d", conf.Conf.Seo.MaxAge))
	etag := ETag(res)
	h.Set("ETag", etag)

	if etag == c.Request.Header.Get("If-None-Match") {
		c.Writer.WriteHeader(http.StatusNotModified)
		return true
	}
	return false
}

// ETag get etag for cache
func ETag(res []byte) string {
	return fmt.Sprintf(`W/"%x-%x"`, len(res), sha1.Sum(res))
}

// logUA log User-Agent and Url
func logUA(c *bm.Context) {
	log.Infov(c,
		log.KV("ua", c.Request.Header.Get("User-Agent")),
		log.KV("url", fmt.Sprintf("%s%s", c.Request.Host, c.Request.URL)),
	)
}
