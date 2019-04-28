package http

import (
	"encoding/json"
	"go-common/app/interface/main/creative/model/app"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"io/ioutil"
	"strconv"
	"strings"
)

func coverList(c *bm.Context) {
	params := c.Request.Form
	fnsStr := params.Get("fns")
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	cvs, err := dataSvc.Covers(c, mid, strings.Split(fnsStr, ","), metadata.String(c, metadata.RemoteIP))
	if err != nil {
		log.Error(" arcSvc.CoverList fnsStr(%s), mid(%d), ip(%s)  error(%v)", fnsStr, mid, metadata.String(c, metadata.RemoteIP), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(cvs, nil)
}

func uploadMaterial(c *bm.Context) {
	req := c.Request
	params := req.Form
	var (
		err error
		aid int64
		bs  []byte
	)
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", aidStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if bs, err = ioutil.ReadAll(req.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	req.Body.Close()
	var editors = &[]*app.Editor{}
	if err = json.Unmarshal(bs, editors); err != nil {
		log.Error("uploadMaterial json.Unmarshal(%s) error(%v)", string(bs), err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(appSvc.UploadMaterial(c, aid, *editors), nil)
}
