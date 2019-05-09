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

var _ Client = &file{}

type watcher struct {
	keys []string
	C    chan Event
}

func newWatcher(keys []string) *watcher {
	return &watcher{keys: keys, C: make(chan Event, 5)}
}

func (w *watcher) HasKey(key string) bool {
	if len(w.keys) == 0 {
		return true
	}
	for _, k := range w.keys {
		if keyNamed(k) == key {
			return true
		}
	}
	return false
}

func (w *watcher) Handle(event Event) {
	select {
	case w.C <- event:
	default:
		log.Printf("paladin: event channel full discard file %s update event", event.Key)
	}
}

// file is file config client.
type file struct {
	values   *Map
	wmu      sync.RWMutex
	notify   *fsnotify.Watcher
	watchers map[*watcher]struct{}
}

// NewFile new a config file client.
// conf = /data/conf/app/
// conf = /data/conf/app/xxx.toml
func NewFile(base string) (Client, error) {
	base = filepath.FromSlash(base)
	raws, err := loadValues(base)
	if err != nil {
		return nil, err
	}
	notify, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	values := new(Map)
	values.Store(raws)
	f := &file{
		values:   values,
		notify:   notify,
		watchers: make(map[*watcher]struct{}),
	}
	go f.watchproc(base)
	return f, nil
}

// Get return value by key.
func (f *file) Get(key string) *Value {
	return f.values.Get(key)
}

// GetAll return value map.
func (f *file) GetAll() *Map {
	return f.values
}

// WatchEvent watch with the specified keys.
func (f *file) WatchEvent(ctx context.Context, keys ...string) <-chan Event {
	w := newWatcher(keys)
	f.wmu.Lock()
	f.watchers[w] = struct{}{}
	f.wmu.Unlock()
	return w.C
}

// Close close watcher.
func (f *file) Close() error {
	if err := f.notify.Close(); err != nil {
		return err
	}
	f.wmu.RLock()
	for w := range f.watchers {
		close(w.C)
	}
	f.wmu.RUnlock()
	return nil
}

// file config daemon to watch file modification
func (f *file) watchproc(base string) {
	if err := f.notify.Add(base); err != nil {
		log.Printf("paladin: create fsnotify for base path %s fail %s, reload function will lose efficacy", base, err)
		return
	}
	log.Printf("paladin: start watch config: %s", base)
	for event := range f.notify.Events {
		// use vim edit config will trigger rename
		switch {
		case event.Op&fsnotify.Write == fsnotify.Write, event.Op&fsnotify.Create == fsnotify.Create:
			if err := f.reloadFile(event.Name); err != nil {
				log.Printf("paladin: load file: %s error: %s, skipped", event.Name, err)
			}
		default:
			log.Printf("paladin: unsupport event %s ingored", event)
		}
	}
}

func (f *file) reloadFile(fpath string) (err error) {
	// NOTE: in some case immediately read file content after receive event
	// will get old content, sleep 100ms make sure get correct content.
	time.Sleep(100 * time.Millisecond)
	value, err := loadValue(fpath)
	if err != nil {
		return
	}
	key := keyNamed(path.Base(fpath))
	raws := f.values.Load()
	raws[key] = value
	f.values.Store(raws)
	f.wmu.RLock()
	n := 0
	for w := range f.watchers {
		if w.HasKey(key) {
			n++
			w.Handle(Event{Event: EventUpdate, Key: key, Value: value.raw})
		}
	}
	f.wmu.RUnlock()
	log.Printf("paladin: reload config: %s events: %d\n", key, n)
	return
}

func loadValues(base string) (map[string]*Value, error) {
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
	if len(paths) == 0 {
		return nil, errors.New("empty config path")
	}
	values := make(map[string]*Value, len(paths))
	for _, fpath := range paths {
		if values[path.Base(fpath)], err = loadValue(fpath); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func loadValue(name string) (*Value, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	content := string(data)
	return &Value{val: content, raw: content}, nil
}
