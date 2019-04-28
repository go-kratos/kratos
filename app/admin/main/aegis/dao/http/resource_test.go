package http

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestSyncResource(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("SyncResource", t, func(ctx convey.C) {
		d.SyncResource(c, &model.Action{}, map[string]interface{}{"state": 1})
	})
}
