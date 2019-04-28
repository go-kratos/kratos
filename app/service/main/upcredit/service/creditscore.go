package service

import (
	"context"
	"go-common/app/service/main/upcredit/common/fsm"
	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/mathutil"
	"go-common/app/service/main/upcredit/model/calculator"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
	"time"
)

//GetCreditScore get credit score
func (c *Service) GetCreditScore(con context.Context, arg *upcrmmodel.GetScoreParam) (results []*upcrmmodel.UpScoreHistory, err error) {
	results, err = c.upcrmdb.GetCreditScore(con, arg)
	if err != nil {
		log.Error("error when get credit score, arg=%+v, err=%+v", arg, err)
	} else {
		log.Info("get credit score, arg=%+v, lens=%d", arg, len(results))
	}
	return
}

func getOrCreateArticleFSM(aid int, m map[int]*calculator.ArticleStateMachine) (artFsm *calculator.ArticleStateMachine) {
	artFsm, ok := m[aid]
	if !ok {
		artFsm = calculator.CreateArticleStateMachineWithInitState()
		m[aid] = artFsm
	}
	return
}

//GetCreditLog get log
func (c *Service) GetCreditLog(con context.Context, arg *upcrmmodel.ArgGetLogHistory) (results []*upcrmmodel.SimpleCreditLogWithContent, err error) {

	var now = time.Now()

	if arg.FromDate.IsZero() {
		arg.FromDate = now.AddDate(-1, 0, 0)
	}
	if arg.ToDate.IsZero() {
		arg.ToDate = now
	}

	if arg.Limit <= 0 {
		arg.Limit = 20
	}

	results, err = c.upcrmdb.GetCreditLog(con, arg)

	if err != nil {
		log.Error("error when get log, arg=%+v, err=%+v", arg, err)
		return
	}

	log.Info("get original log, arg=%+v, lens=%d", arg, len(results))
	var logListTimeAsc []*upcrmmodel.SimpleCreditLogWithContent
	var machineMap = make(map[int]*calculator.ArticleStateMachine)
	for _, l := range results {
		switch l.BusinessType {
		case upcrmmodel.BusinessTypeArticleAudit:
			afsm := getOrCreateArticleFSM(int(l.Oid), machineMap)
			if afsm == nil {
				log.Error("fail to get article's state machine, check reason, id=%d", l.Oid)
				continue
			}
			afsm.OnLog(&l.SimpleCreditLog, func(e *fsm.Event, l1 *upcrmmodel.SimpleCreditLog, a *calculator.ArticleStateMachine) {
				var needLog = false
				switch e.Dst {
				case calculator.StateClose:
					if conf.CreditConfig.ArticleRule.IsRejected(l.Type, l.OpType, l.Reason) {
						needLog = true
					}
				case calculator.StateOpen:
					needLog = true
					if l.Content == "" {
						l.Content = "稿件开放浏览"
					}
				}
				if needLog {
					logListTimeAsc = append(logListTimeAsc, l)
				}
			})
		default:
			logListTimeAsc = append(logListTimeAsc, l)
		}
	}

	log.Info("get refined log, arg=%+v, lens=%d", arg, len(logListTimeAsc))
	// 转化成按时间降序的
	for i, j := 0, len(logListTimeAsc)-1; i < j; i, j = i+1, j-1 {
		logListTimeAsc[i], logListTimeAsc[j] = logListTimeAsc[j], logListTimeAsc[i]
	}

	var max = mathutil.Min(len(logListTimeAsc), arg.Limit)
	results = logListTimeAsc[0:max]
	return
}
