package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/answer/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHistoryES(t *testing.T) {
	Convey("HistoryES", t, func() {
		d.HistoryES(context.Background(), &model.ArgHistory{Mid: 14771787, Pn: 1, Ps: 1})
	})
}
