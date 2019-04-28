package service

import (
	"context"
	"testing"

	"go-common/app/job/main/credit/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_addReplyReport(t *testing.T) {
	var (
		c  = context.TODO()
		mr = &model.Reply{
			Subject: &model.ReplySubject{
				OID: 2222,
			},
			Reply: &model.ReplyMain{
				Like: 1,
				MID:  27515628,
				Content: &struct {
					Message string `json:"message"`
				}{Message: "xxxx"},
			},
			Report: &model.ReplyReport{
				State:  model.ReportStateNew,
				Type:   model.SubTypeArchive,
				Reason: model.ReplyReasonGarbageAds,
				Score:  2,
				RPID:   1111,
				OID:    2222,
			},
		}
	)
	Convey("should return err be nil", t, func() {
		err := s.addReplyReport(c, mr)
		So(err, ShouldBeNil)
	})
}
