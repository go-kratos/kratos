package dao

import (
	"context"

	"go-common/app/job/main/workflow/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

const (
	challBusiness = "workflow_chall_common"
	challIndex    = "workflow_chall_common"

	appealBusiness = "workflow_appeal"
	appealIndex    = "workflow_appeal"
)

// SearchChall .
func (d *Dao) SearchChall(c context.Context, params *model.SearchParams) (res *model.ChallSearchCommonRes, err error) {
	r := d.es.NewRequest(challBusiness).Fields("id").
		Index(challIndex).
		WhereEq("business", params.Business).
		WhereEq("state", params.States).
		WhereEq("business_state", params.BusinessStates)
	if params.AssigneeAdminIDs != "" {
		r = r.WhereEq("assignee_adminid", params.AssigneeAdminIDs)
	}
	if params.AssigneeAdminIDsNot != "" {
		r = r.WhereNot("assignee_adminid", params.AssigneeAdminIDsNot)
	}
	if params.MtimeTo != "" {
		r = r.WhereRange("mtime", "", params.MtimeTo, elastic.RangeScopeLcRo)
	}
	log.Info("search condition %v", r.Params())
	err = r.Scan(c, &res)
	return
}

// SearchAppeal .
func (d *Dao) SearchAppeal(c context.Context, cond model.AppealSearchCond) (res *model.AppealSearchRes, err error) {
	r := d.es.NewRequest(appealBusiness).Index(appealIndex).Fields(cond.Fields...).WhereIn("bid", cond.Bid).
		WhereIn("id", cond.IDs).WhereIn("assign_state", cond.AssignState).WhereIn("audit_state", cond.AuditState).
		WhereIn("transfer_state", cond.TransferState).WhereIn("audit_adminid", cond.AuditAdmin).
		WhereIn("transfer_adminid", cond.TransferAdmin).WhereRange("ctime", cond.CTimeFrom, cond.CTimeTo, elastic.RangeScopeLcRo).
		WhereRange("dtime", cond.DTimeFrom, cond.DTimeTo, elastic.RangeScopeLcRo).WhereRange("ttime", cond.TTimeFrom, cond.TTimeTo, elastic.RangeScopeLcRo).
		WhereRange("mtime", cond.MTimeFrom, cond.MTimeTo, elastic.RangeScopeLcRo).Order(cond.Order, cond.Sort).Pn(cond.PN).Ps(cond.PS)

	log.Info("search condition (%v)", r.Params())
	if err = r.Scan(c, &res); err != nil {
		log.Error("r.Scan(%+v) error(%v)", &res, err)
	}
	if res == nil {
		res = new(model.AppealSearchRes)
	}
	return
}
