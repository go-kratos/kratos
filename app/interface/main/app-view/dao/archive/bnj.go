package archive

import (
	"context"

	"go-common/app/service/main/archive/model/archive"
)

// ViewsRPC is
func (d *Dao) ViewsRPC(c context.Context, aids []int64) (avm map[int64]*archive.View3, err error) {
	avm, err = d.arcRPC.Views3(c, &archive.ArgAids2{Aids: aids})
	return
}

// ViewRPC is
func (d *Dao) ViewRPC(c context.Context, aid int64) (v *archive.View3, err error) {
	v, err = d.arcRPC.View3(c, &archive.ArgAid2{Aid: aid})
	return
}
