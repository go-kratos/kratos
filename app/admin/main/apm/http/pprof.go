package http

import (
	"io/ioutil"
	"net/http"

	"go-common/app/admin/main/apm/conf"
	mpprof "go-common/app/admin/main/apm/model/pprof"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func buildSvg(c *bm.Context) {
	v := new(struct {
		URL      string `form:"url" validate:"required"`
		URI      string `form:"uri" validate:"required"`
		SvgName  string `form:"name" validate:"required"`
		Time     int64  `form:"time" validate:"required"`
		HostName string `form:"hostname" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	if !(v.URI == "profile" || v.URI == "heap") {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = apmSvc.Pprof(v.URL, v.URI, v.SvgName, v.HostName, v.Time, 2)
	c.JSON(nil, err)
}

func readSvg(c *bm.Context) {
	v := new(struct {
		SvgName  string `form:"name" validate:"required"`
		URI      string `form:"uri" validate:"required"`
		HostName string `form:"hostname" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	url := conf.Conf.Pprof.Dir + "/" + v.SvgName + "_" + v.HostName + "_" + v.URI + ".svg"
	data, err := ioutil.ReadFile(url)
	if err != nil {
		log.Error("readfile error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Bytes(http.StatusOK, "image/svg+xml; charset=utf-8", data)
}

func heap(c *bm.Context) {
	v := new(struct {
		URL      string `form:"url" validate:"required"`
		URI      string `form:"uri" validate:"required"`
		SvgName  string `form:"name" validate:"required"`
		HostName string `form:"hostname" validate:"required"`
	})
	var err error
	var data []byte
	if err = c.Bind(v); err != nil {
		return
	}
	if v.URI != "heap" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = apmSvc.Pprof(v.URL, v.URI, v.SvgName, v.HostName, 1, 1)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	url := conf.Conf.Pprof.Dir + "/" + v.SvgName + "_" + v.HostName + "_" + v.URI + ".svg"
	data, err = ioutil.ReadFile(url)
	if err != nil {
		log.Error("readfile error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Bytes(http.StatusOK, "image/svg+xml; charset=utf-8", data)
}

func flame(c *bm.Context) {
	v := new(struct {
		URL      string `form:"url" validate:"required"`
		URI      string `form:"uri" validate:"required"`
		SvgName  string `form:"name" validate:"required"`
		Time     int64  `form:"time"`
		HostName string `form:"hostname" validate:"required"`
	})
	var err error
	var data []byte
	if err = c.Bind(v); err != nil {
		return
	}
	if !(v.URI == "profile" || v.URI == "heap") {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if v.URI == "heap" {
		err = apmSvc.Torch(c, v.URL, v.URI, v.SvgName, v.HostName, 1, 1)
		if err != nil {
			c.JSON(nil, err)
			return
		}
		url := conf.Conf.Pprof.Dir + "/" + v.SvgName + "_" + v.HostName + "_" + v.URI + "_flame.svg"
		data, err = ioutil.ReadFile(url)
		if err != nil {
			log.Error("readfile error(%v)", err)
			c.JSON(nil, err)
			return
		}
		c.Bytes(http.StatusOK, "image/svg+xml; charset=utf-8", data)
	} else {
		err = apmSvc.Torch(c, v.URL, v.URI, v.SvgName, v.HostName, v.Time, 2)
		c.JSON(nil, err)
	}
}

func activeWarning(c *bm.Context) {
	// var (
	// 	err  error
	// 	body []byte
	// )
	// if body, err = ioutil.ReadAll(c.Request.Body); err != nil {
	// 	c.JSON(nil, fmt.Errorf("warning body empty"))
	// 	return
	// }
	// if err = apmSvc.ActiveWarning(c, string(body)); err != nil {
	// 	c.JSON(nil, err)
	// 	return
	// }
	// ioutil.ReadAll You need close body.
	c.JSON(nil, nil)
}

func pprof(c *bm.Context) {
	var (
		err error
		pws = make([]*mpprof.Warn, 0)
		req = &mpprof.Params{}
	)
	if err = c.Bind(req); err != nil {
		c.JSON(nil, err)
		return
	}
	if pws, err = apmSvc.PprofWarn(c, req); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pws, nil)
}
