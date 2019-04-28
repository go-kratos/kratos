package calculator

import (
	"fmt"
	"go-common/app/service/main/upcredit/common/fsm"
	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
)

//Rule rule function
type Rule func(log upcrmmodel.SimpleCreditLog, stat *creditStat, args ...interface{})

var (
	//RuleList rule list
	RuleList = []Rule{rule1}
)

const (
	//StateOpen open
	StateOpen = "st_open"
	//StateClose close
	StateClose = "st_close"

	//EventOpen ev open
	EventOpen = "open"
	//EventClose ev close
	EventClose = "close"
)

//ArticleStateMachine state machine
type ArticleStateMachine struct {
	State     int
	Round     int
	Reason    int
	openCount int // 记录一共open了多少次
	Fsm       *fsm.FSM
}

// ----------------state machine
func (a *ArticleStateMachine) onEnterOpen(e *fsm.Event) {
	var l = e.Args[0].(*upcrmmodel.SimpleCreditLog)
	if l == nil {
		log.Error("event arg is not upcrmmodel.SimpleCreditLog, please check code")
		return
	}

	if a.openCount >= conf.CreditConfig.ArticleRule.ArticleMaxOpenCount {
		fmt.Printf("article max open count exceeded")
		return
	}
	a.openCount++
	a.setState(l)
	var callback, ok = e.Args[1].(OnLogResult)
	if ok && callback != nil {
		(callback)(e, l, a)
	}
	fmt.Printf("article opened\n")
}

func (a *ArticleStateMachine) onEnterClose(e *fsm.Event) {
	var l = e.Args[0].(*upcrmmodel.SimpleCreditLog)
	if l == nil {
		log.Error("event arg is not upcrmmodel.SimpleCreditLog, please check code")
		return
	}
	if a.State == l.OpType && a.Round == l.Type {
		fmt.Printf("nothing happend here! state=%d, round=%d\n", a.State, a.Round)
		return
	}

	a.setState(l)

	var callback, ok = e.Args[1].(OnLogResult)
	if ok && callback != nil {
		callback(e, l, a)
	}
}

//OnLogResult on log function
type OnLogResult func(e *fsm.Event, l *upcrmmodel.SimpleCreditLog, a *ArticleStateMachine)

// OnLog -----------------state machine
func (a *ArticleStateMachine) OnLog(l *upcrmmodel.SimpleCreditLog, callback OnLogResult) {
	var event = ""
	if l.OpType == 0 {
		event = EventOpen
	} else if l.OpType < 0 {
		event = EventClose
	}

	if event != "" {
		a.Fsm.Event(event, l, callback)
	}
}

func (a *ArticleStateMachine) setState(l *upcrmmodel.SimpleCreditLog) {
	a.State = l.OpType
	a.Round = l.Type
	a.Reason = l.Reason
}

//CreateArticleStateMachineWithInitState create init machine
func CreateArticleStateMachineWithInitState() (asm *ArticleStateMachine) {
	var initState = &conf.CreditConfig.ArticleRule.InitState
	return CreateArticleStateMachine(initState.State, initState.Round, initState.Reason)
}

// CreateArticleStateMachine 这里的state, round, reason是稿件的state, round, reason
func CreateArticleStateMachine(state int, round int, reason int) (asm *ArticleStateMachine) {
	asm = &ArticleStateMachine{
		State:  state,
		Round:  round,
		Reason: reason,
	}
	var initstate = StateOpen
	if state < 0 {
		initstate = StateClose
	}

	asm.Fsm = fsm.NewFSM(
		initstate,
		fsm.Events{
			{Name: EventOpen, Src: []string{StateClose}, Dst: StateOpen},
			{Name: EventClose, Src: []string{StateOpen}, Dst: StateClose},
			{Name: EventClose, Src: []string{StateClose}, Dst: StateClose},
		},
		fsm.Callbacks{
			"enter_" + StateOpen:  func(e *fsm.Event) { asm.onEnterOpen(e) },
			"enter_" + StateClose: func(e *fsm.Event) { asm.onEnterClose(e) },
		},
	)
	return asm
}

func rule1(l upcrmmodel.SimpleCreditLog, stat *creditStat, args ...interface{}) {
	switch l.BusinessType {
	case upcrmmodel.BusinessTypeArticleAudit:
		if len(args) >= 1 {
			machineMap, ok := args[0].(map[int]*ArticleStateMachine)
			if !ok {
				log.Error("fail to get machineMap")
			}
			articleFsm := getOrCreateArticleFSM(int(l.Oid), machineMap)

			if articleFsm != nil {
				articleFsm.OnLog(&l, stat.onLogResult)
			}
		}
	}

}
