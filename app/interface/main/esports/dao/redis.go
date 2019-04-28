package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/interface/main/esports/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_keyCont    = "c_%d"
	_keyVideo   = "v_%d"
	_keyCale    = "c_%s_%s"
	_keyFMat    = "fm"
	_keyFVideo  = "fv"
	_keySeason  = "se"
	_keySeasonM = "sp"
	_keyC       = "co_%d"
	_keyCF      = "cf_%d"
	_keyCAF     = "caf_%d"
	_keyS       = "s_%d"
	_keyAct     = "act_%d"
	_keyModule  = "module_%d"
	_keyTop     = "top_%d_%d"
	_keyPoint   = "point_%d_%d_%d"
	_keyKnock   = "knock_%d"
	_keyMAct    = "ma_%d"
	_keyTeam    = "team_%d"
	_keyCSData  = "c_s_data_%d"
	_keyCRecent = "c_recent_%d_%d_%d_%d"
)

func keyCale(stime, etime string) string {
	return fmt.Sprintf(_keyCale, stime, etime)
}

func keyCont(ps int) string {
	return fmt.Sprintf(_keyCont, ps)
}

func keyVideo(ps int) string {
	return fmt.Sprintf(_keyVideo, ps)
}
func keyContID(cid int64) string {
	return fmt.Sprintf(_keyC, cid)
}

func keyCoFav(mid int64) string {
	return fmt.Sprintf(_keyCF, mid)
}
func keyCoAppFav(mid int64) string {
	return fmt.Sprintf(_keyCAF, mid)
}

func keySID(sid int64) string {
	return fmt.Sprintf(_keyS, sid)
}

func keyTeamID(tid int64) string {
	return fmt.Sprintf(_keyTeam, tid)
}

func keyMatchAct(aid int64) string {
	return fmt.Sprintf(_keyAct, aid)
}

func keyCSData(cid int64) string {
	return fmt.Sprintf(_keyCSData, cid)
}

func keyCRecent(param *model.ParamCDRecent) string {
	key := fmt.Sprintf(_keyCRecent, param.CID, param.HomeID, param.AwayID, param.Ps)
	return key
}

func keyMatchModule(mmid int64) string {
	return fmt.Sprintf(_keyModule, mmid)
}

func keyKnock(mdID int64) string {
	return fmt.Sprintf(_keyKnock, mdID)
}

func keyTop(aid, ps int64) string {
	return fmt.Sprintf(_keyTop, aid, ps)
}

func keyPoint(aid, mdID, ps int64) string {
	return fmt.Sprintf(_keyPoint, aid, mdID, ps)
}

func keyMAct(aid int64) string {
	return fmt.Sprintf(_keyMAct, aid)
}

// FMatCache get filter match from cache.
func (d *Dao) FMatCache(c context.Context) (res map[string][]*model.Filter, err error) {
	res, err = d.filterCache(c, _keyFMat)
	return
}

// FVideoCache get filter video from cache.
func (d *Dao) FVideoCache(c context.Context) (res map[string][]*model.Filter, err error) {
	res, err = d.filterCache(c, _keyFVideo)
	return
}

func (d *Dao) filterCache(c context.Context, key string) (rs map[string][]*model.Filter, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var values []byte
	if values, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			log.Error("filterCache (%s) return nil ", key)
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	rs = make(map[string][]*model.Filter)
	if err = json.Unmarshal(values, &rs); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", values, err)
	}
	return
}

// SetFMatCache set  filter match to cache.
func (d *Dao) SetFMatCache(c context.Context, fs map[string][]*model.Filter) (err error) {
	err = d.setFilterCache(c, _keyFMat, fs)
	return
}

// SetFVideoCache set  filter match to cache.
func (d *Dao) SetFVideoCache(c context.Context, fs map[string][]*model.Filter) (err error) {
	err = d.setFilterCache(c, _keyFVideo, fs)
	return
}

