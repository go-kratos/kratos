package service

import (
	"flag"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/openplatform/anti-fraud/conf"

	"context"
	"go-common/app/service/openplatform/anti-fraud/model"
	"testing"
)

var s *Service

func init() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}

	s = New(conf.Conf)
}

func TestGetQuestionBankBind(t *testing.T) {
	Convey("TestGetQuestionBankBind: ", t, func() {
		args := &model.ArgGetBankBind{
			TargetItems:    []string{"asasdds", "asasdds1", "2"},
			TargetItemType: 1,
			Source:         33,
		}
		res, err := s.GetQuestionBankBind(context.TODO(), args)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
