package dao

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/tsuna/gohbase/hrpc"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"
	"strconv"
	"strings"
)

var (
	_hbaseTable  = "live:push_blacklist"
	_hbaseFamily = "blacklist"

	errLinkValueSplit = errors.New("link_value split result nil.")
)

// GetBlackList get blacklist from hbase by target id
func (d *Dao) GetBlackList(c context.Context, task *model.ApPushTask) (mids map[int64]bool, err error) {
	var (
		key         string
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, d.blackListHBaseReadTimeout)
		emptyByte   = []byte("")
		fbytes      = []byte(_hbaseFamily)
	)
	defer cancel()
	split := strings.Split(task.LinkValue, ",")
	if split == nil {
		err = errLinkValueSplit
		return
	}
	key = split[0]
	mids = make(map[int64]bool)
	if result, err = d.blackListHBase.GetStr(ctx, _hbaseTable, key); err != nil {
		log.Error("[dao.blacklist|GetBlackList] d.blackListHBase.Get error(%v) querytable(%v), roomid(%s), task(%v)",
			err, _hbaseTable, key, task)
		return
	}
	if result == nil {
		return
	}

	for _, c := range result.Cells {
		if c != nil && bytes.Equal(c.Family, fbytes) && !bytes.Equal(c.Qualifier, emptyByte) {
			uid, e := strconv.ParseInt(string(c.Qualifier), 10, 64)
			if e != nil {
				continue
			}
			mids[uid] = true
		}
	}
	log.Info("[dao.blacklist|GetBlackList] get blacklist(%v), roomid(%s), task(%v)", mids, key, task)
	return
}
