package dockerlogcollector

import (
	"context"
	"time"
	"strings"
	"path"
	"io/ioutil"

	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pipeline"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerLogCollector struct {
	c      *Config
	client *client.Client
	ctx    context.Context
	cancel context.CancelFunc
}

type configItem struct {
	configPath string
	MergedDir  string
}

func InitDockerLogCollector(ctx context.Context, c *Config) (err error) {
	if err = c.ConfigValidate(); err != nil {
		return err
	}
	collector := new(DockerLogCollector)
	collector.c = c

	collector.ctx, collector.cancel = context.WithCancel(ctx)

	// init docker client
	collector.client, err = client.NewEnvClient()
	if err != nil {
		return err
	}
	go collector.scan()

	return nil
}

func (collector *DockerLogCollector) getConfigs() ([]*configItem, error) {
	var (
		configItems = make([]*configItem, 0)
		mergedDir   string
		ok          bool
	)
	containers, err := collector.client.ContainerList(collector.ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		info, err := collector.client.ContainerInspect(collector.ctx, container.ID)
		if err != nil {
			log.Error("failed to inspect container: %s", container.ID)
			continue
		}
		// get overlay2 info
		if info.GraphDriver.Name != "overlay2" {
			log.Error("only overlay2 is supported")
			continue
		}

		mergedDir, ok = info.GraphDriver.Data["MergedDir"]
		if !ok {
			log.Error("failed to get MergedDir of container:%s", container.ID)
		}

		for _, env := range info.Config.Env {
			if strings.HasPrefix(env, collector.c.

				ConfigEnv) {
				for _, path := range strings.Split(strings.TrimPrefix(env, collector.c.ConfigEnv+"="), ",") {
					configItems = append(configItems, &configItem{path, mergedDir})
				}
			}
		}
	}

	return configItems, nil
}

func (collector *DockerLogCollector) scan() {
	ticker := time.Tick(time.Duration(collector.c.ScanInterval))
	for {
		select {
		case <-ticker:
			configItems, err := collector.getConfigs()
			if err != nil {
				log.Error("failed to scan hostlogcollector config file list: %s", err)
				continue
			}
			for _, item := range configItems {
				configPath := path.Join(item.MergedDir, item.configPath)
				config, err := ioutil.ReadFile(configPath)
				if err != nil {
					log.Error("filed to read hostlogcollector config file %s: %s", configPath, err)
					continue
				}
				if !pipeline.PipelineManagement.PipelineExisted(configPath) {
					ctx := context.WithValue(collector.ctx, "MergedDir", item.MergedDir)
					go pipeline.PipelineManagement.StartPipeline(ctx, configPath, string(config))
				}
			}
		case <-collector.ctx.Done():
			return
		}
	}
}
