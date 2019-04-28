package service

import (
	"fmt"
	"time"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/dao"
	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/util"

	"go-common/library/log"
)

// ToDaoArea .
func ToDaoArea(area string) int {
	if d, ok := conf.Areas[area]; ok {
		return d
	}
	return int(dao.AreaReply)
}

// ToModelArea .
func ToModelArea(area int) string {
	for m, d := range conf.Areas {
		if d == area {
			return m
		}
	}
	return ""
}

// ToDaoState .
func ToDaoState(state string) int {
	switch state {
	case model.StateDefault:
		return dao.StateDefault
	case model.StateDeleted:
		return dao.StateDeleted
	default:
		return dao.StateDefault
	}
}

// ToModelState .
func ToModelState(state int) string {
	switch state {
	case dao.StateDefault:
		return model.StateDefault
	case dao.StateDeleted:
		return model.StateDeleted
	default:
		return ""
	}
}

// ToModelOperation .
func ToModelOperation(op int) string {
	switch op {
	case dao.OperationLimit:
		return model.OperationLimit
	case dao.OperationRestrictLimit:
		return model.OperationRestrictLimit
	case dao.OperationPutToWhiteList:
		return model.OperationPutToWhiteList
	case dao.OperationIgnore:
		return model.OperationIgnore
	default:
		return ""
	}
}

// ToDaoOperation .
func ToDaoOperation(op string) int {
	switch op {
	case model.OperationLimit:
		return dao.OperationLimit
	case model.OperationRestrictLimit:
		return dao.OperationRestrictLimit
	case model.OperationPutToWhiteList:
		return dao.OperationPutToWhiteList
	case model.OperationIgnore:
		return dao.OperationIgnore
	default:
		return dao.OperationLimit
	}
}

// ToModelRules .
func ToModelRules(rules []*dao.Rule) []*model.Rule {
	if rules == nil {
		return nil
	}
	result := make([]*model.Rule, len(rules))
	for i, r := range rules {
		if r == nil {
			result[i] = nil
		} else {
			result[i] = ToModelRule(r)
		}
	}
	return result
}

// ToModelRule .
func ToModelRule(d *dao.Rule) *model.Rule {
	if d == nil {
		return nil
	}
	r := &model.Rule{
		ID:            d.ID,
		Area:          ToModelArea(d.Area),
		AllowedCounts: d.AllowedCounts,
		DurationSec:   d.DurationSec,
	}
	// limit type
	switch d.LimitType {
	case dao.LimitTypeDefaultLimit:
		r.LimitType = model.LimitTypeDefault
	case dao.LimitTypeRestrictLimit:
		r.LimitType = model.LimitTypeRestrict
	}
	// limit scope
	switch d.LimitScope {
	case dao.LimitScopeGlobal:
		r.LimitScope = model.LimitScopeGlobal
	case dao.LimitScopeLocal:
		r.LimitScope = model.LimitScopeLocal
	}
	return r
}

// ToDaoRule .
func ToDaoRule(m *model.Rule) *dao.Rule {
	if m == nil {
		return nil
	}
	d := &dao.Rule{
		ID:            m.ID,
		AllowedCounts: m.AllowedCounts,
		DurationSec:   m.DurationSec,
		Area:          ToDaoArea(m.Area),
	}
	// limit type
	switch m.LimitType {
	case model.LimitTypeDefault:
		d.LimitType = dao.LimitTypeDefaultLimit
	case model.LimitTypeRestrict:
		d.LimitType = dao.LimitTypeRestrictLimit
	}
	// limit scope
	switch m.LimitScope {
	case model.LimitScopeGlobal:
		d.LimitScope = dao.LimitScopeGlobal
	case model.LimitScopeLocal:
		d.LimitScope = dao.LimitScopeLocal
	}
	return d
}

// ToModelKeywords .
func ToModelKeywords(keywords []*dao.Keyword) []*model.Keyword {
	if keywords == nil {
		return nil
	}
	result := make([]*model.Keyword, len(keywords))
	for i, r := range keywords {
		if r == nil {
			result[i] = nil
		} else {
			result[i] = ToModelKeyword(r)
		}
	}
	return result
}

// ToModelKeyword .
func ToModelKeyword(d *dao.Keyword) *model.Keyword {
	if d == nil {
		return nil
	}
	k := &model.Keyword{
		ID:            d.ID,
		Content:       d.Content,
		RegexpName:    d.RegexpName,
		HitCounts:     d.HitCounts,
		OriginContent: d.OriginContent,
		State:         ToModelState(d.State),
		CTime:         util.JSONTime(d.CTime),
		MTime:         util.JSONTime(d.MTime),
		Area:          ToModelArea(d.Area),
	}
	switch d.Tag {
	case dao.KeywordTagDefaultLimit:
		k.Tag = model.KeywordTagDefaultLimit
	case dao.KeywordTagRestrictLimit:
		k.Tag = model.KeywordTagRestrictLimit
	case dao.KeywordTagWhite:
		k.Tag = model.KeywordTagWhite
	case dao.KeywordTagBlack:
		k.Tag = model.KeywordTagBlack
	default:
		log.Error("unknown keyword tag %q", d.Tag)
		return nil
	}
	return k
}

// ToDaoKeywords .
func ToDaoKeywords(ks []*model.Keyword) []*dao.Keyword {
	if ks == nil {
		return nil
	}
	result := make([]*dao.Keyword, 0)
	for _, k := range ks {
		result = append(result, ToDaoKeyword(k))
	}
	return result
}

