package http

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/dao"
	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/util"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// GetRegexp .
func GetRegexp(c *bm.Context) {
	params := c.Request.Form
	_, err := getAdminID(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	id, err := strconv.ParseInt(params.Get(ProtocolRegexpID), 10, 64)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	regexp, err := Svr.GetRegexpByID(c, id)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	c.JSON(regexp, nil)
}

// GetRegexps .
func GetRegexps(c *bm.Context) {
	params := c.Request.Form
	_, area, err := getAdminIDAndArea(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	cond := &Condition{
		Area:       area,
		Pagination: &util.Pagination{},
	}
	if params.Get(ProtocolCurPage) != "" {
		if cond.CurPage, err = strconv.ParseInt(params.Get(ProtocolCurPage), 10, 64); err != nil {
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
	if err = cond.Valid(); err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	regexps, total, err := Svr.GetRegexpsByCond(c, ToServiceCond(cond))
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	res := map[string]interface{}{}
	res[ProtocolTotalCounts] = total
	res[ProtocolData] = regexps
	c.JSONMap(res, nil)

}

// EditRegexp .
func EditRegexp(c *bm.Context) {
	params := c.Request.Form
	if _, err := getAdminID(params); err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	id, err := strconv.ParseInt(params.Get(ProtocolRegexpID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	addOrEditRegexp(c, params, id)
}

// AddRegexp .
func AddRegexp(c *bm.Context) {
	params := c.Request.Form
	if _, _, err := getAdminIDAndArea(params); err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	var id int64
	addOrEditRegexp(c, params, id)
}

// RecoverRegexp .
func RecoverRegexp(c *bm.Context) {
	params := c.Request.Form
	adminID, err := getAdminID(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	id, err := strconv.ParseInt(params.Get(ProtocolRegexpID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	r, err := Svr.GetRegexpByID(c, id)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	if r.State != model.StateDeleted {
		c.JSON(r, nil)
		return
	}
	_, total, err := Svr.GetRegexpsByCond(c, ToServiceCond(&Condition{State: model.StateDefault, Area: r.Area}))
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	if total >= conf.Conf.ServiceOption.MaxRegexpCountsPerArea {
		err = fmt.Errorf("regexp counts cannot exceed %d", conf.Conf.ServiceOption.MaxRegexpCountsPerArea)
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	r.State, r.AdminID = model.StateDefault, adminID
	result, err := Svr.UpsertRegexp(c, r)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(result, nil)
}

// DeleteRegexp .
func DeleteRegexp(c *bm.Context) {
	params := c.Request.Form
	adminID, err := getAdminID(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	id, err := strconv.ParseInt(params.Get(ProtocolRegexpID), 10, 64)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	r, err := Svr.DeleteRegexp(c, id, adminID)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(r, nil)
}

func addOrEditRegexp(c *bm.Context, params url.Values, id int64) {
	adminID, err := getAdminID(params)
	if err != nil {
		errResp(c, ecode.RequestErr, err)
		return
	}
	name := params.Get(ProtocolRegexpName)
	if name == "" {
		err = errors.New("empty regexp name")
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	content := params.Get(ProtocolRegexpContent)
	if content == "" {
		err = errors.New("empty regexp content")
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	if _, err = regexp.Compile(content); err != nil {
		log.Error("%v", err)
		errResp(c, ecode.FilterIllegalRegexp, err)
		return
	}
	op := params.Get(ProtocolRegexpOperation)
	if op == "" {
		op = model.OperationLimit
	}
	area := params.Get(ProtocolArea)
	if id == 0 {
		_, total, err1 := Svr.GetRegexpsByCond(c, ToServiceCond(&Condition{State: model.StateDefault, Area: area}))
		if err1 != nil {
			errResp(c, ecode.ServerErr, err1)
			return
		}
		if total >= conf.Conf.ServiceOption.MaxRegexpCountsPerArea {
			err1 = fmt.Errorf("regexp counts cannot exceed %d", conf.Conf.ServiceOption.MaxRegexpCountsPerArea)
			log.Error("%v", err1)
			errResp(c, ecode.RequestErr, err1)
			return
		}

		oldR, err1 := Svr.GetRegexpByAreaAndContent(c, area, content)
		if err1 == nil {
			if oldR.State != model.StateDefault {
				oldR.State = model.StateDefault
				result, err2 := Svr.UpsertRegexp(c, oldR)
				if err2 != nil {
					errResp(c, ecode.ServerErr, err2)
					return
				}
				c.JSON(result, nil)
				return
			}
			err1 = fmt.Errorf("regexp(%s) already exists", content)
			log.Error("%v", err1)
			errResp(c, ecode.RequestErr, err1)
			return
		}
		if err != dao.ErrResourceNotExist {
			log.Error("%v", err1)
			errResp(c, ecode.ServerErr, err1)
			return
		}
	}
	r := &model.Regexp{
		ID:        id,
		Area:      area,
		AdminID:   adminID,
		Name:      name,
		Content:   content,
		Operation: op,
	}
	result, err := Svr.UpsertRegexp(c, r)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(result, nil)
}
