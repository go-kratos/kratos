package http

import (
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func factors(c *bm.Context) {
	var (
		data []*model.Factors
		err  error
	)
	data, err = spySrv.Factors(c)
	if err != nil {
		log.Error("spySrv.Factors error(%v)", err)
		c.JSON(data, err)
		return
	}
	c.JSON(data, err)
}

func updateFactor(c *bm.Context) {
	var (
		params = c.Request.Form
		idsStr = params.Get("ids")
		name   = params.Get("name")
		fs     []*model.Factor
		err    error
	)
	if name == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	idsArr := strings.Split(idsStr, ",")
	for _, idStr := range idsArr {
		id, theErr := strconv.ParseInt(idStr, 10, 64)
		if theErr != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		factorVal, theErr := strconv.ParseFloat(params.Get("id_"+idStr), 32)
		if theErr != nil || factorVal < 0 || factorVal > 1 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		fs = append(fs, &model.Factor{ID: id, FactorVal: float32(factorVal)})
	}
	err = spySrv.UpdateFactor(c, fs, name)
	if err != nil {
		log.Error("spySrv.UpdateFactor(%v) error(%v)", fs, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func addFactor(c *bm.Context) {
	var (
		params     = c.Request.Form
		nickname   = params.Get("nick_name")
		serviceID  int64
		eventID    int64
		groupID    int64
		riskLevel  int64
		factorVal  float64
		categoryID int64
		err        error
	)
	if serviceID, err = strconv.ParseInt(params.Get("service_id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if eventID, err = strconv.ParseInt(params.Get("event_id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if groupID, err = strconv.ParseInt(params.Get("group_id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if riskLevel, err = strconv.ParseInt(params.Get("risk_level"), 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if factorVal, err = strconv.ParseFloat(params.Get("factor_val"), 32); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if categoryID, err = strconv.ParseInt(params.Get("category_id"), 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if factorVal > 1 || factorVal <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if riskLevel < 1 || riskLevel > 10 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.AddFactor(c, &model.Factor{
		NickName:   nickname,
		ServiceID:  serviceID,
		EventID:    eventID,
		GroupID:    groupID,
		RiskLevel:  int8(riskLevel),
		FactorVal:  float32(factorVal),
		CTime:      time.Now(),
		CategoryID: int8(categoryID),
	})
	if err != nil {
		log.Error("spySrv.AddFactor error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func addEvent(c *bm.Context) {
	var (
		params    = c.Request.Form
		name      = params.Get("name")
		nickName  = params.Get("nick_name")
		err       error
		serviceID int64
		status    int64
	)
	if serviceID, err = strconv.ParseInt(params.Get("service_id"), 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if status, err = strconv.ParseInt(params.Get("status"), 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.AddEvent(c, &model.Event{
		Name:      name,
		NickName:  nickName,
		ServiceID: serviceID,
		Status:    int8(status),
		CTime:     time.Now(),
		MTime:     time.Now(),
	})
	if err != nil {
		log.Error("spySrv.AddEvent error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func addService(c *bm.Context) {
	var (
		params   = c.Request.Form
		name     = params.Get("name")
		nickName = params.Get("nick_name")
		err      error
		status   int64
	)
	if status, err = strconv.ParseInt(params.Get("status"), 10, 8); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.AddService(c, &model.Service{
		Name:     name,
		NickName: nickName,
		Status:   int8(status),
		CTime:    time.Now(),
		MTime:    time.Now(),
	})
	if err != nil {
		log.Error("spySrv.AddService error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func addGroup(c *bm.Context) {
	var (
		params = c.Request.Form
		name   = params.Get("name")
		err    error
	)
	err = spySrv.AddGroup(c, &model.FactorGroup{
		Name:  name,
		CTime: time.Now(),
	})
	if err != nil {
		log.Error("spySrv.AddGroup error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func updateEventName(c *bm.Context) {
	var (
		params   = c.Request.Form
		nickname string
		id       int64
		err      error
	)
	nickname = params.Get("nickname")
	if err != nil || len(nickname) > model.MaxRemarkLen {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	id, err = strconv.ParseInt(params.Get("id"), 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	err = spySrv.UpdateEventName(c, &model.Event{
		ID:       id,
		NickName: nickname,
	})
	if err != nil {
		log.Error("spySrv.UpdateEventName error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
