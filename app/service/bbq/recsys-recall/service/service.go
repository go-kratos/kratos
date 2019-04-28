package service

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/json-iterator/go"

	"go-common/app/job/bbq/recall/proto"
	"go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	"go-common/app/service/bbq/recsys-recall/conf"
	"go-common/app/service/bbq/recsys-recall/dao"
	"go-common/app/service/bbq/recsys-recall/service/index"
	"go-common/library/log"

	"github.com/Dai0522/workpool"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	wp  *workpool.Pool
}

// New init
func New(c *conf.Config) (s *Service) {
	workpoolConf := &workpool.PoolConfig{
		MaxWorkers:     c.WorkPool.MaxWorkers,
		MaxIdleWorkers: c.WorkPool.MaxIdleWorkers,
		MinIdleWorkers: c.WorkPool.MinIdleWorkers,
		KeepAlive:      time.Duration(c.WorkPool.KeepAlive),
	}
	wp, err := workpool.NewWorkerPool(c.WorkPool.Capacity, workpoolConf)
	if err != nil {
		panic(err)
	}

	wp.Start()
	s = &Service{
		c:   c,
		dao: dao.New(c),
		wp:  wp,
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
}

// VideoIndex 获取视频正排信息
func (s *Service) VideoIndex(ctx context.Context, in *v1.VideoIndexRequest) (*v1.VideoIndexResponse, error) {
	var idxs []*proto.ForwardIndex
	for _, v := range in.SVIDs {
		fi := index.Index.Get(uint64(v))
		if fi == nil {
			log.Errorv(ctx, log.KV("forward_index", nil))
			continue
		}
		idxs = append(idxs, fi)
	}

	return &v1.VideoIndexResponse{
		List: idxs,
	}, nil
}

// NewIncomeVideo 更新新发视频标签
func (s *Service) NewIncomeVideo(ctx context.Context, in *v1.NewIncomeVideoRequest) (res *empty.Empty, err error) {
	res = new(empty.Empty)

	svids := make([]uint64, len(in.SVIDs))
	for i := range in.SVIDs {
		svids[i] = uint64(in.SVIDs[i])
	}

	ii := &index.InvertedIndex{
		Data: svids,
	}
	s.dao.SetInvertedIndex(ctx, in.Key, ii.Serialize())

	return
}

// VideosByIndex 获取单个倒排下的视频列表
func (s *Service) VideosByIndex(ctx context.Context, in *v1.VideosByIndexRequest) (res *v1.VideosByIndexResponse, err error) {
	raw, err := s.dao.GetInvertedIndex(ctx, in.Key, true)
	if err != nil {
		return
	}

	var recallList []uint64
	if binary.BigEndian.Uint64(raw[:8]) == 0xffffffffdeadbeef {
		ii := new(index.InvertedIndex)
		if err = ii.Load(raw); err != nil {
			log.Errorv(ctx, log.KV("Tag", in.Key), log.KV("inverted index load", err))
			return
		}
		recallList = ii.Data
	} else {
		if err = jsoniter.Unmarshal(raw, &recallList); err != nil {
			log.Errorv(ctx, log.KV("Tag", in.Key), log.KV("jsoninter", err))
			return
		}
	}

	svidList := make([]int64, len(recallList))
	for i := range recallList {
		svidList[i] = int64(recallList[i])
	}
	res = &v1.VideosByIndexResponse{
		Key:   in.Key,
		SVIDs: svidList,
	}

	return
}
