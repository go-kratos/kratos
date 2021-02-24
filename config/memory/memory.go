package memory

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/google/uuid"
	"sync"
)

var _ config.Source = (*memory)(nil)

type memory struct {
	sync.RWMutex
	KeyValue *config.KeyValue
	Watchers map[string]*watcher
}

func NewSource(opts ...config.Option) config.Source {
	var options config.Options
	for _, o := range opts {
		o(&options)
	}

	s := &memory{
		Watchers: make(map[string]*watcher),
	}

	if options.Context != nil {
		c, ok := options.Context.Value(keyValueKey{}).(*config.KeyValue)
		if ok {
			s.Update(c)
		}
	}

	return s
}

func (m *memory) Load() ([]*config.KeyValue, error) {
	m.RLock()
	kv := &config.KeyValue{
		Key:      m.KeyValue.Key,
		Value:    m.KeyValue.Value,
		Metadata: m.KeyValue.Metadata,
	}
	m.RUnlock()
	return []*config.KeyValue{kv}, nil
}

func (m *memory) Watch() (config.Watcher, error) {
	w := &watcher{
		Id:      uuid.New().String(),
		Updates: make(chan *config.KeyValue, 100),
		Source:  m,
	}

	m.Lock()
	m.Watchers[w.Id] = w
	m.Unlock()
	return w, nil
}

func (m *memory) Write(cs *config.KeyValue) error {
	m.Update(cs)
	return nil
}

// Update allows manual updates of the config data.
func (m *memory) Update(c *config.KeyValue) {
	if c == nil {
		return
	}

	m.Lock()

	m.KeyValue = &config.KeyValue{
		Key:      c.Key,
		Value:    c.Value,
		Metadata: c.Metadata,
	}

	// update watchers
	for _, w := range m.Watchers {
		select {
		case w.Updates <- m.KeyValue:
		default:
		}
	}
	m.Unlock()
}
