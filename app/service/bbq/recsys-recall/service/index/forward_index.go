package index

import (
	"os"
	"sync"
	"syscall"
	"time"

	"go-common/app/job/bbq/recall/proto"
	"go-common/app/service/bbq/recsys-recall/conf"
	"go-common/library/log"
)

var (
	// Index 正排索引
	Index  *ForwardIndex
	loader Loader
)

// ForwardIndex 正排索引
type ForwardIndex struct {
	data *map[uint64]*proto.ForwardIndex
	lock sync.RWMutex
	ch   chan os.Signal
}

// Init .
func Init(c *conf.Config) {
	loader = &LocalLoader{
		path: c.ForwardIndex.LocalPath,
	}
	if c.ForwardIndex.LocalPath == "" {
		loader = &RemoteLoader{
			path:    c.ForwardIndex.RemotePath,
			md5Path: c.ForwardIndex.MD5Path,
			md5:     "",
		}
	}
	d, err := loader.Load()
	if err != nil {
		log.Error("Load ForwardIndex from local error %+v", err)
		panic(err)
	}

	Index = &ForwardIndex{
		data: d,
		ch:   make(chan os.Signal, 1),
	}

	go func(c *conf.Config) {
		if c.ForwardIndex.RemotePath == "" || c.ForwardIndex.MD5Path == "" {
			return
		}
		duration := time.Duration(c.ForwardIndex.ReloadDucation)
		loader := &RemoteLoader{
			path:    c.ForwardIndex.RemotePath,
			md5Path: c.ForwardIndex.MD5Path,
			md5:     "",
		}
		for {
			select {
			case s := <-Index.ch:
				log.Info("ForwardIndex get a signal %s", s.String())
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					return
				case syscall.SIGHUP:
				default:
					return
				}
			case <-time.After(duration):
				data, err := loader.Load()
				if err != nil || data == nil {
					log.Error("ForwardIndex load data %+v", err)
					continue
				}

				Index.lock.Lock()
				Index.data = data
				Index.lock.Unlock()
				log.Info("reload the positive index")
			}
		}
	}(c)
}

// Signal .
func (pi *ForwardIndex) Signal(s os.Signal) {
	pi.ch <- s
}

// Get .
func (pi *ForwardIndex) Get(svid uint64) *proto.ForwardIndex {
	pi.lock.RLocker().Lock()
	defer pi.lock.RUnlock()
	if res, ok := (*pi.data)[svid]; ok {
		return res
	}
	return nil
}
