package dao

import (
	"context"
	"fmt"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixCaseInfo     = "ca_in_%d"
	_prefixVoteCaseInfo = "vc_in_%d_%d"
)

func caseInfoKey(cid int64) string {
	return fmt.Sprintf(_prefixCaseInfo, cid)
}

func voteCaseInfoKey(mid int64, cid int64) string {
	return fmt.Sprintf(_prefixVoteCaseInfo, mid, cid)
}

// DelCaseInfoCache del case info cache info.
func (d *Dao) DelCaseInfoCache(c context.Context, cid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(caseInfoKey(cid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d) error(%v)", cid, err)
	}
	return
}

// DelVoteCaseCache del vote case cache info.
func (d *Dao) DelVoteCaseCache(c context.Context, mid, cid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(voteCaseInfoKey(mid, cid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d,%d) error(%v)", mid, cid, err)
	}
	return
}
