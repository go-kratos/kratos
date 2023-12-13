package consul

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	source          *source
	ch              chan []*config.KeyValue
	wp              *watch.Plan
	fileModifyIndex map[string]uint64
	ctx             context.Context
	cancel          context.CancelFunc
}

func (w *watcher) handle(_ uint64, data interface{}) {
	if data == nil {
		return
	}

	kv, ok := data.(api.KVPairs)
	if !ok {
		return
	}

	pathPrefix := w.source.options.path
	if !strings.HasSuffix(w.source.options.path, "/") {
		pathPrefix = pathPrefix + "/"
	}
	kvs := make([]*config.KeyValue, 0, len(kv))
	for _, item := range kv {
		if index, ok := w.fileModifyIndex[item.Key]; ok && item.ModifyIndex == index {
			continue
		}
		k := strings.TrimPrefix(item.Key, pathPrefix)
		if k == "" {
			continue
		}
		kvs = append(kvs, &config.KeyValue{
			Key:    k,
			Value:  item.Value,
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		})
		w.fileModifyIndex[item.Key] = item.ModifyIndex
	}

	if len(kvs) == 0 {
		return
	}

	w.ch <- kvs
}

func newWatcher(s *source) (*watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	w := &watcher{
		source:          s,
		ch:              make(chan []*config.KeyValue),
		fileModifyIndex: make(map[string]uint64),
		ctx:             ctx,
		cancel:          cancel,
	}

	wp, err := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": s.options.path})
	if err != nil {
		return nil, err
	}

	wp.Handler = w.handle
	w.wp = wp

	// wp.Run is a blocking call and will prevent newWatcher from returning
	go func() {
		err := wp.RunWithClientAndHclog(s.client, nil)
		if err != nil {
			panic(err)
		}
	}()

	return w, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case kv := <-w.ch:
		return kv, nil
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	w.wp.Stop()
	w.cancel()
	return nil
}
