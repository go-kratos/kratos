package service

import (
	"context"
	"crypto/md5"
	"fmt"
	topic "go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/video/conf"
	"go-common/app/service/bbq/video/dao"
	"go-common/library/conf/env"
	"go-common/library/log"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

var (
	hostHashIndex  int64
	autoIncreaseID int64
)

// Service struct
type Service struct {
	c           *conf.Config
	BPSCode     map[string]int
	dao         *dao.Dao
	topicClient topic.TopicClient
}

func init() {
	hostName, err := os.Hostname()
	if err != nil {
		errors.Wrap(err, "获取hostname失败")
		panic(err)
	}
	data := []byte(hostName)
	hash := fmt.Sprintf("%x", md5.Sum(data))
	log.Infov(context.Background(), log.KV("log", "md5="+hash))
	truncateHash := hash[0:8]
	index, err := strconv.ParseInt(truncateHash, 16, 0)
	if err != nil {
		errors.Wrap(err, "解析MD5->int64失败")
		panic(err)
	}
	hostHashIndex = index % 1000
	log.Infov(context.Background(), log.KV("log", fmt.Sprintf("hostname hash index=%03d", hostHashIndex)))
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	s.dao.HostnameRegister(hostHashIndex)
	if env.DeployEnv == env.DeployEnvProd {
		go s.archiveSub(context.Background())
	}
	var err error
	if s.topicClient, err = topic.NewClient(nil); err != nil {
		log.Errorw(context.Background(), "log", "get topic client fail")
		panic(err)
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
	s.archiveSubClose()
}
