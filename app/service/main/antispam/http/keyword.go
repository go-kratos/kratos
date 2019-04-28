package http

import (
	"errors"
	"strconv"
	"time"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/dao"
	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/util"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	// ErrMissID .
	ErrMissID = errors.New("error: id is required")
	// ErrIllegalOp .
	ErrIllegalOp = errors.New("error: this operation on keyword is not allowed")
)

// GetKeywordSenders return keyword's sender list
func GetKeywordSenders(c *bm.Context) {
	params := c.Request.Form
	if _, err := getAdminID(params); err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	id, err := strconv.ParseInt(params.Get(ProtocolKeywordID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	k, err := Svr.GetKeywordByID(c, id)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	if k.State == model.StateDeleted {
		errResp(c, ecode.ServerErr, dao.ErrResourceNotExist)
		return
	}
	l, err := Svr.GetSenderIDsByKeywordID(c, id)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(l, nil)
}

// DeleteKeywords .
func DeleteKeywords(c *bm.Context) {
	params := c.Request.Form
	if _, err := getAdminID(params); err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	if params.Get(ProtocolKeywordIDs) == "" {
		errResp(c, ecode.RequestErr, errors.New("empty ids"))
		return
	}
	ids, err := util.StrToIntSli(params.Get(ProtocolKeywordIDs), ",")
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	if len(ids) == 0 {
		errResp(c, ecode.RequestErr, nil)
		return
	}
	ks, err := Svr.DeleteKeywords(c, ids)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(ks, nil)
}

// UpdateKeyword .
func UpdateKeyword(c *bm.Context) {
	params := c.Request.Form
	_, err := getAdminID(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	op := params.Get(ProtocolKeywordOperation)
	if op != ProtocolKeywordOpDefaultLimit &&
		op != ProtocolKeywordOpRestrictLimit &&
		op != ProtocolKeywordOpBlack &&
		op != ProtocolKeywordOpWhite {
		log.Error("%v", ErrIllegalOp)
		errResp(c, ecode.RequestErr, ErrIllegalOp)
		return
	}
	id, err := strconv.ParseInt(params.Get(ProtocolKeywordID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	k, err := Svr.GetKeywordByID(c, id)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	if k.State == model.StateDeleted {
		errResp(c, ecode.ServerErr, dao.ErrResourceNotExist)
		return
	}
	if k.Tag == op {
		c.JSON(k, nil)
		return
	}
	k, err = Svr.OpKeyword(c, id, op)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(k, nil)
}

// GetKeyword .
func GetKeyword(c *bm.Context) {
	params := c.Request.Form
	_, _, err := getAdminIDAndArea(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	if params.Get(ProtocolKeywordID) == "" {
		errResp(c, ecode.RequestErr, ErrMissID)
		return
	}
	id, err := strconv.ParseInt(params.Get(ProtocolKeywordID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	keyword, err := Svr.GetKeywordByID(c, id)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(keyword, nil)
}

// GetKeywords .
func GetKeywords(c *bm.Context) {
	params := c.Request.Form
	_, area, err := getAdminIDAndArea(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	var (
		ctimeStart, ctimeEnd string
		startTime, endTime   time.Time
	)
	cond := &Condition{
		Pagination: &util.Pagination{},
		Area:       area,
		State:      model.StateDefault,
		Search:     params.Get(ProtocolSearch),
		Order:      params.Get(ProtocolOrder),
		OrderBy:    params.Get(ProtocolOrderBy),
		Tag:        params.Get(ProtocolKeywordTag),
	}
	if ctimeStart = params.Get(ProtocolCTimeStart); ctimeStart != "" {
		startTime, err = time.Parse(util.TimeFormat, ctimeStart)
		if err != nil {
			log.Error("%v", err)
			errResp(c, ecode.RequestErr, err)
			return
		}
		cond.StartTime = &startTime
	}
	if ctimeEnd = params.Get(ProtocolCTimeEnd); ctimeEnd != "" {
		endTime, err = time.Parse(util.TimeFormat, ctimeEnd)
		if err != nil {
			log.Error("%v", err)
			errResp(c, ecode.RequestErr, err)
			return
		}
		cond.EndTime = &endTime
	}
	if ctimeStart != "" && ctimeEnd != "" {
		if startTime.After(endTime) {
			err = errors.New("start_time cannot be bigger than end_time")
			log.Error("%v", err)
			errResp(c, ecode.RequestErr, err)
			return
		}
	}
	var export bool
	if params.Get(ProtocolKeywordExport) != "" {
		export, err = strconv.ParseBool(params.Get(ProtocolKeywordExport))
		if err != nil {
			log.Error("%v", err)
			errResp(c, ecode.RequestErr, err)
			return
		}
	}
	if export {
		cond.PerPage = conf.Conf.ServiceOption.MaxExportRows
	} else {
		if params.Get(ProtocolCurPage) != "" {
			cond.CurPage, err = strconv.ParseInt(params.Get(ProtocolCurPage), 10, 64)
			if err != nil {
				log.Error("%v", err)
				errResp(c, ecode.RequestErr, err)
				return
			}
		}
		if perPage := params.Get(ProtocolPerPage); perPage != "" {
			if cond.PerPage, err = strconv.ParseInt(perPage, 10, 64); err != nil {
				log.Error("%v", err)
				errResp(c, ecode.RequestErr, err)
				return
			}
		}
	}
	if err = cond.Valid(); err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	keywords, total, err := Svr.GetKeywordsByCond(c, ToServiceCond(cond))
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	res := map[string]interface{}{}
	res[ProtocolTotalCounts] = total
	res[ProtocolData] = keywords
	c.JSONMap(res, nil)
}
