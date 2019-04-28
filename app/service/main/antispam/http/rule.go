package http

import (
	"errors"
	"fmt"
	"strconv"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// GetRule .
func GetRule(c *bm.Context) {
	params := c.Request.Form
	_, area, err := getAdminIDAndArea(params)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	if params.Get(ProtocolRuleLimitType) == "" ||
		params.Get(ProtocolRuleLimitScope) == "" {
		err = errors.New("either limit_type or limit_scope is nil")
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	rule, err := Svr.GetRuleByAreaAndLimitTypeAndScope(c,
		area, params.Get(ProtocolRuleLimitType), params.Get(ProtocolRuleLimitScope))
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(rule, nil)
}

// GetRules .
func GetRules(c *bm.Context) {
	params := c.Request.Form
	_, area, err := getAdminIDAndArea(params)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	rules, err := Svr.GetRuleByArea(c, area)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(rules, nil)
}

// AddRule .
func AddRule(c *bm.Context) {
	params := c.Request.Form
	_, area, err := getAdminIDAndArea(params)
	if err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	typ := params.Get(ProtocolRuleLimitType)
	if typ != model.LimitTypeDefault && typ != model.LimitTypeRestrict {
		err = fmt.Errorf("illegal limit type %q", typ)
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	scope := params.Get(ProtocolRuleLimitScope)
	if scope != model.LimitScopeLocal && scope != model.LimitScopeGlobal {
		err = fmt.Errorf("illegal limit scope %q", scope)
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	r := &model.Rule{
		Area:       area,
		LimitType:  typ,
		LimitScope: scope,
	}
	allowedCounts := params.Get(ProtocolRuleAllowedCounts)
	if r.AllowedCounts, err = strconv.ParseInt(allowedCounts, 10, 64); err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	if r.DurationSec, err = strconv.ParseInt(params.Get(ProtocolRuleDuration), 10, 64); err != nil {
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	if r.DurationSec <= 0 || r.AllowedCounts <= 0 {
		err = fmt.Errorf("both durationSec(%d) and allowedCounts(%d) must be greater than 0", r.AllowedCounts, r.DurationSec)
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	if r.DurationSec > conf.Conf.MaxDurationSec || r.AllowedCounts > conf.Conf.MaxAllowedCounts {
		err = fmt.Errorf("either durationSec(%d) or allowedCounts(%d) exceed maxDurationSec(%d), maxAllowedCounts(%d)",
			r.AllowedCounts, r.DurationSec, conf.Conf.MaxDurationSec, conf.Conf.MaxAllowedCounts)
		log.Error("%v", err)
		errResp(c, ecode.RequestErr, err)
		return
	}
	r, err = Svr.UpsertRule(c, r)
	if err != nil {
		errResp(c, ecode.ServerErr, err)
		return
	}
	c.JSON(r, nil)
}
