package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
)

var (
	// Svr .
	Svr       service.Service
	verifySvc *verify.Verify
	authSvc   *auth.Auth
)

// Init .
func Init(c *conf.Config, s service.Service) {
	Svr = s
	verifySvc = verify.New(c.Verify)
	authSvc = auth.New(c.Auth)
	engine := bm.DefaultServer(c.BM)

	interRouter(engine)

	if err := engine.Start(); err != nil {
		log.Error("engine.Start() error(%v)", err)
		panic(err)
	}
}

func interRouter(e *bm.Engine) {
	e.GET("/monitor/ping", ping)

	e.GET("/register", register)
	e.GET("/x/internal/antispam/filter", authSvc.Guest, Filter)

	regexps := e.Group("/x/internal/antispam/regexps")
	regexps.GET("", verifySvc.Verify, GetRegexps)
	regexps.GET("/one", verifySvc.Verify, GetRegexp)
	regexps.POST("/add", verifySvc.Verify, AddRegexp)
	regexps.POST("/edit", verifySvc.Verify, EditRegexp)
	regexps.POST("/del", verifySvc.Verify, DeleteRegexp)
	regexps.POST("/recover", verifySvc.Verify, RecoverRegexp)

	rules := e.Group("/x/internal/antispam/rules")
	rules.GET("", verifySvc.Verify, GetRules)
	rules.GET("/one", verifySvc.Verify, GetRule)
	rules.POST("/add", verifySvc.Verify, AddRule)

	keywords := e.Group("/x/internal/antispam/keywords")
	keywords.GET("", verifySvc.Verify, GetKeywords)
	keywords.GET("/senders", verifySvc.Verify, GetKeywordSenders)
	keywords.GET("/one", verifySvc.Verify, GetKeyword)
	keywords.POST("/dels", verifySvc.Verify, DeleteKeywords)
	keywords.POST("/action", verifySvc.Verify, UpdateKeyword)
}

func ping(c *bm.Context) {
	if err := Svr.Ping(c); err != nil {
		log.Error("antispam service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(struct{}{}, nil)
}

func getAdminIDAndArea(params url.Values) (adminID int64, area string, err error) {
	adminID, err = getAdminID(params)
	if err != nil {
		return 0, "", err
	}
	area, err = parseArea(params)
	if err != nil {
		return 0, "", err
	}
	return adminID, area, nil
}

func parseArea(params url.Values) (string, error) {
	area := params.Get(ProtocolArea)
	if _, ok := conf.Areas[area]; !ok {
		err := fmt.Errorf("invalid area(%s)", area)
		log.Error("%v", err)
		return "", err
	}
	return area, nil
}

func getAdminID(params url.Values) (int64, error) {
	adminIDStr := params.Get(ProtocolAdminID)
	if adminIDStr == "" {
		err := errors.New("empty admin id")
		log.Error("%v", err)
		return 0, err
	}
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	return adminID, nil
}

func errResp(c *bm.Context, code interface{}, err error) {
	c.JSONMap(map[string]interface{}{
		ProtocolData:    code,
		ProtocolMessage: err.Error(),
	}, nil)
}
