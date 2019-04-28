package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/job/live/recommend-job/internal/conf"
	"go-common/app/service/live/recommend/recconst"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

// UserAreaJob 把用户分区偏好导入到redis
type UserAreaJob struct {
	JobConf    *conf.JobConfig
	RedisConf  *redis.Config
	HadoopConf *conf.HadoopConfig
}

// Run ...
func (j *UserAreaJob) Run() {
	log.Info("UserAreaJob Start")
	processFile(j.JobConf, j.HadoopConf, j.RedisConf, writeUserAreaToRedis)
	log.Info("UserAreaJob End")
}

func writeUserAreaToRedis(line string, pool *redis.Pool) (err error) {
	var split = strings.Split(line, ",")
	var uid = split[0]
	uid = strings.Trim(uid, "\"")
	var areaIds = split[1]
	areaIds = strings.Trim(areaIds, "\"")
	var ctx = context.Background()
	var conn = pool.Get(ctx)
	defer conn.Close()
	uidInt, _ := strconv.Atoi(uid)
	var key = fmt.Sprintf(recconst.UserAreaKey, uidInt)
	_, err = conn.Do("SETEX", key, 86400*7, areaIds)
	if err != nil {
		log.Error("writeUserAreaToRedis err +%v, key=%s", err, key)
	}
	return
}
