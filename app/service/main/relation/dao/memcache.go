package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/service/main/relation/model"
	gmc "go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_prefixFollowing      = "pb_a_"      // key of following
	_prefixFollower       = "pb_f_"      // key of follower
	_prefixTagCount       = "rs_tmtc_%d" // key of relation tag by mid & tag's count
	_prefixTags           = "tags_"      // user tag info.
	_prefixGlobalHotRec   = "gh_rec"     // global hot rec
	_prefixStat           = "c_"         // key of stat
	_prefixFollowerNotify = "f_notify_"
	_emptyExpire          = 20 * 24 * 3600
	_recExpire            = 5 * 24 * 3600
)

func statKey(mid int64) string {
	return _prefixStat + strconv.FormatInt(mid, 10)
}

func tagsKey(mid int64) string {
	return _prefixTags + strconv.FormatInt(mid, 10)
}
func followingKey(mid int64) string {
	return _prefixFollowing + strconv.FormatInt(mid, 10)
}

func followerKey(mid int64) string {
	return _prefixFollower + strconv.FormatInt(mid, 10)
}

func tagCountKey(mid int64) string {
	return fmt.Sprintf(_prefixTagCount, mid)
}

func globalHotKey() string {
	return _prefixGlobalHotRec
}

func followerNotifySetting(mid int64) string {
	return _prefixFollowerNotify + strconv.FormatInt(mid, 10)
}

// pingMC ping memcache.
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Set(&gmc.Item{Key: "ping", Value: []byte{1}, Expiration: d.mcExpire}); err != nil {
		log.Error("conn.Store(set, ping, 1) error(%v)", err)
	}
	conn.Close()
	return
}

// SetFollowingCache set following cache.
func (d *Dao) SetFollowingCache(c context.Context, mid int64, followings []*model.Following) (err error) {
	return d.setFollowingCache(c, followingKey(mid), followings)
}

// FollowingCache get following cache.
func (d *Dao) FollowingCache(c context.Context, mid int64) (followings []*model.Following, err error) {
	return d.followingCache(c, followingKey(mid))
}

// DelFollowingCache delete following cache.
func (d *Dao) DelFollowingCache(c context.Context, mid int64) (err error) {
	return d.delFollowingCache(c, followingKey(mid))
}

// SetFollowerCache set follower cache.
func (d *Dao) SetFollowerCache(c context.Context, mid int64, followers []*model.Following) (err error) {
	return d.setFollowingCache(c, followerKey(mid), followers)
}

// FollowerCache get follower cache.
func (d *Dao) FollowerCache(c context.Context, mid int64) (followers []*model.Following, err error) {
	return d.followingCache(c, followerKey(mid))
}

// DelFollowerCache delete follower cache.
func (d *Dao) DelFollowerCache(c context.Context, mid int64) (err error) {
	return d.delFollowingCache(c, followerKey(mid))
}

// setFollowingCache set following cache.
func (d *Dao) setFollowingCache(c context.Context, key string, followings []*model.Following) (err error) {
	expire := d.followerExpire
	if len(followings) == 0 {
		expire = _emptyExpire
	}
	item := &gmc.Item{Key: key, Object: &model.FollowingList{FollowingList: followings}, Expiration: expire, Flags: gmc.FlagProtobuf}
	conn := d.mc.Get(c)
	if err = conn.Set(item); err != nil {
		log.Error("setFollowingCache err(%v)", err)
	}
	conn.Close()
	return
}

// followingCache get following cache.
func (d *Dao) followingCache(c context.Context, key string) (followings []*model.Following, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("d.followingCache err(%v)", err)
		return
	}
	followingList := &model.FollowingList{}
	if err = conn.Scan(item, followingList); err != nil {
		log.Error("d.followinfCache err(%v)", err)
	}
	followings = followingList.FollowingList
	if followings == nil {
		followings = []*model.Following{}
	}
	return
}

// delFollowingCache delete following cache.
func (d *Dao) delFollowingCache(c context.Context, key string) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", key, err)
		}
	}
	conn.Close()
	return
}

// TagCountCache tag count cache
func (d *Dao) TagCountCache(c context.Context, mid int64) (tagCount []*model.TagCount, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	res, err := conn.Get(tagCountKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Error("mc.Get error(%v)", err)
		return
	}
	tc := &model.TagCountList{}
	if err = conn.Scan(res, tc); err != nil {
		log.Error("conn.Scan error(%v)", err)
	}
	tagCount = tc.TagCountList
	return
}

// SetTagCountCache set tag count cache
func (d *Dao) SetTagCountCache(c context.Context, mid int64, tagCount []*model.TagCount) (err error) {
	item := &gmc.Item{Key: tagCountKey(mid), Object: &model.TagCountList{TagCountList: tagCount}, Expiration: d.mcExpire, Flags: gmc.FlagProtobuf}
	conn := d.mc.Get(c)
	if err = conn.Set(item); err != nil {
		log.Error("setTagMidFidCache(%s) error(%v)", tagCountKey(mid), err)
	}
	conn.Close()
	return
}

// DelTagCountCache del tag count cache
func (d *Dao) DelTagCountCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(tagCountKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", tagCountKey(mid), err)
		}
	}
	conn.Close()
	return
}

