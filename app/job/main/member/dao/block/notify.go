package block

import (
	"context"

	"github.com/pkg/errors"
)

// AccountNotify is
func (d *Dao) AccountNotify(c context.Context, mid int64) (err error) {
	action := "blockUser"
	if err = d.notifyFunc(c, mid, action); err != nil {
		err = errors.Errorf("mid(%d) s.accountNotify.Send(%+v) error(%+v)", mid, action, err)
	}
	return
}
