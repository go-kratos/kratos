package paladin

import (
	"context"
	"errors"
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
	baseDir  string
	values   *Map
	watchChs map[string][]chan Event
	mx       sync.Mutex
	wg       sync.WaitGroup
	done     chan struct{}
}

// NewFile new a config file client.
// conf = /data/conf/app/
// conf = /data/conf/app/xxx.toml
func NewFile(base string) (Client, error) {
	base = filepath.FromSlash(base)
	paths, err := readAllPaths(base)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, errors.New("empty config path")
	}
	raws, err := loadValuesFromPaths(paths)
	if err != nil {
		return nil, err
	}
	values := new(Map)
	values.Store(raws)
	fc := &file{
		baseDir:  base,
		values:   values,
		watchChs: make(map[string][]chan Event),
		done:     make(chan struct{}, 1),
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
		log.Printf("paladin: create file watcher fail! reload function will lose efficacy error: %s", err)
		return
	}
	if err = fswatcher.Add(f.baseDir); err != nil {
		log.Printf("paladin: create fsnotify for base path %s fail %s, reload function will lose efficacy", f.baseDir, err)
		return
	}
	log.Printf("paladin: start watch config: %s", f.baseDir)
	for event := range fswatcher.Events {
		// use vim edit config will trigger rename
		switch {
		case event.Op&fsnotify.Write == fsnotify.Write, event.Op&fsnotify.Create == fsnotify.Create:
			f.reloadFile(event.Name)
		default:
			log.Printf("paladin: unsupport event %s ingored", event)
		}
	}
}

func (f *file) reloadFile(name string) {
	// NOTE: in some case immediately read file content after receive event
	// will get old content, sleep 100ms make sure get correct content.
	time.Sleep(100 * time.Millisecond)
	key := filepath.Base(name)
	value, err := loadValue(name)
	if err != nil {
		log.Printf("paladin: load file: %s error: %s, skipped", name, err)
		return
	}
	raws := f.values.Load()
	raws[name] = value
	f.values.Store(raws)
	f.mx.Lock()
	chs := f.watchChs[key]
	f.mx.Unlock()
	for _, ch := range chs {
		select {
		case ch <- Event{Event: EventUpdate, Value: value.raw}:
		default:
			log.Printf("paladin: event channel full discard file %s update event", name)
		}
	}
	log.Printf("paladin: reload config: %s notify: %d\n", name, len(chs))
}

func readAllPaths(base string) ([]string, error) {
	fi, err := os.Stat(base)
	if err != nil {
		return nil, fmt.Errorf("paladin: check local config file fail! error: %s", err)
	}
	var paths []string
	if fi.IsDir() {
		files, err := ioutil.ReadDir(base)
		if err != nil {
			return nil, fmt.Errorf("paladin: read dir %s error: %s", base, err)
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
