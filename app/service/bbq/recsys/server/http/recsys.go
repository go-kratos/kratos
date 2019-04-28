package http

import (
	"encoding/json"
	"go-common/app/service/bbq/recsys/api/grpc/v1"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
)

// start this just a example
func start(c *bm.Context) {
	arg := new(v1.RecsysRequest)

	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Start(c, arg))
}

func reqRecsys(c *bm.Context) {
	res, _ := ioutil.ReadAll(c.Request.Body)
	arg := new(v1.RecsysRequest)
	json.Unmarshal(res, &arg)

	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Start(c, arg))
}

func relatedRecsys(c *bm.Context) {
	res, _ := ioutil.ReadAll(c.Request.Body)
	arg := new(v1.RecsysRequest)
	json.Unmarshal(res, &arg)

	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.RelatedRecService(c, arg))
}

func upsRecsys(c *bm.Context) {
	res, _ := ioutil.ReadAll(c.Request.Body)
	arg := new(v1.RecsysRequest)
	json.Unmarshal(res, &arg)

	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.UpsRecService(c, arg))
}
