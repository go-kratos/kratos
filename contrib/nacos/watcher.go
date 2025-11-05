package nacos

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/go-kratos/kratos/v2/config"
)

type ConfigWatcher struct {
	dataID             string
	group              string
	content            chan string
	cancelListenConfig cancelListenConfigFunc

	ctx    context.Context
	cancel context.CancelFunc
}

type cancelListenConfigFunc func(params vo.ConfigParam) (err error)

func newWatcher(ctx context.Context, dataID string, group string, cancelListenConfig cancelListenConfigFunc) *ConfigWatcher {
	ctx, cancel := context.WithCancel(ctx)
	w := &ConfigWatcher{
		dataID:             dataID,
		group:              group,
		cancelListenConfig: cancelListenConfig,
		content:            make(chan string, 100),

		ctx:    ctx,
		cancel: cancel,
	}
	return w
}

func (w *ConfigWatcher) Next() ([]*config.KeyValue, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case content := <-w.content:
		k := w.dataID
		return []*config.KeyValue{
			{
				Key:    k,
				Value:  []byte(content),
				Format: strings.TrimPrefix(filepath.Ext(k), "."),
			},
		}, nil
	}
}

func (w *ConfigWatcher) Close() error {
	err := w.cancelListenConfig(vo.ConfigParam{
		DataId: w.dataID,
		Group:  w.group,
	})
	w.cancel()
	return err
}

func (w *ConfigWatcher) Stop() error {
	return w.Close()
}
