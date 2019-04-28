package file

import (
	"os"
	"bufio"
	"time"
	"io/ioutil"
	"math/rand"
	"path"
	"strconv"
	"sync"
	"fmt"
	"errors"
	"encoding/binary"
	"encoding/json"
	"sort"
	"bytes"
	"io"
	"strings"
	"go-common/app/service/ops/log-agent/event"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"

	"github.com/fsnotify/fsnotify"
)

const (
	_formatUpdated      = "2006-01-02 15:04:05"
	_logMagicSize       = 2
	_logHeadSize        = 6
	_logLenSize         = 4
	_logIdSize          = 6
	_logLancerHeaderLen = 19
)

var (
	errLogNotFound  = errors.New("log not found")
	errMagicInvaild = errors.New("log magic invalid")
	logMagic        = []byte{0xAC, 0xBE}
	_logType        = []byte{0, 1}
	_logLength      = []byte{0, 0, 0, 0}
	local, _        = time.LoadLocation("Local")
)

// Index index.
type Index struct {
	Name    string `json:"name"`
	Offset  int64  `json:"offset"`
	Updated string `json:"updated"`
}

type FileCache struct {
	c           *Config
	next        chan string
	storageFull bool
	writeChan   chan *event.ProcessorEvent
	readChan    chan *event.ProcessorEvent
	eLock       sync.RWMutex
	logs        map[string]os.FileInfo
	wh          *fsnotify.Watcher
}

func NewFileCache(c *Config) (f *FileCache, err error) {
	if err = c.ConfigValidate(); err != nil {
		return nil, err
	}
	f = new(FileCache)
	f.c = c
	f.storageFull = false
	f.next = make(chan string, 1)
	f.writeChan = make(chan *event.ProcessorEvent)
	f.readChan = make(chan *event.ProcessorEvent)
	f.logs = make(map[string]os.FileInfo)

	if _, err := os.Stat(f.c.Storage); os.IsNotExist(err) {
		if err = os.MkdirAll(f.c.Storage, 0755); err != nil {
			return nil, err
		}
	}

	if err = f.nextFile(); err != nil {
		return nil, err
	}
	if err = f.watch(); err != nil {
		return
	}
	if err = f.loadFiles(); err != nil {
		return
	}

	go f.watchproc()
	go f.writeProcess()
	go f.readProcess()
	return f, nil
}

func (f *FileCache) WriteToCache(e *event.ProcessorEvent) {
	f.writeChan <- e
}

func (f *FileCache) ReadFromCache() (e *event.ProcessorEvent) {
	e = <-f.readChan
	return
}

// loadFiles loadFiles
func (f *FileCache) loadFiles() (err error) {
	var (
		fi  os.FileInfo
		fis []os.FileInfo
	)
	if fis, err = ioutil.ReadDir(f.c.Storage); err != nil {
		log.Error("ioutil.ReadDir(%s) error(%v)", f.c.Storage, err)
		return
	}
	for _, fi = range fis {
		name := path.Join(f.c.Storage, fi.Name())
		if !fi.IsDir() && strings.HasSuffix(name, f.c.Suffix) {
			f.eLock.Lock()
			f.logs[name] = fi
			f.eLock.Unlock()
			log.Info("loadFile: %s, size: %d", name, fi.Size())
		}
	}
	return
}

