package http

import (
	"strconv"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// doc get document
func doc(c *bm.Context) {
	params := c.Request.Form
	tsStr := params.Get("timestamp")
	moduleKey := params.Get("module_key")
	checkSumStr := params.Get("check_sum")
	// check login
	midIf, _ := c.Get("mid")
	mid := midIf.(int64)
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt64(%v) error(%v)", tsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	checkSum, err := strconv.ParseInt(checkSumStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt64(%v) error(%v)", checkSumStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// doc
	data, err := kvoSvr.Document(c, mid, moduleKey, ts, checkSum)
	if err != nil {
		if ecode.Cause(err) != ecode.NotModified {
			log.Error("kvoSvr.Document(%v,%v) error(%v)", mid, moduleKey, err)
		}
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

// addDoc add document
func addDoc(c *bm.Context) {
	params := c.Request.Form
	data := params.Get("data")
	tsStr := params.Get("timestamp")
	checkSumStr := params.Get("check_sum")
	moduleKey := params.Get("module_key")
	// check login
	midIf, _ := c.Get("mid")
	mid := midIf.(int64)
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt64(%v) error(%v)", tsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	checkSum, err := strconv.ParseInt(checkSumStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt64(%v) error(%v)", checkSumStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// add doc
	resp, err := kvoSvr.AddDocument(c, mid, moduleKey, data, ts, checkSum, time.Now())
	if err != nil {
		if ecode.Cause(err) != ecode.NotModified {
			log.Error("kvoSvr.AddDocument(%v,%v,%v,%v,%v) error(%v)", mid, moduleKey, data, ts, checkSum, err)
		}
		c.JSON(nil, err)
		return
	}
	c.JSON(resp, nil)
}
