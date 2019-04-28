package service

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/openplatform/abtest/model"
	"go-common/app/service/openplatform/abtest/model/jump"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

//resetGroupCache 使用数据库信息重置分组缓存
func (s *Service) resetGroupCache(c context.Context, group int, ver int64) (err error) {
	var src []*model.AB
	if src, err = s.d.ActByGroup(c, group); err != nil {
		return
	}
	err = s.setGroupCache(c, group, src, ver)
	return
}

//groupCache 获取分组缓存
func (s *Service) groupCache(c context.Context, group int) (cache map[int]*model.AB, err error) {
	var (
		ok  bool
		ver int64
	)
	for cache, ok = s.readGroupCache(c, group); !ok; cache, ok = s.readGroupCache(c, group) {
		if ver, err = s.d.RedisVersionID(c, group); err != nil {
			if err == redis.ErrNil {
				if err = s.d.SetnxRedisVersionID(c, group, time.Now().Unix()); err != nil {
					log.Error("[service.abtest|groupCache] s.d.SetnxRedisVersionID err: %v", err)
					return
				}
			} else {
				log.Error("[service.abtest|groupCache] s.d.RedisVersionID err: %v", err)
				return
			}
		}
		s.resetGroupCache(c, group, ver)
	}
	return
}

func syncStart(s *Service) {
	ctx := context.TODO()
	for {
		s.SyncVersionID(ctx)
		time.Sleep(time.Second * 30)
	}
}

//SyncVersionID 同步版本
func (s *Service) SyncVersionID(c context.Context) {
	var (
		err          error
		redisVer     int64
		cacheVerList map[int]int64
	)
	if cacheVerList, err = s.VersionIDListCache(c); err != nil {
		log.Error("[service.abtest|SyncVersionID] VersionIDListCache err: %v", err)
		return
	}
	for group, cacheVer := range cacheVerList {
		if redisVer, err = s.d.RedisVersionID(c, group); err != nil {
			if err == redis.ErrNil {
				if err = s.d.SetnxRedisVersionID(c, group, time.Now().Unix()); err != nil {
					log.Error("[service.abtest|groupCache] s.d.SetnxRedisVersionID err: %v", err)
					return
				}
			} else {
				log.Error("[service.abtest|groupCache] s.d.RedisVersionID err: %v", err)
				return
			}
		}
		if redisVer > cacheVer {
			if err = s.resetGroupCache(c, group, redisVer); err != nil {
				log.Error("[service.abtest|SyncVersionID] resetGroupCache(group:%d) err: %v", group, err)
				continue
			}
		} else if redisVer < cacheVer {
			s.d.UpdateRedisVersionID(c, group, cacheVer)
		}
	}
}

// versionByKey 根据key选择现行AB测试配置
func (s *Service) versionByKey(c context.Context, group int, key string) (abMap map[int]int, err error) {
	var cache map[int]*model.AB
	abMap = make(map[int]int)
	if cache, err = s.groupCache(c, group); err != nil {
		log.Error("[service.abtest|versionByKey] s.groupCache err: %v", err)
		return
	}
	for _, ele := range cache {
		if ele.Status != 1 || ele.Group != group {
			continue
		}
		if key != "" {
			//key不为空时散列后计算得到result
			score := int(jump.Hash(jump.Md5(strings.Join([]string{key, string(ele.Seed)}, "")), ele.Stra.Precision))
			if abMap[ele.ID], err = ele.Stra.Version(score); err != nil {
				log.Error("[service.abtest|versionByKey] stra.Version err: %v", err)
				return
			}
		} else {
			abMap[ele.ID] = ele.Result
		}
	}
	return
}

//VersionID returns current version ID
func (s *Service) VersionID(c context.Context, group int) (ver int64, err error) {
	return s.versionID(c, group), nil
}