// SetTagsCache set user tags cache.
func (d *Dao) SetTagsCache(c context.Context, mid int64, tags *model.Tags) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: tagsKey(mid), Object: tags, Expiration: d.mcExpire, Flags: gmc.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		log.Error("SetTagsCache err %v", err)
	}
	return
}

// TagsCache get user tags.
func (d *Dao) TagsCache(c context.Context, mid int64) (tags map[int64]*model.Tag, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	res, err := conn.Get(tagsKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		return
	}
	tag := new(model.Tags)
	if err = conn.Scan(res, tag); err != nil {
		return
	}
	tags = tag.Tags
	return
}

// DelTagsCache del user tags cache.
func (d *Dao) DelTagsCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	if err = conn.Delete(tagsKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		} else {
			log.Error("conn.Delete(%s) error(%v)", tagCountKey(mid), err)
		}
	}
	conn.Close()
	return
}

// SetGlobalHotRecCache set global hot recommend cache.
func (d *Dao) SetGlobalHotRecCache(c context.Context, fids []int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{Key: globalHotKey(), Object: &model.GlobalRec{Fids: fids}, Expiration: _recExpire, Flags: gmc.FlagProtobuf}
	if err = conn.Set(item); err != nil {
		log.Error("SetGlobalHotRecCache err %v", err)
	}
	return
}

// GlobalHotRecCache get global hot recommend.
func (d *Dao) GlobalHotRecCache(c context.Context) (fs []int64, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	res, err := conn.Get(globalHotKey())
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
		}
		return
	}
	gh := new(model.GlobalRec)
	if err = conn.Scan(res, gh); err != nil {
		return
	}
	fs = gh.Fids
	return
}

// SetStatCache set stat cache.
func (d *Dao) SetStatCache(c context.Context, mid int64, st *model.Stat) error {
	conn := d.mc.Get(c)
	defer conn.Close()
	item := &gmc.Item{
		Key:        statKey(mid),
		Object:     st,
		Expiration: d.mcExpire,
		Flags:      gmc.FlagProtobuf,
	}
	if err := conn.Set(item); err != nil {
		log.Error("Failed to set stat cache: mid: %d stat: %+v: %+v", mid, st, err)
		return err
	}
	return nil
}

// statCache get stat cache.
func (d *Dao) statCache(c context.Context, mid int64) (*model.Stat, error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	item, err := conn.Get(statKey(mid))
	if err != nil {
		if err == gmc.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	stat := &model.Stat{}
	if err := conn.Scan(item, stat); err != nil {
		log.Error("Failed to get stat cache: mid: %d: %+v", mid, err)
		return nil, err
	}
	return stat, nil
}

// statsCache get multi stat cache.
func (d *Dao) statsCache(c context.Context, mids []int64) (map[int64]*model.Stat, []int64, error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	keys := make([]string, 0, len(mids))
	for _, mid := range mids {
		keys = append(keys, statKey(mid))
	}
	items, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("Failed to get multi stat: keys: %+v: %+v", keys, err)
		return nil, nil, err
	}

	stats := make(map[int64]*model.Stat, len(mids))
	for _, item := range items {
		stat := &model.Stat{}
		if err := conn.Scan(item, stat); err != nil {
			log.Error("Failed to scan item: key: %s item: %+v: %+v", item.Key, item, err)
			continue
		}
		stats[stat.Mid] = stat
	}

	missed := make([]int64, 0, len(mids))
	for _, mid := range mids {
		if _, ok := stats[mid]; !ok {
			missed = append(missed, mid)
		}
	}

	return stats, missed, nil
}

// DelStatCache delete stat cache.
func (d *Dao) DelStatCache(c context.Context, mid int64) error {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err := conn.Delete(statKey(mid)); err != nil {
		if err == gmc.ErrNotFound {
			return nil
		}
		log.Error("Failed to delete stat cache: mid: %d: %+v", mid, err)
		return err
	}
	return nil
}

// GetFollowerNotifyCache get data from mc
func (d *Dao) GetFollowerNotifyCache(c context.Context, mid int64) (res *model.FollowerNotifySetting, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := followerNotifySetting(mid)
	reply, err := conn.Get(key)
	if err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Errorv(c, log.KV("GetFollowerNotifyCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	res = &model.FollowerNotifySetting{}
	err = conn.Scan(reply, res)
	if err != nil {
		log.Errorv(c, log.KV("GetFollowerNotifyCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// SetFollowerNotifyCache Set data to mc
func (d *Dao) SetFollowerNotifyCache(c context.Context, mid int64, val *model.FollowerNotifySetting) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := followerNotifySetting(mid)
	item := &gmc.Item{
		Key:        key,
		Object:     val,
		Expiration: 86400,
		Flags:      gmc.FlagJSON,
	}
	if err = conn.Set(item); err != nil {
		log.Errorv(c, log.KV("SetFollowerNotifyCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// DelFollowerNotifyCache Del data from mc
func (d *Dao) DelFollowerNotifyCache(c context.Context, mid int64) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := followerNotifySetting(mid)
	if err = conn.Delete(key); err != nil {
		if err == gmc.ErrNotFound {
			err = nil
			return
		}
		log.Errorv(c, log.KV("DelFollowerNotifyCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
