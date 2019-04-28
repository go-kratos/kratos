package tag

import (
	"context"
	"testing"

	tag "go-common/app/interface/main/tag/model"
)

var (
	svr *Service
)

func TestTagRPC(t *testing.T) {
	var (
		ctx = context.TODO()
	)
	svr = New2(nil)
	arg := &tag.ArgBind{Oid: 1, Mid: 12345, Type: tag.PicResType, Names: []string{"User"}}
	t.Log(svr.UpBind(ctx, arg))
}