//Version 获取用户AB测试配置
func (s *Service) Version(c context.Context, group int, key string, ver *model.Version, appkey string) (res *model.Version, err error) {
	res = ver
	nVer := s.versionID(c, group)
	if nVer > ver.VersionID {
		res.VersionID = nVer
		res.Data, err = s.versionByKey(c, group, key)
	}
	if s.c.Stat == 1 {
		go s.stat(c, res, appkey)
		// log.Info("key: %s, result: %s", key, res.Data[s.c.Stra])
	}
	return
}

//Ab 获取单个实验
func (s *Service) Ab(c context.Context, id int, group int) (res *model.AB, err error) {
	if group == 0 {
		return s.d.Ab(c, id)
	}
	return s.d.AbByIDAndGroup(c, id, group)
}

//ListAb 获取实验列表
func (s *Service) ListAb(c context.Context, page, pageSize int, mstatus string, group int) (res []*model.AB, total int, err error) {
	offset := (page - 1) * pageSize
	if group == 0 {
		return s.d.ListAb(c, offset, pageSize, mstatus)
	}
	return s.d.ListAbByGroup(c, offset, pageSize, mstatus, group)
}

//AddAb 添加实验
func (s *Service) AddAb(c context.Context, ab *model.AB) (res map[string]interface{}, err error) {
	var (
		newID int64
		seed  int
		stra  []byte
	)
	seed = rand.Intn(1000000000)
	if stra, err = json.Marshal(ab.Stra); err != nil {
		return
	}
	if newID, err = s.d.AddAb(c, ab.Name, ab.Desc, string(stra), seed, ab.Result, ab.Group, ab.Author); err != nil {
		return
	}
	res = map[string]interface{}{
		"newid": newID,
	}
	return
}

//UpdateAb 更新实验
func (s *Service) UpdateAb(c context.Context, id int, ab *model.AB) (res bool, err error) {
	var stra []byte
	if stra, err = json.Marshal(ab.Stra); err != nil {
		return
	}
	if _, err = s.d.UpAb(c, id, ab.Name, ab.Desc, string(stra), ab.Result, ab.Modifier, ab.Version+1, ab.Status, ab.Group); err != nil {
		return
	}
	if err = s.resetGroupCache(c, ab.Group, time.Now().Unix()); err != nil {
		return
	}
	res = true
	go s.statU(c)
	return
}

//UpdateStatus 更新实验状态
func (s *Service) UpdateStatus(c context.Context, id, status int, username string, group int) (res bool, err error) {
	if _, err = s.d.UpStatus(c, id, status, username, group); err != nil {
		return
	}
	if err = s.resetGroupCache(c, group, time.Now().Unix()); err != nil {
		return
	}
	res = true
	return
}

//DeleteAb 删除实验
func (s *Service) DeleteAb(c context.Context, id int) (res bool, err error) {
	var row int64
	if row, err = s.d.DelAb(c, id); err != nil {
		return
	}
	if row != 0 {
		res = true
	}
	return
}

func (s *Service) stat(c context.Context, data *model.Version, appkey string) {
	for k, v := range data.Data {
		key := "STAT:" + appkey + ":" + strconv.Itoa(k) + ":" + strconv.Itoa(v)
		if _, ok := s.keyList.Load(key); !ok {
			s.keyList.Store(key, model.Empty{})
		}
		s.d.Incr(c, key)
	}
}

func (s *Service) statU(c context.Context) {
	s.keyList.Range(func(key interface{}, value interface{}) bool {
		s.d.Move(c, key.(string))
		return true
	})
}

//Total total of abtest result
func (s *Service) Total(c context.Context) (res map[string]map[string]int, err error) {
	res0 := make(map[string]int)
	res1 := make(map[string]int)
	res = make(map[string]map[string]int)
	s.keyList.Range(func(key interface{}, value interface{}) bool {
		if v, err := s.d.GetFromRedis(c, key.(string)); err == nil {
			res0[key.(string)] = v
		}
		if v, err := s.d.GetFromRedis(c, "O:"+key.(string)); err == nil {
			res1[key.(string)] = v
		}
		res["current"] = res0
		res["last"] = res1
		return true
	})
	return
}
