package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"go-common/app/service/main/rank/conf"
	"go-common/app/service/main/rank/dao"
	"go-common/app/service/main/rank/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_bucketSize int64 = 100000 // 10w
)

// Service struct
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	rmap   map[int][]*model.Field
	mux    *sync.RWMutex
	waiter *sync.WaitGroup
	// databus
	statViewSub  *databus.Databus
	archiveSub   *databus.Databus
	archiveTvSub *databus.Databus
	procChan     chan *model.StatViewMsg
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		mux:    new(sync.RWMutex),
		waiter: new(sync.WaitGroup),
		// rank map
		rmap: make(map[int][]*model.Field),
		// databus
		statViewSub:  databus.New(c.Databus.StatView),
		archiveSub:   databus.New(c.Databus.Archive),
		archiveTvSub: databus.New(c.Databus.UgcTvBinlog),

		procChan: make(chan *model.StatViewMsg, 1024),
	}
	if err := os.MkdirAll(s.c.Rank.FilePath, 0644); err != nil {
		panic(err)
	}
	if s.c.Rank.SwitchIncr {
		s.incr()
	}
	s.flush()
	s.waiter.Add(1)
	go s.dump()
	fmt.Println("map len:", len(s.rmap))
	return
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
	if err := s.statViewSub.Close(); err != nil {
		log.Error("s.statViewSub.Close() error(%v)", err)
	}
	if err := s.archiveSub.Close(); err != nil {
		log.Error("s.statViewSub.Close() error(%v)", err)
	}
	if err := s.archiveTvSub.Close(); err != nil {
		log.Error("s.archiveTvSub.Close() error(%v)", err)
	}
}

func (s *Service) mod(oid int64) int {
	return int(oid % _bucketSize)
}

func (s *Service) bucket(oid int64) int {
	return int(oid / _bucketSize)
}

func (s *Service) incr() {
	for i := 0; i < runtime.NumCPU(); i++ {
		s.waiter.Add(1)
		go s.batchProc(s.procChan)
	}
	s.waiter.Add(1)
	go s.consumeStatView()
	s.waiter.Add(1)
	go s.consumeArchive()
	s.waiter.Add(1)
	go s.consumeArchiveTv()
}

func (s *Service) flush() error {
	c := context.Background()
	files, err := ioutil.ReadDir(s.c.Rank.FilePath)
	if err != nil {
		log.Error("ioutil.ReadDir() error(%v)", err)
	}
	if len(files) < 2 { // one snapshot.pb and a timestamp.txt at least
		if !s.c.Rank.SwitchAll {
			return nil
		}
		return s.all(c, 0, 0)
	}
	maxID, err := s.dao.MaxOid(context.Background())
	if err != nil {
		log.Error("s.dao.MaxOid() error(%v)", err)
		return err
	}
	for i := 0; i <= int(maxID/_bucketSize); i++ {
		if s.load(i) != nil {
			s.all(c, int64(i)*_bucketSize, int64(i+1)*_bucketSize)
			continue
		}
	}
	fi, err := os.Open(s.c.Rank.FilePath + "timestamp.txt")
	if err != nil {
		log.Error(" os.Open(%s) error(%v)", s.c.Rank.FilePath+"timestamp.txt", err)
		return err
	}
	defer fi.Close()
	data, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		return err
	}
	begin, _ := strconv.ParseInt(string(data[:]), 10, 64)
	if begin <= 1 {
		begin = time.Now().Unix() - 3600 // 1h data
	}
	return s.patch(c, time.Unix(begin, 0), time.Now())
}

func (s *Service) load(index int) error {
	fi, err := os.Open(s.c.Rank.FilePath + fmt.Sprintf(s.c.Rank.FileName, index))
	if err != nil {
		log.Error(" os.Open() error(%v)", err)
		return err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		return err
	}
	fs := new(model.Fields)
	if err := fs.Unmarshal(fd); err != nil {
		log.Error("fs.Unmarshal() error(%v)", err)
		return err
	}
	s.rmap[index] = fs.Fields
	return nil
}

func (s *Service) dump() error {
	for {
		time.Sleep(time.Duration(s.c.Rank.Ticker))
		now := time.Now().Unix()
		for k, sli := range s.rmap {
			csli := make([]*model.Field, 0)
			for _, v := range sli {
				cv := *v
				csli = append(csli, &cv)
			}
			fs := new(model.Fields)
			fs.Fields = csli
			buf, err := fs.Marshal()
			if err != nil {
				log.Error("fs.Marshal() error(%v)", err)
			}
			path := fmt.Sprintf("/tmp/"+s.c.Rank.FileName, k)
			if err = ioutil.WriteFile(path, buf, 0644); err != nil {
				log.Error("ioutil.WriteFile(%s) error(%v)", path, err)
			}
			newPath := s.c.Rank.FilePath + fmt.Sprintf(s.c.Rank.FileName, k)
			s.fileCopy(path, newPath)
		}
		if err := ioutil.WriteFile(s.c.Rank.FilePath+"timestamp.txt", []byte(fmt.Sprintf("%d", now)), 0644); err != nil {
			log.Error("ioutil.WriteFile(%d) error(%v)", now, err)
		}
		log.Info("rmap dump successful now(%d)", now)
	}
}

func (s *Service) fileCopy(path, newPath string) {
	originalFile, err := os.Open(path)
	if err != nil {
		log.Error("os.Open(%s) error(%v)", path, err)
	}
	defer originalFile.Close()
	newFile, err := os.Create(newPath)
	if err != nil {
		log.Error("os.Create(%s) error(%v)", path, err)
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		log.Error("os.Copy() error(%v)", err)
	}
	err = newFile.Sync()
	if err != nil {
		log.Error("newFile.Sync() error(%v)", err)
	}
}

func (s *Service) field(oid int64) (f *model.Field) {
	bk := s.bucket(oid)
	md := s.mod(oid)
	s.mux.RLock()
	r, ok := s.rmap[bk]
	s.mux.RUnlock()
	if !ok {
		s.mux.Lock()
		r, ok = s.rmap[bk]
		if !ok {
			for i := 0; i < int(_bucketSize); i++ {
				r = append(r, new(model.Field))
			}
			s.rmap[bk] = r
		}
		s.mux.Unlock()
	}
	if f = r[md]; f == nil {
		f = new(model.Field)
		r[md] = f
	}
	return
}

func (s *Service) setField(oid int64, f *model.Field) {
	bk := s.bucket(oid)
	md := s.mod(oid)
	s.mux.RLock()
	r, ok := s.rmap[bk]
	s.mux.RUnlock()
	if !ok {
		s.mux.Lock()
		r, ok = s.rmap[bk]
		if !ok {
			for i := 0; i < int(_bucketSize); i++ {
				r = append(r, new(model.Field))
			}
			s.rmap[bk] = r
		}
		s.mux.Unlock()
	}
	r[md] = f
}
