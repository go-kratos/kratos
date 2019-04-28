package calculator

import (
	"flag"
	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"path/filepath"
	"testing"
)

func init() {
	dir, _ := filepath.Abs("../../cmd/upcredit-service.toml")
	flag.Set("conf", dir)
	conf.Init()
}

var (
	logs = []upcrmmodel.SimpleCreditLog{
		{Type: 1, OpType: -10, Reason: 0},
		{Type: 1, OpType: -1, Reason: 0},
		{Type: 1, OpType: -3, Reason: 0},
		{Type: 1, OpType: -1, Reason: 0},
		{Type: 1, OpType: -1, Reason: 0},
		{Type: 1, OpType: 0, Reason: 0},
		{Type: 1, OpType: 0, Reason: 0},
		{Type: 1, OpType: -9, Reason: 0},
		{Type: 1, OpType: 0, Reason: 0},
		{Type: 1, OpType: -30, Reason: 0},
	}
)

func TestArcFSM(t *testing.T) {
	var stat = creditStat{}
	var article = CreateArticleStateMachine(logs[0].OpType, logs[0].Type, logs[0].Reason)
	for i := 1; i < len(logs); i++ {
		article.OnLog(&logs[i], stat.onLogResult)
	}
	stat.CalcRelativeScore()
	stat.CalcTotalScore()
	t.Logf("stat: %+v", stat)
}

func TestArcFsmInitState(t *testing.T) {
	var fsm = CreateArticleStateMachineWithInitState()
	var init = conf.CreditConfig.ArticleRule.InitState
	if fsm.Round != init.Round ||
		fsm.Reason != init.Reason ||
		fsm.State != init.State {
		t.Errorf("fail to pass init state!")
	}
}
