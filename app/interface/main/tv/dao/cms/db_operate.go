package cms

import (
	"context"
	"fmt"

	"go-common/library/xstr"
)

const (
	_unshelveArcs = "UPDATE ugc_archive SET valid = 0 WHERE aid IN (%s) AND valid = 1 AND deleted = 0"
)

// UnshelveArcs unshelves the arcs
func (d *Dao) UnshelveArcs(c context.Context, ids []int64) (err error) {
	_, err = d.db.Exec(c, fmt.Sprintf(_unshelveArcs, xstr.JoinInts(ids)))
	return
}
