package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	schmdl "go-common/app/service/main/riot-search/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

const _search = "http://api.bilibili.co/x/internal/riot-search/arc/ids"

// SearchArcs return archive ids by aids.
func (d *Dao) SearchArcs(c context.Context, keyword string, ids []int64, pn, ps int) (res *schmdl.IDsResp, err error) {
	params := url.Values{}
	params.Set("ids", xstr.JoinInts(ids))
	params.Set("keyword", keyword)
	params.Set("pn", strconv.Itoa(pn))
	params.Set("ps", strconv.Itoa(ps))
	ip := metadata.String(c, metadata.RemoteIP)
	var (
		resp = &struct {
			Code int             `json:"code"`
			Data *schmdl.IDsResp `json:"data"`
		}{}
	)
	if err = d.httpClient.Post(c, _search, ip, params, &resp); err != nil {
		log.Error("s.httpClient.Post(%s) error(%v)", _search+"?"+params.Encode(), err)
		return
	}
	log.Info("searchArcs(%s) error(%v)", _search+"?"+params.Encode(), err)
	if resp.Code != 0 {
		err = fmt.Errorf("code is:%d", resp.Code)
		log.Error("s.httpClient.Post(%s) error(%v)", _search+"?"+params.Encode(), err)
		return
	}
	return resp.Data, nil
}