func (d *Dao) setFilterCache(c context.Context, key string, fs map[string][]*model.Filter) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	var bs []byte
	if bs, err = json.Marshal(fs); err != nil {
		log.Error("json.Marshal(%v) error(%v)", fs, err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET,%s,%s) error(%v)", key, string(bs), err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.filterExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s,%d) error(%v)", key, d.filterExpire, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Recevie(%d) error(%v0", i, err)
		}
	}
	return
}

// ContestCache get all contest from cache.
func (d *Dao) ContestCache(c context.Context, ps int) (res []*model.Contest, total int, err error) {
	key := keyCont(ps)
	res, total, err = d.cosCache(c, key)
	return
}

// FavCoCache get fav contest from cache.
func (d *Dao) FavCoCache(c context.Context, mid int64) (res []*model.Contest, total int, err error) {
	key := keyCoFav(mid)
	res, total, err = d.cosCache(c, key)
	return
}

// FavCoAppCache get fav contest from cache.
func (d *Dao) FavCoAppCache(c context.Context, mid int64) (res []*model.Contest, total int, err error) {
	key := keyCoAppFav(mid)
	res, total, err = d.cosCache(c, key)
	return
}

func (d *Dao) cosCache(c context.Context, key string) (res []*model.Contest, total int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		cont := &model.Contest{}
		if err = json.Unmarshal(bs, cont); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, cont)
	}
	total = from(num)
	return
}

// SetContestCache set  all contest to cache.
func (d *Dao) SetContestCache(c context.Context, ps int, contests []*model.Contest, total int) (err error) {
	key := keyCont(ps)
	err = d.setCosCache(c, key, contests, total)
	return
}

// SetFavCoCache set  fav contest to cache.
func (d *Dao) SetFavCoCache(c context.Context, mid int64, contests []*model.Contest, total int) (err error) {
	key := keyCoFav(mid)
	err = d.setCosCache(c, key, contests, total)
	return
}

// SetAppFavCoCache set  fav contest to cache.
func (d *Dao) SetAppFavCoCache(c context.Context, mid int64, contests []*model.Contest, total int) (err error) {
	key := keyCoAppFav(mid)
	err = d.setCosCache(c, key, contests, total)
	return
}

