package dao

import (
	"context"
	"fmt"
	"strconv"

	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_opinion        = "op_v2_"    // user opinion prefix.
	_labourIsAnswer = "labour_%d" // key of labourIsAnswer
	_juryInfo       = "jy_"       // key of jury info
	_caseVoteTop    = "ca_vo_top_%d"
)

// user opinion key.
func opinionKey(opid int64) string {
	return _opinion + strconv.FormatInt(opid, 10)
}

// labourKey.
func labourKey(mid int64) string {
	return fmt.Sprintf(_labourIsAnswer, mid)
}

// user jury info key.
func juryInfoKey(mid int64) string {
	return _juryInfo + strconv.FormatInt(mid, 10)
}

func caseVoteTopKey(mid int64) string {
	return fmt.Sprintf(_caseVoteTop, mid)
}

// DelOpinionCache del opinion cache info.
func (d *Dao) DelOpinionCache(c context.Context, vid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(opinionKey(vid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d) error(%v)", vid, err)
	}
	return
}

// DelAnswerStateCache del answer state cache info.
func (d *Dao) DelAnswerStateCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(labourKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d) error(%v)", mid, err)
	}
	return
}

// DelJuryInfoCache del jury cache info.
func (d *Dao) DelJuryInfoCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(juryInfoKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d) error(%v)", mid, err)
	}
	return
}

// DelCaseVoteTopCache del case vote total cache.
func (d *Dao) DelCaseVoteTopCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(caseVoteTopKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Delete(%d) error(%v)", mid, err)
	}
	return
}
