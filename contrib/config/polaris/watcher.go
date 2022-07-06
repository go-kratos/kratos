package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"
)

type Watcher struct {
	configFile polaris.ConfigFile
	fullPath   string
}

type eventChan struct {
	closed bool
	event  chan model.ConfigFileChangeEvent
}

var eventChanMap = make(map[string]eventChan)

func getFullPath(namespace string, fileGroup string, fileName string) string {
	return fmt.Sprintf("%s/%s/%s", namespace, fileGroup, fileName)
}

func recieve(event model.ConfigFileChangeEvent) {
	meta := event.ConfigFileMetadata
	ec := eventChanMap[getFullPath(meta.GetNamespace(), meta.GetFileGroup(), meta.GetFileName())]
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	if !ec.closed {
		ec.event <- event
	}
}

func newWatcher(configFile polaris.ConfigFile) *Watcher {
	configFile.AddChangeListener(recieve)

	fullPath := getFullPath(configFile.GetNamespace(), configFile.GetFileGroup(), configFile.GetFileName())
	if _, ok := eventChanMap[fullPath]; !ok {
		eventChanMap[fullPath] = eventChan{
			closed: false,
			event:  make(chan model.ConfigFileChangeEvent),
		}
	}
	w := &Watcher{
		configFile: configFile,
		fullPath:   fullPath,
	}
	return w
}

func (w *Watcher) Next() ([]*config.KeyValue, error) {
	ec := eventChanMap[w.fullPath]
	event := <-ec.event
	return []*config.KeyValue{
		{
			Key:    w.configFile.GetFileName(),
			Value:  []byte(event.NewValue),
			Format: strings.TrimPrefix(filepath.Ext(w.configFile.GetFileName()), "."),
		},
	}, nil
}

func (w *Watcher) Stop() error {
	ec := eventChanMap[w.fullPath]
	if !ec.closed {
		ec.closed = true
		close(ec.event)
	}
	return nil
}
