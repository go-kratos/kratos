package server

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"go-common/app/interface/main/broadcast/conf"
	pb "go-common/app/service/main/broadcast/api/grpc/v1"
	"go-common/library/conf/env"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil"

	"github.com/google/uuid"
	"github.com/zhenjl/cityhash"
)

var (
	_maxInt          = 1<<31 - 1
	_roomOnlineValue atomic.Value
)

const (
	_clientHeartbeat    = time.Second * 90
	_minHeartbeatSecond = 600  // 10m
	_maxHeartbeatSecond = 1200 // 20m
	_roomsShard         = 32
)

// Server .
type Server struct {
	c         *conf.Config
	round     *Round    // accept round store
	buckets   []*Bucket // subkey bucket
	bucketIdx uint32

	serverID  string
	rpcClient pb.ZergClient
	httpCli   *bm.Client
	backoff   *netutil.BackoffConfig
}

// NewServer returns a new Server.
func NewServer(c *conf.Config) *Server {
	var err error
	s := new(Server)
	s.c = c
	s.serverID = env.Hostname
	s.rpcClient, err = pb.NewClient(c.WardenClient)
	if err != nil {
		panic(err)
	}
	s.httpCli = bm.NewClient(c.HTTPClient)
	s.round = NewRound(conf.Conf)
	// init bucket
	s.buckets = make([]*Bucket, c.Bucket.Size)
	s.bucketIdx = uint32(c.Bucket.Size)
	for i := 0; i < c.Bucket.Size; i++ {
		s.buckets[i] = NewBucket(c.Bucket)
	}
	s.backoff = &netutil.BackoffConfig{
		MaxDelay:  5 * time.Second,
		BaseDelay: 1.0 * time.Second,
		Factor:    1.6,
		Jitter:    0.2,
	}
	s.loadOnline()
	go s.onlineproc()
	return s
}

// Buckets return all buckets.
func (s *Server) Buckets() []*Bucket {
	return s.buckets
}

// Bucket get the bucket by subkey.
func (s *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	if s.c.Broadcast.Debug {
		log.Info("%s hit channel bucket index: %d use cityhash", subKey, idx)
	}
	return s.buckets[idx]
}

// NextKey generate a server key.
func (s *Server) NextKey() string {
	u, err := uuid.NewRandom()
	if err == nil {
		return u.String()
	}
	return fmt.Sprintf("%s-%d", s.serverID, time.Now().UnixNano())
}

// RandServerHearbeat rand server heartbeat.
func (s *Server) RandServerHearbeat() time.Duration {
	return time.Duration(_minHeartbeatSecond+rand.Intn(_maxHeartbeatSecond-_minHeartbeatSecond)) * time.Second
}

// Close close the server.
func (s *Server) Close() (err error) {
	return
}

func (s *Server) onlineproc() {
	var retry int
	for {
		if err := s.loadOnline(); err != nil {
			retry++
			time.Sleep(s.backoff.Backoff(retry))
			continue
		}
		retry = 0
		time.Sleep(time.Duration(s.c.Broadcast.OnlineTick))
	}
}

func (s *Server) loadOnline() (err error) {
	roomCountShard := make(map[uint32]map[string]int32)
	for _, bucket := range s.buckets {
		for roomID, count := range bucket.RoomsCount() {
			hash := cityhash.CityHash32([]byte(roomID), uint32(len(roomID))) % _roomsShard
			roomCount, ok := roomCountShard[hash]
			if !ok {
				roomCount = make(map[string]int32)
				roomCountShard[hash] = roomCount
			}
			roomCount[roomID] += count
		}
	}
	allRoomsCount := make(map[string]int32)
	for i := uint32(0); i < _roomsShard; i++ {
		var mergedRoomsCount map[string]int32
		mergedRoomsCount, err = s.RenewOnline(context.Background(), s.serverID, int32(i), roomCountShard[i])
		if err != nil {
			log.Error("s.RenewOnline(%s, %d, %d) error(%v)", s.serverID, i, len(roomCountShard[i]), err)
			return
		}
		for roomID, count := range mergedRoomsCount {
			allRoomsCount[roomID] = count
		}
	}
	for _, bucket := range s.buckets {
		bucket.UpRoomsCount(allRoomsCount)
	}
	_roomOnlineValue.Store(allRoomsCount)
	return
}

func roomOnline(rid string) int32 {
	online, ok := _roomOnlineValue.Load().(map[string]int32)
	if !ok {
		return 0
	}
	return online[rid]
}
