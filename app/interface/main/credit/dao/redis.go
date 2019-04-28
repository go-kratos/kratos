package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_caseObtainMIDIdx  = "c_o_v2_%d"
	_caseVoteCIDMIDIdx = "c_v_c_m_v2_%d_%s"
	_grantCaseKey      = "gr_ca_li_v2"
)

func caseObtainMIDKey(mid int64) string {
	return fmt.Sprintf(_caseObtainMIDIdx, mid)
}

func caseVoteCIDMIDKey(mid int64) string {
	return fmt.Sprintf(_caseVoteCIDMIDIdx, mid, time.Now().Format("2006_01_02"))
}

// CaseObtainMID get voted cids by MID.
func (d *Dao) CaseObtainMID(c context.Context, mid int64, isToday bool) (cases map[int64]*model.SimCase, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var _setKey string
	if isToday {
		_setKey = caseVoteCIDMIDKey(mid)
	} else {
		_setKey = caseObtainMIDKey(mid)
	}
	var ms []string
	if ms, err = redis.Strings(conn.Do("SMEMBERS", _setKey)); err != nil {
		if err != redis.ErrNil {
			return
		}
		err = nil
	}
	cases = make(map[int64]*model.SimCase)
	for _, s := range ms {
		if s == "" {
			continue
		}
		sc := &model.SimCase{}
		if err = json.Unmarshal([]byte(s), sc); err != nil {
			err = errors.WithStack(err)
			return
		}
		cases[sc.ID] = sc
	}
	return
}

// IsExpiredObtainMID voted cids by MID of key is expired.
func (d *Dao) IsExpiredObtainMID(c context.Context, mid int64, isToday bool) (ok bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var _setKey string
	if isToday {
		_setKey = caseVoteCIDMIDKey(mid)
	} else {
		_setKey = caseObtainMIDKey(mid)
	}
	ok, err = redis.Bool(conn.Do("EXPIRE", _setKey, d.redisExpire))
	return
}

// GrantCases get granting case cids.
func (d *Dao) GrantCases(c context.Context) (mcases map[int64]*model.SimCase, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var ms map[string]string
	if ms, err = redis.StringMap(conn.Do("HGETALL", _grantCaseKey)); err != nil {
		if err != redis.ErrNil {
			return
		}
		err = nil
	}
	mcases = make(map[int64]*model.SimCase)
	for m, s := range ms {
		if s == "" {
			continue
		}
		cid, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", m, err)
			err = nil
			continue
		}
		cc := &model.SimCase{}
		if err = json.Unmarshal([]byte(s), cc); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", s, err)
			err = nil
			continue
		}
		mcases[cid] = cc
	}
	return
}

// LoadVoteCaseMID load mid vote cases.
func (d *Dao) LoadVoteCaseMID(c context.Context, mid int64, mcases map[int64]*model.SimCase, isToday bool) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var _setKey string
	if isToday {
		_setKey = caseVoteCIDMIDKey(mid)
	} else {
		_setKey = caseObtainMIDKey(mid)
	}
	args := redis.Args{}.Add(_setKey)
	for _, mcase := range mcases {
		var mcStr []byte
		if mcStr, err = json.Marshal(mcase); err != nil {
			err = errors.WithStack(err)
			return
		}
		args = args.Add(string(mcStr))
	}
	if err = conn.Send("SADD", args...); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", _setKey, d.redisExpire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < 2; i++ {
		_, err = conn.Receive()
	}
	return
}

// SetVoteCaseMID set vote case by mid on today.
func (d *Dao) SetVoteCaseMID(c context.Context, mid int64, sc *model.SimCase) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var scStr []byte
	if scStr, err = json.Marshal(sc); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = conn.Send("SADD", caseObtainMIDKey(mid), string(scStr)); err != nil {
		return
	}
	if err = conn.Send("SADD", caseVoteCIDMIDKey(mid), string(scStr)); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", caseObtainMIDKey(mid), d.redisExpire); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", caseVoteCIDMIDKey(mid), d.redisExpire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	for i := 0; i < 4; i++ {
		_, err = conn.Receive()
	}
	return
}
