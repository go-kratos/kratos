package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/app/infra/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// push config update
func push(c *bm.Context) {
	var (
		err      error
		svr      string
		buildVer string
		ver      int64
		env      string
	)
	query := c.Request.Form
	verStr := query.Get("version")
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if buildVer = query.Get("build_ver"); buildVer == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ver, err = strconv.ParseInt(verStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	service := &model.Service{Name: svr, BuildVersion: buildVer, Version: ver, Env: env}
	// update & write cache
	c.JSON(nil, confSvc.Push(c, service))
}

// hosts client hosts
func hosts(c *bm.Context) {
	var (
		err  error
		svr  string
		data []*model.Host
		env  string
	)
	query := c.Request.URL.Query()
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, err = confSvc.Hosts(c, svr, env); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// versions client versions which the configuration is complete
func versions(c *bm.Context) {
	var (
		err  error
		svr  string
		data *model.Versions
		env  string
		bver string
	)
	query := c.Request.URL.Query()
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if bver = query.Get("build"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, err = confSvc.VersionSuccess(c, svr, env, bver); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// versions client versions which the configuration is complete
func versionIng(c *bm.Context) {
	var (
		err  error
		svr  string
		data []int64
		env  string
	)
	query := c.Request.URL.Query()
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if data, err = confSvc.VersionIng(c, svr, env); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// config get config file
func config(c *bm.Context) {
	var (
		err      error
		svr      string
		host     string
		buildVer string
		version  int64
		env      string
		token    string
	)

	query := c.Request.URL.Query()
	verStr := query.Get("version")
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if token = query.Get("token"); token == "" {
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
	if version, err = strconv.ParseInt(verStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	service := &model.Service{Name: svr, BuildVersion: buildVer, Env: env, Token: token, Version: version, Host: host}
	data, err := confSvc.Config(c, service)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// file get one file value
func file(c *bm.Context) {
	var (
		err      error
		svr      string
		buildVer string
		env      string
		token    string
		file     string
		ver      int64
		data     string
	)

	query := c.Request.URL.Query()
	// params
	if buildVer = query.Get("build"); buildVer == "" {
		data = "build is null"
	}
	if !strings.HasPrefix(buildVer, "shsb") && !strings.HasPrefix(buildVer, "shylf") &&
		query.Get("zone") != "" && query.Get("env") != "" && query.Get("treeid") != "" {
		file2(c)
		return
	}
	if verStr := query.Get("version"); verStr == "" {
		ver = model.UnknownVersion
	} else {
		if ver, err = strconv.ParseInt(verStr, 10, 64); err != nil {
			data = "version must be num"
		}
	}
	if svr = query.Get("service"); svr == "" {
		data = "service is null"
	}
	if env = query.Get("environment"); env == "" {
		data = "environment is null"
	}
	if token = query.Get("token"); token == "" {
		data = "token is null"
	}
	if file = query.Get("fileName"); file == "" {
		data = "fileName is null"
	}
	service := &model.Service{Name: svr, BuildVersion: buildVer, Env: env, File: file, Token: token, Version: ver}
	if data == "" {
		if data, err = confSvc.File(c, service); err != nil {
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
func check(c *bm.Context) {
	var (
		err      error
		svr      string
		host     string
		buildVer string
		ip       string
		ver      int64
		env      string
		token    string
		appoint  int64
		query    = c.Request.URL.Query()
	)
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if token = query.Get("token"); token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
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
	evt, err := confSvc.CheckVersion(c, rhost, env, token)
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
	confSvc.Unsub(svr, host, env)
}

//clear  host in redis
func clearhost(c *bm.Context) {
	var (
		svr string
		env string
	)
	query := c.Request.Form
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, confSvc.ClearHost(c, svr, env))
}

// versions client versions which the configuration is complete
func builds(c *bm.Context) {
	var (
		svr     string
		bs, bs2 []string
		env     string
	)
	query := c.Request.URL.Query()
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	bs, _ = confSvc.Builds(c, svr, env)
	bs2, _ = confSvc2.TmpBuilds(svr, env)
	bs = append(bs, bs2...)
	if len(bs) == 0 {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(bs, nil)
}

func addConfigs(c *bm.Context) {
	var (
		svr    string
		env    string
		token  string
		user   string
		data   map[string]string
		err    error
		values = c.Request.PostForm
	)
	// params
	if svr = values.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = values.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if token = values.Get("token"); token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if user = values.Get("user"); user == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(values.Get("data")), &data); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, confSvc.AddConfigs(c, svr, env, token, user, data))
}

func copyConfigs(c *bm.Context) {
	var (
		svr    string
		env    string
		token  string
		build  string
		user   string
		err    error
		ver    int64
		values = c.Request.PostForm
	)
	// params
	if svr = values.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = values.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if token = values.Get("token"); token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if user = values.Get("user"); user == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if build = values.Get("build"); build == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ver, err = confSvc.CopyConfigs(c, svr, env, token, user, build); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ver, nil)
}

func updateConfigs(c *bm.Context) {
	var (
		svr    string
		env    string
		token  string
		ver    int64
		user   string
		data   map[string]string
		err    error
		values = c.Request.PostForm
	)
	// params
	if svr = values.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = values.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if ver, err = strconv.ParseInt(values.Get("version"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if token = values.Get("token"); token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if user = values.Get("user"); user == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = json.Unmarshal([]byte(values.Get("data")), &data); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, confSvc.UpdateConfigs(c, svr, env, token, user, ver, data))
}

// configN get config namespace file
func configN(c *bm.Context) {
	var (
		err      error
		svr      string
		host     string
		buildVer string
		version  int64
		env      string
		token    string
	)

	query := c.Request.URL.Query()
	verStr := query.Get("version")
	// params
	if svr = query.Get("service"); svr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if env = query.Get("environment"); env == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if token = query.Get("token"); token == "" {
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
	if version, err = strconv.ParseInt(verStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	service := &model.Service{Name: svr, BuildVersion: buildVer, Env: env, Token: token, Version: version, Host: host}
	data, err := confSvc.Config2(c, service)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
