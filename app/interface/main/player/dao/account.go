package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/player/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const _blockTimeURI = "/api/member/getBlockAndMoralStatus"

// BlockTime get user block time from account by mid
func (d *Dao) BlockTime(c context.Context, mid int64) (res *model.BlockTime, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var rs struct {
		Code int              `json:"code"`
		Data *model.BlockTime `json:"data"`
	}
	if err = d.client.Get(c, d.blockTimeURL, ip, params, &rs); err != nil {
		log.Error("d.client.Get(%s,%d) error(%v)", d.blockTimeURL, mid, err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		log.Error("d.client.Get(%s,%d) error code(%d)", d.blockTimeURL, mid, rs.Code)
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.Data
	return
}
