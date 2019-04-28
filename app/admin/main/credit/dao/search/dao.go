package search

import (
	"context"
	"time"

	"go-common/app/admin/main/credit/conf"
	"go-common/app/admin/main/credit/model"
	"go-common/app/admin/main/credit/model/blocked"
	"go-common/app/admin/main/credit/model/search"
	"go-common/library/database/elastic"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao .
type Dao struct {
	elastic *elastic.Elastic
}

// New .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		elastic: elastic.NewElastic(nil),
	}
	return
}

// Blocked get search blocked.
func (d *Dao) Blocked(c context.Context, arg *blocked.ArgBlockedSearch) (ids []int64, pager *blocked.Pager, err error) {
	req := d.elastic.NewRequest(blocked.BusinessBlockedInfo).Index(blocked.TableBlockedInfo).Fields("id")
	if arg.Keyword != blocked.SearchDefaultString {
		req.WhereLike([]string{"origin_content"}, []string{arg.Keyword}, true, elastic.LikeLevelHigh)
	}
	if arg.UID != blocked.SearchDefaultNum {
		req.WhereEq("uid", arg.UID)
	}
	if arg.OPID != blocked.SearchDefaultNum {
		req.WhereEq("oper_id", arg.OPID)
	}
	if arg.OriginType != blocked.SearchDefaultNum {
		req.WhereEq("origin_type", arg.OriginType)
	}
	if arg.BlockedType != blocked.SearchDefaultNum {
		req.WhereEq("blocked_type", arg.BlockedType)
	}
	if arg.PublishStatus != blocked.SearchDefaultNum {
		req.WhereEq("publish_status", arg.PublishStatus)
	}
	req.WhereRange("punish_time", arg.Start, arg.End, elastic.RangeScopeLcRc)
	req.WhereEq("status", blocked.SearchDefaultStatus)
	req.Pn(arg.PN).Ps(arg.PS).Order(arg.Order, arg.Sort)
	var res *search.ReSearchData
	if err = req.Scan(c, &res); err != nil {
		err = errors.Errorf("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	ids, pager = pagerExtra(res)
	return
}

// Publish get search publish.
func (d *Dao) Publish(c context.Context, arg *blocked.ArgPublishSearch) (ids []int64, pager *blocked.Pager, err error) {
	req := d.elastic.NewRequest(blocked.BusinessBlockedPublish).Index(blocked.TableBlockedPublish).Fields("id")
	if arg.Keyword != blocked.SearchDefaultString {
		req.WhereLike([]string{"title", "sub_title"}, []string{arg.Keyword}, true, elastic.LikeLevelHigh)
	}
	if arg.Type != blocked.SearchDefaultNum {
		req.WhereEq("ptype", arg.Type)
	}
	req.WhereRange("show_time", arg.ShowFrom, arg.ShowTo, elastic.RangeScopeLcRc)
	req.WhereEq("status", blocked.SearchDefaultStatus)
	req.Pn(arg.PN).Ps(arg.PS).Order(arg.Order, arg.Sort)
	var res *search.ReSearchData
	if err = req.Scan(c, &res); err != nil {
		err = errors.Errorf("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	ids, pager = pagerExtra(res)
	return
}

// Case get search case.
func (d *Dao) Case(c context.Context, arg *blocked.ArgCaseSearch) (ids []int64, pager *blocked.Pager, err error) {
	req := d.elastic.NewRequest(blocked.BusinessBlockedCase).Index(blocked.TableBlockedCase).Fields("id")
	if arg.Keyword != blocked.SearchDefaultString {
		req.WhereLike([]string{"origin_content"}, []string{arg.Keyword}, true, elastic.LikeLevelHigh)
	}
	if arg.OriginType != blocked.SearchDefaultNum {
		req.WhereEq("origin_type", arg.OriginType)
	}
	if arg.Status != blocked.SearchDefaultNum {
		req.WhereEq("status", arg.Status)
	}
	if arg.CaseType != blocked.SearchDefaultNum {
		req.WhereEq("case_type", arg.CaseType)
	}
	if arg.UID != blocked.SearchDefaultNum {
		req.WhereEq("mid", arg.UID)
	}
	if arg.OPID != blocked.SearchDefaultNum {
		req.WhereEq("oper_id", arg.OPID)
	}
	req.WhereRange("start_time", arg.TimeFrom, arg.TimeTo, elastic.RangeScopeLcRc)
	req.Pn(arg.PN).Ps(arg.PS).Order(arg.Order, arg.Sort)
	var res *search.ReSearchData
	if err = req.Scan(c, &res); err != nil {
		err = errors.Errorf("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	ids, pager = pagerExtra(res)
	return
}

// Jury get search jury.
func (d *Dao) Jury(c context.Context, arg *blocked.ArgJurySearch) (ids []int64, pager *blocked.Pager, err error) {
	req := d.elastic.NewRequest(blocked.BusinessBlockedJury).Index(blocked.TableBlockedJury).Fields("id")
	if arg.UID != blocked.SearchDefaultNum {
		req.WhereEq("mid", arg.UID)
	}
	if arg.Status != blocked.SearchDefaultNum {
		req.WhereEq("status", arg.Status)
	}
	if arg.Black != blocked.SearchDefaultNum {
		req.WhereEq("black", arg.Black)
	}
	req.WhereRange("expired", arg.ExpiredFrom, arg.ExpiredTo, elastic.RangeScopeLcRc)
	req.Pn(arg.PN).Ps(arg.PS).Order(arg.Order, arg.Sort)
	var res *search.ReSearchData
	if err = req.Scan(c, &res); err != nil {
		err = errors.Errorf("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	ids, pager = pagerExtra(res)
	return
}

// Opinion get search opinion.
func (d *Dao) Opinion(c context.Context, arg *blocked.ArgOpinionSearch) (ids []int64, pager *blocked.Pager, err error) {
	req := d.elastic.NewRequest(blocked.BusinessBlockedOpinion).Index(blocked.TableBlockedOpinion).Fields("id")
	if arg.UID != blocked.SearchDefaultNum {
		req.WhereEq("mid", arg.UID)
	}
	if arg.CID != blocked.SearchDefaultNum {
		req.WhereEq("cid", arg.CID)
	}
	if arg.Vote != blocked.SearchDefaultNum {
		req.WhereEq("vote", arg.Vote)
	}
	if arg.State != blocked.SearchDefaultNum {
		req.WhereEq("state", arg.State)
	}
	req.Pn(arg.PN).Ps(arg.PS).Order(arg.Order, arg.Sort)
	var res *search.ReSearchData
	if err = req.Scan(c, &res); err != nil {
		err = errors.Errorf("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	ids, pager = pagerExtra(res)
	return
}

// KPIPoint get search kpi point.
func (d *Dao) KPIPoint(c context.Context, arg *blocked.ArgKpiPointSearch) (ids []int64, pager *blocked.Pager, err error) {
	var (
		dayFromT, dayToT time.Time
		dayFrom, dayTo   string
	)
	req := d.elastic.NewRequest(blocked.BusinessBlockedKpiPoint).Index(blocked.TableBlockedKpiPoint).Fields("id")
	if arg.UID != blocked.SearchDefaultNum {
		req.WhereEq("mid", arg.UID)
	}
	if arg.Start != blocked.SearchDefaultString {
		if dayFromT, err = time.ParseInLocation(model.TimeFormatSec, arg.Start, time.Local); err != nil {
			err = errors.Errorf("time.ParseInLocation(%s) error(%v)", arg.Start, err)
			return
		}
		dayFrom = dayFromT.Format(model.TimeFormatDay)
	}
	if arg.End != blocked.SearchDefaultString {
		if dayToT, err = time.ParseInLocation(model.TimeFormatSec, arg.End, time.Local); err != nil {
			err = errors.Errorf("time.ParseInLocation(%s) error(%v)", arg.End, err)
			return
		}
		dayTo = dayToT.Format(model.TimeFormatDay)
	}
	req.WhereRange("day", dayFrom, dayTo, elastic.RangeScopeLcRc)
	req.Pn(arg.PN).Ps(arg.PS).Order(arg.Order, arg.Sort)
	var res *search.ReSearchData
	if err = req.Scan(c, &res); err != nil {
		err = errors.Errorf("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	ids, pager = pagerExtra(res)
	return
}

// SearchUpdate update about seach data.
func (d *Dao) SearchUpdate(c context.Context, appid string, table string, data []interface{}) (err error) {
	us := d.elastic.NewUpdate(appid)
	for _, v := range data {
		us.AddData(table, v)
	}
	if err = us.Do(c); err != nil {
		log.Info("appid(%s) table(%s) params(%s) ip(%s) search blocked update error(%v)", appid, table, us.Params(), metadata.String(c, metadata.RemoteIP), err)
	}
	return
}

func pagerExtra(res *search.ReSearchData) (ids []int64, pager *blocked.Pager) {
	for _, v := range res.Result {
		ids = append(ids, v.ID)
	}
	if res.Page != nil {
		pager = &blocked.Pager{
			Total: res.Page.Total,
			PN:    res.Page.PN,
			PS:    res.Page.PS,
			Sort:  res.Sort,
			Order: res.Order,
		}
	}
	return
}
