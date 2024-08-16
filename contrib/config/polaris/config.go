package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/polarismesh/polaris-go"

	"github.com/go-kratos/kratos/v2/config"
)

// Option is polaris config option.
type Option func(o *options)

type options struct {
	namespace  string
	fileGroup  string
	fileName   string
	configFile polaris.ConfigFile
}

// WithNamespace with polaris config namespace
func WithNamespace(namespace string) Option {
	return func(o *options) {
		o.namespace = namespace
	}
}

// WithFileGroup with polaris config fileGroup
func WithFileGroup(fileGroup string) Option {
	return func(o *options) {
		o.fileGroup = fileGroup
	}
}

// WithFileName with polaris config fileName
func WithFileName(fileName string) Option {
	return func(o *options) {
		o.fileName = fileName
	}
}

type source struct {
	client  polaris.ConfigAPI
	options *options
}

func New(client polaris.ConfigAPI, opts ...Option) (config.Source, error) {
	options := &options{
		namespace: "default",
		fileGroup: "",
		fileName:  "",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.fileGroup == "" {
		return nil, errors.New("fileGroup invalid")
	}

	if options.fileName == "" {
		return nil, errors.New("fileName invalid")
	}

	return &source{
		client:  client,
		options: options,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	configFile, err := s.client.GetConfigFile(s.options.namespace, s.options.fileGroup, s.options.fileName)
	if err != nil {
		fmt.Println("fail to get config.", err)
		return nil, err
	}

	content := configFile.GetContent()
	k := s.options.fileName

	s.options.configFile = configFile

	return []*config.KeyValue{
		{
			Key:    k,
			Value:  []byte(content),
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		},
	}, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s.options.configFile), nil
}
