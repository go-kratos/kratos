package upcrmservice

import (
	"context"
	"fmt"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/log"
)

func formatLog(log upcrmmodel.SimpleCreditLogWithContent) (result upcrmmodel.CreditLogInfo) {

	var timeStr = log.CTime.Time().Format(upcrmmodel.TimeFmtDate)
	switch log.BusinessType {
	case upcrmmodel.BusinessTypeArticleAudit:
		result.Log = fmt.Sprintf("[%s][aid=%d]%s", timeStr, log.Oid, log.Content)
	default:
		result.Log = fmt.Sprintf("[%s]%s", timeStr, log.Content)
	}

	result.Time = log.CTime
	return
}

//CreditLogQueryUp query credit log
func (s *Service) CreditLogQueryUp(context context.Context, arg *upcrmmodel.CreditLogQueryArgs) (result upcrmmodel.CreditLogUpResult, err error) {
	if arg.Limit <= 0 {
		arg.Limit = 20
	} else if arg.Limit >= 100 {
		arg.Limit = 100
	}

	logs, e := s.crmdb.GetCreditLog(arg.Mid, arg.Limit)
	err = e
	if err != nil {
		log.Error("credit log get fail, err=%+v", err)
		return
	}
	for _, v := range logs {
		//if v == nil {
		//	continue
		//}
		result.Logs = append(result.Logs, formatLog(v))
	}
	log.Info("credit log get ok, mid=%d, length=%d", arg.Mid, len(logs))
	return
}
