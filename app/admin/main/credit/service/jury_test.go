package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/credit/model/blocked"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SetCaseConf(t *testing.T) {
	Convey("return someting", t, func() {
		cc := new(blocked.ArgCaseConf)
		cc.CaseCheckHours = 2
		cc.CaseGiveHours = 24
		cc.CaseJudgeRadio = 60
		cc.CaseLoadMax = 20
		cc.CaseLoadSwitch = 1
		cc.CaseObtainMax = 100
		cc.CaseVoteMax = 50
		cc.CaseVoteMin = 50
		cc.JuryApplyMax = 200
		cc.JuryVoteRadio = 40
		err := s.SetCaseConf(context.TODO(), cc)
		fmt.Println(cc)
		fmt.Println(err)
		So(err, ShouldBeNil)
		So(cc, ShouldNotBeEmpty)
	})
}

func Test_CaseConfig(t *testing.T) {
	Convey("return someting", t, func() {
		cv := s.CaseConfig(blocked.ConfigCaseGiveHours)
		fmt.Println(cv)
		So(cv, ShouldNotBeEmpty)
	})
}
func Test_VotenumConf(t *testing.T) {
	Convey("return someting", t, func() {
		cc, err := s.VotenumConf(context.TODO())
		fmt.Println(cc)
		fmt.Println(err)
		So(err, ShouldBeNil)
		So(cc, ShouldNotBeEmpty)
	})
}

func Test_SetVotenumConf(t *testing.T) {
	Convey("return someting", t, func() {
		vn := new(blocked.ArgVoteNum)
		vn.OID = 1
		vn.RateS = 3
		vn.RateA = 2
		vn.RateB = 1
		vn.RateC = 1
		vn.RateD = 1
		err := s.SetVotenumConf(context.TODO(), vn)
		So(err, ShouldBeNil)
	})
}