func (f *FileCache) writeProcess() {
	var (
		err       error
		n, total  int
		lengthbuf = make([]byte, 4)
		cur       *os.File
		wr        = bufio.NewWriterSize(nil, f.c.WriteBuffer)
		tk        = time.Tick(time.Duration(f.c.CacheFlushInterval))
		timestamp = []byte(fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	)
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case next := <-f.next:
			if cur != nil && wr != nil {
				wr.Flush()
				cur.Close()
			}
			f, err := os.OpenFile(next, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Error("os.OpenFile(%s) error(%v)", next, err)
				continue
			}
			cur = f
			wr.Reset(f)
			total = 0
		case <-tk:
			if wr != nil && cur != nil {
				wr.Flush()
			}
			f.checkStorageSize()
		case e := <-f.writeChan:
			if f.storageFull {
				flowmonitor.Fm.AddEvent(e, "log-agent.output.lancer", "ERROR", "file cache storgefull")
				event.PutEvent(e)
				continue
			}

			if total > f.c.FileBytes && len(f.next) == 0 {
				if err := f.nextFile(); err != nil {
					log.Error("c.nextFile() error(%v)", err)
				}
			}

			binary.BigEndian.PutUint32(lengthbuf, uint32(e.Length+_logLancerHeaderLen))
			// write logMagic
			if n, err = wr.Write(logMagic); err != nil {
				goto HERE
			}
			total += n
			// write length
			if n, err = wr.Write(lengthbuf); err != nil {
				goto HERE
			}
			total += n
			// write log
			if n, err = wr.Write([]byte(e.LogId)); err != nil {
				goto HERE
			}

			if n, err = wr.Write(timestamp); err != nil {
				goto HERE
			}

			if n, err = wr.Write(e.Bytes()); err != nil {
				goto HERE
			}

			total += n
			flowmonitor.Fm.AddEvent(e, "log-agent.output.lancer", "OK", "write file cache ok")
			event.PutEvent(e)
			continue

		HERE: // write file cache  error
			flowmonitor.Fm.AddEvent(e, "log-agent.output.lancer", "ERROR", "write file cache failed")
			event.PutEvent(e)
			log.Error("wr.Write() error(%v)", err)
			if cur != nil && wr != nil {
				wr.Flush()
				cur.Close()
			}
			name := f.nextFileName()
			f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Error("os.OpenFile(%s) error(%v)", name, err)
				continue
			}
			cur = f
			wr.Reset(f)
			total = 0
			continue

		}
	}
}

// index index
func (f *FileCache) index() (idx *Index, err error) {
	f.eLock.RLock()
	length := len(f.logs)
	f.eLock.RUnlock()
	if length == 0 {
		err = errLogNotFound
		return
	}
	i, err := os.OpenFile(f.c.Index, os.O_RDONLY, 0666)
	if err != nil {
		log.Error("os.OpenFile(%s) error(%v)", f.c.Index, err)
		return
	}
	defer i.Close()
	b, err := ioutil.ReadAll(i)
	if err != nil {
		log.Error("ioutil.ReadAll(%s) error(%v)", f.c.Index, err)
		return
	}
	idx = &Index{}
	if err = json.Unmarshal(b, idx); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", b, err)
		return
	}
	return
}

// nextFile return first filename.
// sorted by name.
func (f *FileCache) nextReadFile() (name string) {
	var names []string
	f.eLock.RLock()
	for name = range f.logs {
		names = append(names, name)
	}
	f.eLock.RUnlock()
	if len(names) > 0 {
		sort.Strings(names)
		name = names[0]
	}
	return
}

// loadRemain loadRemain
func (f *FileCache) loadRemain() (i *Index, w *os.File, err error) {
	if i, err = f.index(); err != nil {
		next := f.nextReadFile()
		if next == "" {
			err = errLogNotFound
			return
		}
		i = &Index{
			Name:    next,
			Updated: time.Now().Format(_formatUpdated),
		}
	}
	if w, err = f.openLog(i); err != nil {
		log.Warn("a.openLog(%v) error(%v)", i, err)
		return
	}
	return
}

// openLog open the log file
func (f *FileCache) openLog(idx *Index) (w *os.File, err error) {
	if w, err = os.OpenFile(idx.Name, os.O_RDONLY, 0666); err != nil {
		log.Error("os.OpenFile(%s) error(%v)", idx.Name, err)
		return
	}
	if _, err = w.Seek(idx.Offset, os.SEEK_SET); err != nil {
		log.Error("f.Seek(%d) error(%v)", idx.Offset, err)
		return
	}
	return
}

// watch watch
func (f *FileCache) watch() (err error) {
	if f.wh, err = fsnotify.NewWatcher(); err != nil {
		log.Error("fsnotify.NewWatcher() error(%v)", err)
		return
	}
	if err = f.wh.Add(f.c.Storage); err != nil {
		log.Error("wh.Watch(%s) error(%v)", err)
	}
	return
}

// watchproc observe the directory file changes
func (f *FileCache) watchproc() {
	var evt fsnotify.Event
	for {
		evt = <-f.wh.Events
		if evt.Op&fsnotify.Create == fsnotify.Create {
			if !strings.HasSuffix(evt.Name, f.c.Suffix) {
				log.Warn("create invalid file: %s", evt.Name)
				continue
			}
			fi, err := os.Stat(evt.Name)
			if err != nil {
				log.Error("os.Stat(%s) error(%v)", evt.Name, err)
				continue
			}
			f.eLock.Lock()
			f.logs[evt.Name] = fi
			f.eLock.Unlock()
			log.Info("create file: %s", evt.Name)
		}
		if evt.Op&fsnotify.Remove == fsnotify.Remove {
			f.eLock.Lock()
			delete(f.logs, evt.Name)
			f.eLock.Unlock()
			log.Info("remove file: %s", evt.Name)
		}
	}
}

