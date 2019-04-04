package paladin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	defaultChSize = 10
)

var _ Client = &file{}

// file is file config client.
type file struct {
	values *Map
	rawVal map[string]*Value

	watchChs map[string][]chan Event
	mx       sync.Mutex
	wg       sync.WaitGroup

	base string
	done chan struct{}
}

func readAllPaths(base string) ([]string, error) {
	fi, err := os.Stat(base)
	if err != nil {
		return nil, fmt.Errorf("check local config file fail! error: %s", err)
	}
	// dirs or file to paths
	var paths []string
	if fi.IsDir() {
		files, err := ioutil.ReadDir(base)
		if err != nil {
			return nil, fmt.Errorf("read dir %s error: %s", base, err)
		}
		for _, file := range files {
			if !file.IsDir() {
				paths = append(paths, path.Join(base, file.Name()))
			}
		}
	} else {
		paths = append(paths, base)
	}
	return paths, nil
}

func loadValuesFromPaths(paths []string) (map[string]*Value, error) {
	// laod config file to values
	var err error
	values := make(map[string]*Value, len(paths))
	for _, fpath := range paths {
		if values[path.Base(fpath)], err = loadValue(fpath); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func loadValue(fpath string) (*Value, error) {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	content := string(data)
	return &Value{val: content, raw: content}, nil
}

// NewFile new a config file client.
// conf = /data/conf/app/
// conf = /data/conf/app/xxx.toml
func NewFile(base string) (Client, error) {
	// paltform slash
	base = filepath.FromSlash(base)

	paths, err := readAllPaths(base)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("empty config path")
	}

	rawVal, err := loadValuesFromPaths(paths)
	if err != nil {
		return nil, err
	}

	valMap := &Map{}
	valMap.Store(rawVal)
	fc := &file{
		values:   valMap,
		rawVal:   rawVal,
		watchChs: make(map[string][]chan Event),

		base: base,
		done: make(chan struct{}, 1),
	}

	fc.wg.Add(1)
	go fc.daemon()

	return fc, nil
}

// Get return value by key.
func (f *file) Get(key string) *Value {
	return f.values.Get(key)
}

// GetAll return value map.
func (f *file) GetAll() *Map {
	return f.values
}

// WatchEvent watch multi key.
func (f *file) WatchEvent(ctx context.Context, keys ...string) <-chan Event {
	f.mx.Lock()
	defer f.mx.Unlock()
	ch := make(chan Event, defaultChSize)
	for _, key := range keys {
		f.watchChs[key] = append(f.watchChs[key], ch)
	}
	return ch
}

// Close close watcher.
func (f *file) Close() error {
	f.done <- struct{}{}
	f.wg.Wait()
	return nil
}

// file config daemon to watch file modification
func (f *file) daemon() {
	defer f.wg.Done()
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("create file watcher fail! reload function will lose efficacy error: %s", err)
		return
	}
	if err = fswatcher.Add(f.base); err != nil {
		log.Printf("create fsnotify for base path %s fail %s, reload function will lose efficacy", f.base, err)
		return
	}
	log.Printf("start watch filepath: %s", f.base)
	for event := range fswatcher.Events {
		switch event.Op {
		// use vim edit config will trigger rename
		case fsnotify.Write, fsnotify.Create:
			f.reloadFile(event.Name)
		case fsnotify.Chmod:
		default:
			log.Printf("unsupport event %s ingored", event)
		}
	}
}

func (f *file) reloadFile(name string) {
	// NOTE: in some case immediately read file content after receive event
	// will get old content, sleep 100ms make sure get correct content.
	time.Sleep(100 * time.Millisecond)
	key := filepath.Base(name)
	val, err := loadValue(name)
	if err != nil {
		log.Printf("load file %s error: %s, skipped", name, err)
		return
	}
	f.rawVal[key] = val
	f.values.Store(f.rawVal)

	f.mx.Lock()
	chs := f.watchChs[key]
	f.mx.Unlock()

	for _, ch := range chs {
		select {
		case ch <- Event{Event: EventUpdate, Value: val.raw}:
		default:
			log.Printf("event channel full discard file %s update event", name)
		}
	}
}