// ToDaoKeyword .
func ToDaoKeyword(k *model.Keyword) *dao.Keyword {
	if k == nil {
		return nil
	}
	d := &dao.Keyword{
		ID:            k.ID,
		Content:       k.Content,
		RegexpName:    k.RegexpName,
		Area:          ToDaoArea(k.Area),
		State:         ToDaoState(k.State),
		OriginContent: k.OriginContent,
		CTime:         time.Time(k.CTime),
		HitCounts:     k.HitCounts,
	}
	switch k.Tag {
	case model.KeywordTagDefaultLimit:
		d.Tag = dao.KeywordTagDefaultLimit
	case model.KeywordTagRestrictLimit:
		d.Tag = dao.KeywordTagRestrictLimit
	case model.KeywordTagWhite:
		d.Tag = dao.KeywordTagWhite
	case model.KeywordTagBlack:
		d.Tag = dao.KeywordTagBlack
	default:
		log.Error("Unknown tag %q", k.Tag)
		return nil
	}
	return d
}

// ToModelRegexps .
func ToModelRegexps(rs []*dao.Regexp) []*model.Regexp {
	if rs == nil {
		return nil
	}
	result := make([]*model.Regexp, len(rs))
	for i, r := range rs {
		if r == nil {
			result[i] = nil
		} else {
			result[i] = ToModelRegexp(r)
		}
	}
	return result
}

// ToModelRegexp .
func ToModelRegexp(d *dao.Regexp) *model.Regexp {
	if d == nil {
		return nil
	}
	return &model.Regexp{
		ID:        d.ID,
		Area:      ToModelArea(d.Area),
		AdminID:   d.AdminID,
		Name:      d.Name,
		Content:   d.Content,
		State:     ToModelState(d.State),
		Operation: ToModelOperation(d.Operation),
		CTime:     util.JSONTime(d.CTime),
		MTime:     util.JSONTime(d.MTime),
	}
}

// ToDaoRegexps .
func ToDaoRegexps(regs []*model.Regexp) []*dao.Regexp {
	if regs == nil {
		return nil
	}
	result := make([]*dao.Regexp, 0)
	for _, reg := range regs {
		result = append(result, ToDaoRegexp(reg))
	}
	return result
}

// ToDaoRegexp .
func ToDaoRegexp(m *model.Regexp) *dao.Regexp {
	if m == nil {
		return nil
	}
	return &dao.Regexp{
		ID:        m.ID,
		Area:      ToDaoArea(m.Area),
		State:     ToDaoState(m.State),
		Name:      m.Name,
		AdminID:   m.AdminID,
		Content:   m.Content,
		Operation: ToDaoOperation(m.Operation),
		CTime:     time.Time(m.CTime),
		MTime:     time.Time(m.MTime),
	}
}

// ToDaoCond .
func ToDaoCond(cond *Condition) *dao.Condition {
	if cond == nil {
		return nil
	}
	res := &dao.Condition{
		Pagination: cond.Pagination,
		HitCounts:  cond.HitCounts,
		Search:     cond.Search,
		Offset:     cond.Offset,
		Limit:      cond.Limit,
		Order:      cond.Order,
		OrderBy:    cond.OrderBy,
		Contents:   cond.Contents,
	}
	if len(cond.Area) > 0 {
		res.Area = fmt.Sprintf("%d", ToDaoArea(cond.Area))
	}
	if len(cond.State) > 0 {
		res.State = fmt.Sprintf("%d", ToDaoState(cond.State))
	}
	if res.OrderBy == "" {
		res.OrderBy = "id"
	}
	if cond.StartTime != nil {
		res.StartTime = cond.StartTime.Format(util.TimeFormat)
	}
	if cond.EndTime != nil {
		res.EndTime = cond.EndTime.Format(util.TimeFormat)
	}
	if cond.LastModifiedTime != nil {
		res.LastModifiedTime = cond.LastModifiedTime.Format(util.TimeFormat)
	}
	for _, tag := range cond.Tags {
		switch tag {
		case model.KeywordTagBlack:
			res.Tags = append(res.Tags, fmt.Sprintf("%d", dao.KeywordTagBlack))
		case model.KeywordTagWhite:
			res.Tags = append(res.Tags, fmt.Sprintf("%d", dao.KeywordTagWhite))
		case model.KeywordTagDefaultLimit:
			res.Tags = append(res.Tags, fmt.Sprintf("%d", dao.KeywordTagDefaultLimit))
		case model.KeywordTagRestrictLimit:
			res.Tags = append(res.Tags, fmt.Sprintf("%d", dao.KeywordTagRestrictLimit))
		}
	}
	switch cond.LimitType {
	case model.LimitTypeBlack:
		res.LimitType = fmt.Sprintf("%d", dao.LimitTypeBlack)
	case model.LimitTypeWhite:
		res.LimitType = fmt.Sprintf("%d", dao.LimitTypeWhite)
	case model.LimitTypeDefault:
		res.LimitType = fmt.Sprintf("%d", dao.LimitTypeDefaultLimit)
	case model.LimitTypeRestrict:
		res.LimitType = fmt.Sprintf("%d", dao.LimitTypeRestrictLimit)
	}

	switch cond.LimitScope {
	case model.LimitScopeGlobal:
		res.LimitScope = fmt.Sprintf("%d", dao.LimitScopeGlobal)
	case model.LimitScopeLocal:
		res.LimitScope = fmt.Sprintf("%d", dao.LimitScopeLocal)
	}
	return res
}