// setIndex setIndex
func (f *FileCache) setIndex(idx *Index) (err error) {
	w, err := os.OpenFile(f.c.Index, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error("os.OpenFile(%s) error(%v)", f.c.Index, err)
		return
	}
	defer w.Close()
	b, err := json.Marshal(idx)
	if err != nil {
		log.Error("json.Marshal(%v)", idx)
		return
	}
	if _, err = w.Write(b); err != nil {
		log.Error("f.Write(%s) error(%v)", b, err)
	}
	return
}

// tailLog check the log format and get log from reader
func (f *FileCache) tailLog(rr *bufio.Reader) (b []byte, err error) {
	var (
		t []byte
	)
	// peek magic
	for {
		if b, err = rr.Peek(_logMagicSize); err != nil {
			return
		}
		if bytes.Equal(b, logMagic) {
			break
		}
		rr.Discard(1)
	}
	// peek length
	if t, err = rr.Peek(_logHeadSize); err != nil {
		if err != io.EOF {
			log.Error("rr.Peek(len:%d) error(%v)", _logLenSize, err)
		}
		return
	}
	// peek body
	l := int(binary.BigEndian.Uint32(t[_logMagicSize:_logHeadSize]))
	if t, err = rr.Peek(_logHeadSize + l); err != nil {
		if err != io.EOF {
			log.Error("rr.Peek(%d) error(%v)", l, err)
		}
		return
	}
	b = t[_logHeadSize:]
	rr.Discard(l + _logHeadSize)
	return
}

// readproc read data and encapsulation protocol from file
func (f *FileCache) readProcess() {
	var (
		err      error
		idx      *Index
		rr       = bufio.NewReaderSize(nil, f.c.ReadBuffer)
		lastTime int64
		length   int
		cur      *os.File
	)
	if idx, cur, err = f.loadRemain(); err == nil {
		rr.Reset(cur)
	}
	for {
		if time.Now().Unix()-lastTime > 5 {
			if idx != nil {
				f.setIndex(idx)
			}
			lastTime = time.Now().Unix()
		}
		f.eLock.RLock()
		length = len(f.logs)
		f.eLock.RUnlock()
		// check is available for observing file
		if length == 0 {
			if cur != nil {
				cur.Close()
				cur = nil
			}
			time.Sleep(time.Second * 1)
			continue
		}
		// read first file from observing logs
		if cur == nil {
			next := f.nextReadFile()
			idx = &Index{
				Name:    next,
				Updated: time.Now().Format(_formatUpdated),
			}
			if cur, err = f.openLog(idx); err != nil {
				log.Error("a.openLog(%v) error(%v)", idx, err)
				continue
			}
			rr.Reset(cur)
			f.setIndex(idx)
		}
		// tail a log from thos.OpenFilee buffer
		b, err := f.tailLog(rr)
		if err != nil {
			if err == io.EOF {
				if length > 1 {
					cur.Close()
					cur = nil
					os.Remove(idx.Name)
					f.eLock.Lock()
					delete(f.logs, idx.Name)
					f.eLock.Unlock()
				} else {
					time.Sleep(time.Second * 1)
				}
				continue
			}
			log.Error("read log error(%v)", err)
			rr.Discard(1)
			continue
		}

		idx.Offset += int64(len(b)) + _logHeadSize

		if len(b) <= _logLancerHeaderLen {
			continue
		}
		e := event.GetEvent()
		e.Write(b[_logLancerHeaderLen:])
		e.LogId = string(b[:_logIdSize])
		f.readChan <- e
	}
}

// check storage size
func (f *FileCache) checkStorageSize() {
	var size int64
	if entries, err := ioutil.ReadDir(f.c.Storage); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				size += entry.Size()
			}
		}
	}
	if size > int64(f.c.StorageMaxMB*1024*1024) {
		log.Error("storage is full, discard log")
		flowmonitor.Fm.Add("log-agent", "log-agent.output.file-cache", strconv.FormatInt(time.Now().Unix()/100*100, 10), "ERROR", "storage full")
		f.storageFull = true
	} else {
		f.storageFull = false
	}
}

func (f *FileCache) nextFileName() string {
	return path.Join(f.c.Storage, strconv.FormatInt(time.Now().Unix(), 10)+f.c.Suffix)
}

// nextFile set first log filename.
// sorted by name.
func (f *FileCache) nextFile() (err error) {
	f.next <- f.nextFileName()
	return
}
