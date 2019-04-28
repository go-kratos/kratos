package common

import (
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/admin/main/feed/model/common"
	searchModel "go-common/app/admin/main/search/model"
	bm "go-common/library/net/http/blademaster"
)

const (
	logURL = "/x/admin/search/log"
)

//LogAction log action
func (s *Service) LogAction(c *bm.Context, typ, ps, pn int64, ctimeFrom, ctimeTo, uName string) (res *common.LogManagers, err error) {
	var (
		items []*common.LogManager
	)
	res = &common.LogManagers{}
	params := url.Values{}
	params.Set("appid", "log_audit")
	params.Set("business", strconv.FormatUint(common.BusinessID, 10))
	params.Set("order", "ctime")
	params.Set("type", strconv.FormatInt(typ, 10))
	params.Set("ps", strconv.FormatInt(ps, 10))
	params.Set("pn", strconv.FormatInt(pn, 10))
	if ctimeFrom != "" {
		params.Set("ctime_from", ctimeFrom)
	}
	if ctimeTo != "" {
		params.Set("ctime_to", ctimeTo)
	}
	if uName != "" {
		params.Set("uname", uName)
	}
	type log struct {
		Code int                       `json:"code"`
		Data *searchModel.SearchResult `json:"data"`
	}
	l := &log{}
	if err = s.client.Get(c, s.managerURL+logURL, "", params, l); err != nil {
		return
	}
	var logS []*common.LogSearch
	for _, v := range l.Data.Result {
		log := &common.LogSearch{}
		if err = json.Unmarshal(v, log); err != nil {
			return
		}
		logS = append(logS, log)
	}
	for _, v := range logS {
		tmp := &common.LogManager{
			OID:       v.OID,
			Uname:     v.Uname,
			UID:       v.UID,
			Type:      v.Type,
			ExtraData: v.ExtraData,
			Action:    v.Action,
			CTime:     v.CTime,
		}
		items = append(items, tmp)
	}
	res.Item = items
	res.Page.TotalItems = int(l.Data.Page.Total)
	res.Page.PageSize = l.Data.Page.Ps
	res.Page.CurrentPage = l.Data.Page.Pn
	return
}
