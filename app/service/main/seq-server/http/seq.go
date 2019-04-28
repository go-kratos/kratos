package http

import (
	"strconv"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// seqID return id
func seqID(c *bm.Context) {
	params := c.Request.Form
	bsIDStr := params.Get("businessID")
	// check params
	bsID, err := strconv.ParseInt(bsIDStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	token := params.Get("token")
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(seqSvc.ID(c, bsID, token))
}

// seqID32 return id32
func seqID32(c *bm.Context) {
	params := c.Request.Form
	bsIDStr := params.Get("businessID")
	// check params
	bsID, err := strconv.ParseInt(bsIDStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	token := params.Get("token")
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err := seqSvc.ID32(c, bsID, token)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(id, nil)
}

// maxSeq update maxseq.
func maxSeq(c *bm.Context) {
	params := c.Request.Form
	bsIDStr := params.Get("businessID")
	maxSeqStr := params.Get("maxseq")
	stepStr := params.Get("step")
	token := params.Get("token")
	// check params
	bsID, err := strconv.ParseInt(bsIDStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	maxSeq, err := strconv.ParseInt(maxSeqStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	step, err := strconv.ParseInt(stepStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if token == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = seqSvc.UpMaxSeq(c, bsID, maxSeq, step, token); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(maxSeq, nil)
}
