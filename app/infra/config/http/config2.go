package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-common/app/infra/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// versions client versions which the configuration is complete
func versions2(c *bm.Context) {
	var (
		err   error
		svr   string
		data  *model.Versions
		bver  string
		token string
		env   string
		zone  string
	)
	query := c.Request.URL.Query()
	// params
	svr = query.Get("service")
	if svr == "" {
		token = query.Get("token")
		zone = query.Get("zone")
		env = query.Get("env")
		if zone == "" || env == "" || token == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if svr, err = confSvc2.AppService(zone, env, token); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if bver = query.Get("build"); bver == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, err = confSvc2.VersionSuccess(c, svr, bver); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// config get config file
func config2(c *bm.Context) {
	var (
		err      error
		svr      string
		buildVer string
		version  int64
		token    string
		ids      []int64
		zone     string
		env      string
	)
	query := c.Request.URL.Query()
	// params
	if token = query.Get("token"); token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	svr = query.Get("service")
	if svr == "" {
		zone = query.Get("zone")
		env = query.Get("env")
		if zone == "" || env == "" || token == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if svr, err = confSvc2.AppService(zone, env, token); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if buildVer = query.Get("build"); buildVer == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if version, err = strconv.ParseInt(query.Get("version"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if idsStr := query.Get("ids"); len(idsStr) != 0 {
		if err = json.Unmarshal([]byte(query.Get("ids")), &ids); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	data, err := confSvc2.Config(c, svr, token, version, ids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// file get one file value
func file2(c *bm.Context) {
	var (
		err      error
		svr      string
		buildVer string
		token    string
		file     string
		ver      int64
		treeID   string
		env      string
		zone     string
		data     string
	)

	query := c.Request.URL.Query()
	// params
	if verStr := query.Get("version"); verStr == "" {
		ver = model.UnknownVersion
	} else {
		if ver, err = strconv.ParseInt(verStr, 10, 64); err != nil {
			data = "version must be num"
		}
	}
	if env = query.Get("env"); env == "" {
		data = "env is null"
	}
	if zone = query.Get("zone"); zone == "" {
		data = "zone is null"
	}
	if token = query.Get("token"); token == "" {
		data = "token is null"
	}
	if treeID = query.Get("treeid"); treeID == "" {
		if svr, err = confSvc2.AppService(zone, env, token); err != nil {
			data = "appid is null"
		}
	} else {
		svr = fmt.Sprintf("%s_%s_%s", treeID, env, zone)
	}
	if buildVer = query.Get("build"); buildVer == "" {
		data = "build is null"
	}
	if file = query.Get("fileName"); file == "" {
		data = "fileName is null"
	}
	service := &model.Service{Name: svr, BuildVersion: buildVer, File: file, Token: token, Version: ver}
	if data == "" {
		if data, err = confSvc2.File(c, service); err != nil {
			data = err.Error()
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if _, err = c.Writer.Write([]byte(data)); err != nil {
		log.Error("Response().Write(%v) error(%v)", data, err)
	}
}

// check check config version
func check2(c *bm.Context) {
	var (
		err      error
		svr      string
		host     string
		buildVer string
		ip       string
		ver      int64
		token    string
		appoint  int64
		zone     string
		env      string
		query    = c.Request.URL.Query()
	)
	// params
	if token = query.Get("token"); token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	svr = query.Get("service")
	if svr == "" {
		zone = query.Get("zone")
		env = query.Get("env")
		if zone == "" || env == "" || token == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if svr, err = confSvc2.AppService(zone, env, token); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if ip = query.Get("ip"); ip == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if host = query.Get("hostname"); host == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if buildVer = query.Get("build"); buildVer == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ver, err = strconv.ParseInt(query.Get("version"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	appoint, _ = strconv.ParseInt(query.Get("appoint"), 10, 64)
	// check config version
	rhost := &model.Host{Service: svr, Name: host, BuildVersion: buildVer, IP: ip, ConfigVersion: ver, Appoint: appoint, Customize: query.Get("customize")}
	evt, err := confSvc2.CheckVersion(c, rhost, token)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	// wait for version change
	select {
	case e := <-evt:
		c.JSON(e, nil)
	case <-time.After(time.Duration(cnf.PollTimeout)):
		c.JSON(nil, ecode.NotModified)
	case <-c.Writer.(http.CloseNotifier).CloseNotify():
		c.JSON(nil, ecode.NotModified)
	}
	confSvc2.Unsub(svr, host)
}

//clear  host in redis
func clearhost2(c *bm.Context) {
	var (
		svr   string
		zone  string
		env   string
		token string
		err   error
	)
	query := c.Request.Form
	svr = query.Get("service")
	if svr == "" {
		token = query.Get("token")
		zone = query.Get("zone")
		env = query.Get("env")
		if zone == "" || env == "" || token == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if svr, err = confSvc2.AppService(zone, env, token); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(nil, confSvc2.ClearHost(c, svr))
}

// versions client versions which the configuration is complete
func builds2(c *bm.Context) {
	var (
		err   error
		svr   string
		data  []string
		token string
		env   string
		zone  string
	)
	query := c.Request.URL.Query()
	// params
	svr = query.Get("service")
	if svr == "" {
		token = query.Get("token")
		zone = query.Get("zone")
		env = query.Get("env")
		if zone == "" || env == "" || token == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if svr, err = confSvc2.AppService(zone, env, token); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if data, err = confSvc2.Builds(c, svr); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// config get config file
func latest(c *bm.Context) {
	var (
		err      error
		svr      string
		buildVer string
		version  int64
		token    string
		ids      []int64
		zone     string
		env      string
		verStr   string
	)
	query := c.Request.URL.Query()
	// params
	if token = query.Get("token"); token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	svr = query.Get("service")
	if svr == "" {
		zone = query.Get("zone")
		env = query.Get("env")
		if zone == "" || env == "" || token == "" {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if svr, err = confSvc2.AppService(zone, env, token); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if buildVer = query.Get("build"); buildVer == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	verStr = query.Get("version")
	if len(verStr) > 0 {
		if version, err = strconv.ParseInt(verStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		rhost := &model.Host{Service: svr, BuildVersion: buildVer}
		version, err = confSvc2.CheckLatest(c, rhost, token)
		if err != nil {
			c.JSON(nil, err)
			return
		}
	}
	data, err := confSvc2.ConfigCheck(c, svr, token, version, ids)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
