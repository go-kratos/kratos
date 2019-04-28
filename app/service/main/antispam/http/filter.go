package http

import (
	"strconv"

	"go-common/app/service/main/antispam/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Filter .
func Filter(c *bm.Context) {
	params := c.Request.Form

	senderID, err := strconv.ParseInt(params.Get(ProtocolKeywordSenderID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	oid, err := strconv.ParseInt(params.Get(ProtocolKeywordSubjectID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	susp := &model.Suspicious{
		SenderId: senderID,
		Content:  params.Get(ProtocolKeywordContent),
		Area:     params.Get(ProtocolArea),
		OId:      oid,
	}
	result, err := Svr.Filter(c, susp)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	c.JSON(result, nil)
}
