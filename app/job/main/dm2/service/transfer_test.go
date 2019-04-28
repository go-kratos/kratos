package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/job/main/dm2/model"
)

func TestTransferDMS(t *testing.T) {
	Convey("test NewCommentList", t, func() {
		ll, err := svr.transferDMS(context.TODO(), 1, 1012, 0, 10)
		So(err, ShouldBeNil)
		So(ll, ShouldNotBeEmpty)
	})
}

func TestTransfer(t *testing.T) {
	trans := &model.Transfer{
		ID:      265,
		FromCid: 1012,
		ToCid:   1211,
		Mid:     0,
		Dmid:    123,
		Offset:  0,
		State:   0,
	}
	Convey("transfer", t, func() {
		svr.transfer(context.TODO(), trans)
	})
}
