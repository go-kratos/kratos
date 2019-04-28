package dao

import (
	"context"
	"fmt"
	"strconv"

	model "go-common/app/interface/main/credit/model"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_noticeInfo     = "notice"    // key of creditInfo
	_reasonList     = "reason"    // key of creditInfo
	_opinion        = "op_v2_"    // user opinion prefix.
	_questionids    = "question"  // key of creditInfo
	_labourIsAnswer = "labour_%d" // key of labourIsAnswer
	_juryInfo       = "jy_"       // key of jury info
)

func noticeKey() string {
	return _noticeInfo
}

func reasonKey() string {
	return _reasonList
}

// user opinion key.
func opinionKey(opid int64) string {
	return _opinion + strconv.FormatInt(opid, 10)
}

func questionKey(mid int64) string {
	return _questionids + strconv.FormatInt(mid, 10)
}

func labourKey(mid int64) string {
	return fmt.Sprintf(_labourIsAnswer, mid)
}

// user jury info key.
func juryInfoKey(mid int64) string {
	return _juryInfo + strconv.FormatInt(mid, 10)
}

func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Set(&gmc.Item{Key: "ping", Value: []byte{1}, Expiration: d.commonExpire}); err != nil {
		log.Error("conn.Store(set, ping, 1) error(%v)", err)
	}
	conn.Close()
	return
}

// NoticeInfoCache get notice info cache.
func (d *Dao) NoticeInfoCache(c context.Context) (dt *model.Notice, err error) {
	var conn = d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(noticeKey())
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	dt = &model.Notice{}
	err = conn.Scan(item, &dt)
	return
}

// SetReasonListCache get notice info cache.
func (d *Dao) SetReasonListCache(c context.Context, dt []*model.Reason) (err error) {
	var (
		item = &gmc.Item{Key: reasonKey(), Object: dt, Expiration: d.minCommonExpire, Flags: gmc.FlagJSON}
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// SetNoticeInfoCache get Notice info cache.
func (d *Dao) SetNoticeInfoCache(c context.Context, n *model.Notice) (err error) {
	var (
		item = &gmc.Item{Key: noticeKey(), Object: n, Expiration: d.minCommonExpire, Flags: gmc.FlagJSON}
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// DelNoticeInfoCache delete notice info cache.
func (d *Dao) DelNoticeInfoCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(noticeKey()); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
	}
	return
}

// ReasonListCache get Reason info cache.
func (d *Dao) ReasonListCache(c context.Context) (dt []*model.Reason, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(reasonKey())
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	dt = make([]*model.Reason, 0)
	err = conn.Scan(item, &dt)
	return
}

// DelReasonListCache delete reason info cache.
func (d *Dao) DelReasonListCache(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(reasonKey()); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
	}
	return
}

// AddOpinionCache add opinion cache.
func (d *Dao) AddOpinionCache(c context.Context, op *model.Opinion) (err error) {
	var (
		item = &gmc.Item{Key: opinionKey(op.OpID), Object: op, Expiration: d.commonExpire, Flags: gmc.FlagJSON}
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// OpinionsCache get opids opinion detail.
func (d *Dao) OpinionsCache(c context.Context, opids []int64) (ops map[int64]*model.Opinion, miss []int64, err error) {
	var (
		keys []string
		conn = d.mc.Get(c)
	)
	for _, id := range opids {
		keys = append(keys, opinionKey(id))
	}
	defer conn.Close()
	res, err := conn.GetMulti(keys)
	if err != nil {
		return
	}
	ops = make(map[int64]*model.Opinion, len(res))
	for _, id := range opids {
		if r, ok := res[opinionKey(id)]; ok {
			op := &model.Opinion{}
			if err = conn.Scan(r, &op); err != nil {
				return
			}
			ops[op.OpID] = op
		} else {
			miss = append(miss, id)
		}
	}
	return
}

// SetQsCache set labour questions cache.
func (d *Dao) SetQsCache(c context.Context, mid int64, qsid *model.QsCache) (err error) {
	var (
		item = &gmc.Item{Key: questionKey(mid), Object: qsid, Expiration: d.commonExpire, Flags: gmc.FlagJSON}
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// GetQsCache get labour qestions info cache.
func (d *Dao) GetQsCache(c context.Context, mid int64) (qsid *model.QsCache, err error) {
	var (
		reply *gmc.Item
		conn  = d.mc.Get(c)
	)
	defer conn.Close()
	reply, err = conn.Get(questionKey(mid))
	if err != nil || reply == nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	qsid = &model.QsCache{}
	err = conn.Scan(reply, &qsid)
	return
}

// DelQsCache delete labour questions cache.
func (d *Dao) DelQsCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(questionKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
	}
	return
}

// SetAnswerStateCache set labour answer state cache.
func (d *Dao) SetAnswerStateCache(c context.Context, mid int64, state int8) (err error) {
	var (
		conn = d.mc.Get(c)
		item = &gmc.Item{Key: labourKey(mid), Object: state, Expiration: d.userExpire, Flags: gmc.FlagJSON}
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// GetAnswerStateCache get labour answer state cache.
func (d *Dao) GetAnswerStateCache(c context.Context, mid int64) (state int8, found bool, err error) {
	var (
		reply *gmc.Item
		conn  = d.mc.Get(c)
	)
	defer conn.Close()
	if reply, err = conn.Get(labourKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	if err = conn.Scan(reply, &state); err != nil {
		return
	}
	found = true
	return
}

// JuryInfoCache .
func (d *Dao) JuryInfoCache(c context.Context, mid int64) (bj *model.BlockedJury, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(juryInfoKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	bj = &model.BlockedJury{}
	err = conn.Scan(item, &bj)
	return
}

// SetJuryInfoCache .
func (d *Dao) SetJuryInfoCache(c context.Context, mid int64, bj *model.BlockedJury) (err error) {
	var (
		conn = d.mc.Get(c)
		item = &gmc.Item{Key: juryInfoKey(mid), Object: bj, Expiration: d.userExpire, Flags: gmc.FlagJSON}
	)
	defer conn.Close()
	err = conn.Set(item)
	return
}

// DelJuryInfoCache del jury cache info.
func (d *Dao) DelJuryInfoCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Delete(juryInfoKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
	}
	return
}
