package paladin

import (
	"context"
	"flag"
	"go-common/library/log"
)

var (
	// DefaultClient default client.
	DefaultClient Client
	confPath      string
	vars          = make(map[string][]Setter) // NOTE: no thread safe
)

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config client.
func Init() (err error) {
	if confPath != "" {
		DefaultClient, err = NewFile(confPath)
	} else {
		DefaultClient, err = NewSven()
	}
	if err != nil {
		return
	}
	go func() {
		for event := range DefaultClient.WatchEvent(context.Background()) {
			if event.Event != EventUpdate && event.Event != EventAdd {
				continue
			}
			if sets, ok := vars[event.Key]; ok {
				for _, s := range sets {
					if err := s.Set(event.Value); err != nil {
						log.Error("paladin: vars:%v event:%v error(%v)", s, event, err)
					}
				}
			}
		}
	}()
	return
}

// Watch watch on a key. The configuration implements the setter interface, which is invoked when the configuration changes.
func Watch(key string, s Setter) error {
	v := DefaultClient.Get(key)
	str, err := v.Raw()
	if err != nil {
		return err
	}
	if err := s.Set(str); err != nil {
		return err
	}
	vars[key] = append(vars[key], s)
	return nil
}

// WatchEvent watch on multi keys. Events are returned when the configuration changes.
func WatchEvent(ctx context.Context, keys ...string) <-chan Event {
	return DefaultClient.WatchEvent(ctx, keys...)
}

// Get return value by key.
func Get(key string) *Value {
	return DefaultClient.Get(key)
}

// GetAll return all config map.
func GetAll() *Map {
	return DefaultClient.GetAll()
}

// Keys return values key.
func Keys() []string {
	return DefaultClient.GetAll().Keys()
}

// Close close watcher.
func Close() error {
	return DefaultClient.Close()
}