// DelFavCoCache delete fav contests cache.
func (d *Dao) DelFavCoCache(c context.Context, mid int64) (err error) {
	key := keyCoFav(mid)
	keyApp := keyCoAppFav(mid)
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL plaKey(%s) error(%v))", key, err)
		return
	}
	if err = conn.Send("DEL", keyApp); err != nil {
		log.Error("conn.Send(DEL pladKey(%s) error(%v))", keyApp, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush() error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

func (d *Dao) setCosCache(c context.Context, key string, contests []*model.Contest, total int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	args := redis.Args{}.Add(key)
	for sort, contest := range contests {
		bs, _ := json.Marshal(contest)
		args = args.Add(combine(int64(sort), total)).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("conn.Send(ZADD, %s, %v) error(%v)", key, args, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.filterExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.filterExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// CalendarCache get all calendar from cache.
func (d *Dao) CalendarCache(c context.Context, p *model.ParamFilter) (res []*model.Calendar, err error) {
	var (
		key  = keyCale(p.Stime, p.Etime)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		cale := &model.Calendar{}
		if err = json.Unmarshal(bs, cale); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, cale)
	}
	return
}

// SetCalendarCache set  all calendar to cache.
func (d *Dao) SetCalendarCache(c context.Context, p *model.ParamFilter, cales []*model.Calendar) (err error) {
	var (
		key  = keyCale(p.Stime, p.Etime)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	args := redis.Args{}.Add(key)
	for sort, cale := range cales {
		bs, _ := json.Marshal(cale)
		args = args.Add(sort).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("conn.Send(ZADD, %s, %v) error(%v)", key, args, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.filterExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.filterExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// VideoCache get all video from cache.
func (d *Dao) VideoCache(c context.Context, ps int) (res []*arcmdl.Arc, total int, err error) {
	var (
		key  = keyVideo(ps)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		cont := &arcmdl.Arc{}
		if err = json.Unmarshal(bs, cont); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, cont)
	}
	total = from(num)
	return
}

// SetVideoCache set  all contest to cache.
func (d *Dao) SetVideoCache(c context.Context, ps int, videos []*arcmdl.Arc, total int) (err error) {
	var (
		key  = keyVideo(ps)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	args := redis.Args{}.Add(key)
	for sort, video := range videos {
		bs, _ := json.Marshal(video)
		args = args.Add(combine(int64(sort), total)).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("conn.Send(ZADD, %s, %v) error(%v)", key, args, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.filterExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.filterExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}
func (d *Dao) seasonsCache(c context.Context, key string, start, end int) (res []*model.Season, total int, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, start, end, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		object := &model.Season{}
		if err = json.Unmarshal(bs, object); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, object)
	}
	total = from(num)
	return
}

func (d *Dao) setSeasonsCache(c context.Context, key string, seasons []*model.Season, total int) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	for sort, season := range seasons {
		bs, _ := json.Marshal(season)
		if err = conn.Send("ZADD", key, combine(int64(sort), total), bs); err != nil {
			log.Error("conn.Send(ZADD, %s, %s) error(%v)", key, string(bs), err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("conn.Send(Expire, %s, %d) error(%v)", key, d.listExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// SeasonCache get season list from cache.
func (d *Dao) SeasonCache(c context.Context, start, end int) (res []*model.Season, total int, err error) {
	res, total, err = d.seasonsCache(c, _keySeason, start, end)
	return
}

// SetSeasonCache set season list cache.
func (d *Dao) SetSeasonCache(c context.Context, seasons []*model.Season, total int) (err error) {
	err = d.setSeasonsCache(c, _keySeason, seasons, total)
	return
}

// SeasonMCache get season list from cache.
func (d *Dao) SeasonMCache(c context.Context, start, end int) (res []*model.Season, total int, err error) {
	res, total, err = d.seasonsCache(c, _keySeasonM, start, end)
	return
}

// SetSeasonMCache set season list cache.
func (d *Dao) SetSeasonMCache(c context.Context, seasons []*model.Season, total int) (err error) {
	err = d.setSeasonsCache(c, _keySeasonM, seasons, total)
	return
}

func from(i int64) int {
	return int(i & 0xffff)
}

func combine(sort int64, count int) int64 {
	return sort<<16 | int64(count)
}

// CacheEpContests .
func (d *Dao) CacheEpContests(c context.Context, ids []int64) (res map[int64]*model.Contest, err error) {
	var (
		key  string
		args = redis.Args{}
		bss  [][]byte
	)
	for _, cid := range ids {
		key = keyContID(cid)
		args = args.Add(key)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if bss, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("CacheEpContests conn.Do(MGET,%s) error(%v)", key, err)
		}
		return
	}
	res = make(map[int64]*model.Contest, len(ids))
	for _, bs := range bss {
		con := new(model.Contest)
		if bs == nil {
			continue
		}
		if err = json.Unmarshal(bs, con); err != nil {
			log.Error("CacheEpContests json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		res[con.ID] = con
	}
	return
}

// AddCacheEpContests .
func (d *Dao) AddCacheEpContests(c context.Context, data map[int64]*model.Contest) (err error) {
	if len(data) == 0 {
		return
	}
	var (
		bs      []byte
		keyID   string
		keyIDs  []string
		argsCid = redis.Args{}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, v := range data {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("json.Marshal err(%v)", err)
			continue
		}
		keyID = keyContID(v.ID)
		keyIDs = append(keyIDs, keyID)
		argsCid = argsCid.Add(keyID).Add(string(bs))
	}
	if err = conn.Send("MSET", argsCid...); err != nil {
		log.Error("AddCacheMatchSubjects conn.Send(MSET) error(%v)", err)
		return
	}
	count := 1
	for _, v := range keyIDs {
		count++
		if err = conn.Send("EXPIRE", v, d.listExpire); err != nil {
			log.Error("AddCacheMatchSubjects conn.Send(Expire, %s, %d) error(%v)", v, d.listExpire, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// CacheEpSeasons .
func (d *Dao) CacheEpSeasons(c context.Context, ids []int64) (res map[int64]*model.Season, err error) {
	var (
		key  string
		args = redis.Args{}
		bss  [][]byte
	)
	for _, sid := range ids {
		key = keySID(sid)
		args = args.Add(key)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if bss, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("CacheEpSeasons conn.Do(MGET,%s) error(%v)", key, err)
		}
		return
	}
	res = make(map[int64]*model.Season, len(ids))
	for _, bs := range bss {
		sea := new(model.Season)
		if bs == nil {
			continue
		}
		if err = json.Unmarshal(bs, sea); err != nil {
			log.Error("CacheEpSeasons json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		res[sea.ID] = sea
	}
	return
}

// CacheEpTeams .
func (d *Dao) CacheEpTeams(c context.Context, ids []int64) (res map[int64]*model.Team, err error) {
	var (
		key  string
		args = redis.Args{}
		bss  [][]byte
	)
	for _, tid := range ids {
		key = keyTeamID(tid)
		args = args.Add(key)
	}
	conn := d.redis.Get(c)
	defer conn.Close()
	if bss, err = redis.ByteSlices(conn.Do("MGET", args...)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("CacheEpTeams conn.Do(MGET,%s) error(%v)", key, err)
		}
		return
	}
	res = make(map[int64]*model.Team, len(ids))
	for _, bs := range bss {
		team := new(model.Team)
		if bs == nil {
			continue
		}
		if err = json.Unmarshal(bs, team); err != nil {
			log.Error("CacheEpTeams json.Unmarshal(%s) error(%v)", string(bs), err)
			err = nil
			continue
		}
		res[team.ID] = team
	}
	return
}

// AddCacheEpTeams .
func (d *Dao) AddCacheEpTeams(c context.Context, data map[int64]*model.Team) (err error) {
	if len(data) == 0 {
		return
	}
	var (
		bs      []byte
		keyID   string
		keyIDs  []string
		argsCid = redis.Args{}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, v := range data {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("AddCacheEpTeams.json.Marshal err(%v)", err)
			continue
		}
		keyID = keyTeamID(v.ID)
		keyIDs = append(keyIDs, keyID)
		argsCid = argsCid.Add(keyID).Add(string(bs))
	}
	if err = conn.Send("MSET", argsCid...); err != nil {
		log.Error("AddCacheEpTeams conn.Send(MSET) error(%v)", err)
		return
	}
	count := 1
	for _, v := range keyIDs {
		count++
		if err = conn.Send("EXPIRE", v, d.listExpire); err != nil {
			log.Error("AddCacheEpTeams conn.Send(Expire, %s, %d) error(%v)", v, d.listExpire, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddCacheEpSeasons .
func (d *Dao) AddCacheEpSeasons(c context.Context, data map[int64]*model.Season) (err error) {
	if len(data) == 0 {
		return
	}
	var (
		bs      []byte
		keyID   string
		keyIDs  []string
		argsCid = redis.Args{}
	)
	conn := d.redis.Get(c)
	defer conn.Close()
	for _, v := range data {
		if bs, err = json.Marshal(v); err != nil {
			log.Error("json.Marshal err(%v)", err)
			continue
		}
		keyID = keySID(v.ID)
		keyIDs = append(keyIDs, keyID)
		argsCid = argsCid.Add(keyID).Add(string(bs))
	}
	if err = conn.Send("MSET", argsCid...); err != nil {
		log.Error("AddCacheEpSeasons conn.Send(MSET) error(%v)", err)
		return
	}
	count := 1
	for _, v := range keyIDs {
		count++
		if err = conn.Send("EXPIRE", v, d.listExpire); err != nil {
			log.Error("AddCacheEpSeasons conn.Send(Expire, %s, %d) error(%v)", v, d.listExpire, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// GetActPageCache get act from cache.
func (d *Dao) GetActPageCache(c context.Context, id int64) (act *model.ActivePage, err error) {
	var (
		bs   []byte
		key  = keyMatchAct(id)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			act = nil
		} else {
			log.Error("conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	act = new(model.ActivePage)
	if err = json.Unmarshal(bs, act); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// GetCSingleData get contest single data
func (d *Dao) GetCSingleData(c context.Context, id int64) (data *model.ContestDataPage, err error) {
	var (
		bs   []byte
		key  = keyCSData(id)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			data = nil
		} else {
			log.Error("GetCSingleData conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	data = new(model.ContestDataPage)
	if err = json.Unmarshal(bs, data); err != nil {
		log.Error("GetCSingleData json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// AddCSingleData add act first page value
func (d *Dao) AddCSingleData(c context.Context, id int64, act *model.ContestDataPage) (err error) {
	var (
		bs   []byte
		key  = keyCSData(id)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(act); err != nil {
		log.Error("AddCSingleData json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("AddCSingleData conn.Send(SET,%s,%d) error(%v)", key, id, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddCSingleData conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}

// GetCRecent get contest recent data
func (d *Dao) GetCRecent(c context.Context, param *model.ParamCDRecent) (data []*model.Contest, err error) {
	var (
		bs   []byte
		key  = keyCRecent(param)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			data = nil
		} else {
			log.Error("GetCRecent conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	data = make([]*model.Contest, 0)
	if err = json.Unmarshal(bs, &data); err != nil {
		log.Error("GetCRecent json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// AddCRecent add contest recent data
func (d *Dao) AddCRecent(c context.Context, param *model.ParamCDRecent, data []*model.Contest) (err error) {
	var (
		bs   []byte
		key  = keyCRecent(param)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(data); err != nil {
		log.Error("AddCRecent json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("AddCRecent conn.Send(SET,%s,%v) error(%v)", key, param, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddCRecent conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddCRecent conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddCRecent conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}

// AddActPageCache add act first page value
func (d *Dao) AddActPageCache(c context.Context, aid int64, act *model.ActivePage) (err error) {
	var (
		bs   []byte
		key  = keyMatchAct(aid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(act); err != nil {
		log.Error("json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("conn.Send(SET,%s,%d) error(%v)", key, aid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}

// GetActModuleCache get module from cache.
func (d *Dao) GetActModuleCache(c context.Context, mmid int64) (res []*arcmdl.Arc, err error) {
	var (
		bs   []byte
		key  = keyMatchModule(mmid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
		} else {
			log.Error("GetModuleCache conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(bs, &res); err != nil {
		log.Error("GetModuleCache json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// AddActModuleCache add act first page cache
func (d *Dao) AddActModuleCache(c context.Context, mmid int64, module []*arcmdl.Arc) (err error) {
	var (
		bs   []byte
		key  = keyMatchModule(mmid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(module); err != nil {
		log.Error("AddActModuleCache json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("AddActModuleCache conn.Send(SET,%s,%d) error(%v)", key, mmid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddActModuleCache conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddActModuleCache add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddActModuleCache add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}

// GetActTopCache get act top value cache
func (d *Dao) GetActTopCache(c context.Context, aid, ps int64) (res []*model.Contest, total int, err error) {
	key := keyTop(aid, ps)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("GetActTopCache conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		cont := &model.Contest{}
		if err = json.Unmarshal(bs, cont); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, cont)
	}
	total = from(num)
	return
}

// AddActTopCache add act top cache
func (d *Dao) AddActTopCache(c context.Context, aid, ps int64, tops []*model.Contest, total int) (err error) {
	key := keyTop(aid, ps)
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("AddActTopCache conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	args := redis.Args{}.Add(key)
	for sort, contest := range tops {
		bs, _ := json.Marshal(contest)
		args = args.Add(combine(int64(sort), total)).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("AddActTopCache conn.Send(ZADD, %s, %v) error(%v)", key, args, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddActTopCache conn.Send(Expire, %s, %d) error(%v)", key, d.listExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("AddActTopCache conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddActTopCache conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// GetActPointsCache get act point value
func (d *Dao) GetActPointsCache(c context.Context, aid, mdID, ps int64) (res []*model.Contest, total int, err error) {
	key := keyPoint(aid, mdID, ps)
	conn := d.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
	if err != nil {
		log.Error("GetActTopCache conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	var num int64
	for len(values) > 0 {
		bs := []byte{}
		if values, err = redis.Scan(values, &bs, &num); err != nil {
			log.Error("redis.Scan(%v) error(%v)", values, err)
			return
		}
		cont := &model.Contest{}
		if err = json.Unmarshal(bs, cont); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", bs, err)
			return
		}
		res = append(res, cont)
	}
	total = from(num)
	return
}

// AddActPointsCache add act point data cache
func (d *Dao) AddActPointsCache(c context.Context, aid, mdID, ps int64, points []*model.Contest, total int) (err error) {
	key := keyPoint(aid, mdID, ps)
	conn := d.redis.Get(c)
	defer conn.Close()
	count := 0
	if err = conn.Send("DEL", key); err != nil {
		log.Error("AddActTopCache conn.Send(DEL, %s) error(%v)", key, err)
		return
	}
	count++
	args := redis.Args{}.Add(key)
	for sort, contest := range points {
		bs, _ := json.Marshal(contest)
		args = args.Add(combine(int64(sort), total)).Add(bs)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("AddActTopCache conn.Send(ZADD, %s, %v) error(%v)", key, args, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddActTopCache conn.Send(Expire, %s, %d) error(%v)", key, d.listExpire, err)
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("AddActTopCache conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddActTopCache conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// GetActKnockoutCache add act knockout cache value
func (d *Dao) GetActKnockoutCache(c context.Context, mdID int64) (res [][]*model.TreeList, err error) {
	var (
		bs   []byte
		key  = keyKnock(mdID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
		} else {
			log.Error("GetActKnockoutCache conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(bs, &res); err != nil {
		log.Error("GetActKnockoutCache json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// AddActKnockoutCache add act knockout cache value
func (d *Dao) AddActKnockoutCache(c context.Context, mdID int64, knock [][]*model.TreeList) (err error) {
	var (
		bs   []byte
		key  = keyKnock(mdID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(knock); err != nil {
		log.Error("AddActKnockoutCache json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("AddActKnockoutCache conn.Send(SET,%s,%d) error(%v)", key, mdID, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddActKnockoutCache conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddActKnockoutCache add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddActModuleCache add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}

// AddActKnockCacheTime add act knockout cache value time
func (d *Dao) AddActKnockCacheTime(c context.Context, mdID int64) (err error) {
	var (
		key  = keyKnock(mdID)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddActKnockCacheTime conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddActKnockCacheTime add conn.Flush error(%v)", err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		log.Error("AddActKnockCacheTime add error(%v)", err)
		return
	}
	return
}

// GetMActCache get act cache value
func (d *Dao) GetMActCache(c context.Context, aid int64) (res *model.Active, err error) {
	var (
		bs   []byte
		key  = keyMAct(aid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
		} else {
			log.Error("GetMActCache conn.Do(GET,%s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(bs, &res); err != nil {
		log.Error("GetMActCache json.Unmarshal(%s) error(%v)", string(bs), err)
	}
	return
}

// AddMActCache add act cache value
func (d *Dao) AddMActCache(c context.Context, aid int64, act *model.Active) (err error) {
	var (
		bs   []byte
		key  = keyMAct(aid)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if bs, err = json.Marshal(act); err != nil {
		log.Error("AddMActCache json.Marshal() error(%v)", err)
		return
	}
	if err = conn.Send("SET", key, bs); err != nil {
		log.Error("AddMActCache conn.Send(SET,%s,%d) error(%v)", key, aid, err)
		return
	}
	if err = conn.Send("EXPIRE", key, d.listExpire); err != nil {
		log.Error("AddMActCache conn.Send(EXPIRE,%s) error(%v)", key, err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("AddMActCache add conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("AddMActCache add conn.Receive()%d error(%v)", i+1, err)
			return
		}
	}
	return
}
