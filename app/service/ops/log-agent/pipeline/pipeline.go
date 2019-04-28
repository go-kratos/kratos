package pipeline

import (
	"fmt"
	"context"

	"go-common/app/service/ops/log-agent/output/lancerlogstream"
	"go-common/app/service/ops/log-agent/output/lancergrpc"
	"go-common/app/service/ops/log-agent/input/sock"
	"go-common/app/service/ops/log-agent/input/file"
	"go-common/app/service/ops/log-agent/processor/classify"
	"go-common/app/service/ops/log-agent/processor/jsonLog"
	"go-common/app/service/ops/log-agent/processor/fileLog"
	"go-common/app/service/ops/log-agent/processor/lengthCheck"
	"go-common/app/service/ops/log-agent/processor/sample"
	"go-common/app/service/ops/log-agent/processor/httpstream"
	"go-common/app/service/ops/log-agent/processor/grok"
	"go-common/app/service/ops/log-agent/output/stdout"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
)

var inputConfigDecodeFactory = make(map[string]configDecodeFunc)
var processorConfigDecodeFactory = make(map[string]configDecodeFunc)
var outputConfigDecodeFactory = make(map[string]configDecodeFunc)

func init() {
	RegisterInputConfigDecodeFunc("sock", sock.DecodeConfig)
	RegisterInputConfigDecodeFunc("file", file.DecodeConfig)
	RegisterProcessorConfigDecodeFunc("classify", classify.DecodeConfig)
	RegisterProcessorConfigDecodeFunc("jsonLog", jsonLog.DecodeConfig)
	RegisterProcessorConfigDecodeFunc("lengthCheck", lengthCheck.DecodeConfig)
	RegisterProcessorConfigDecodeFunc("sample", sample.DecodeConfig)
	RegisterProcessorConfigDecodeFunc("httpStream", httpstream.DecodeConfig)
	RegisterProcessorConfigDecodeFunc("fileLog", fileLog.DecodeConfig)
	RegisterProcessorConfigDecodeFunc("grok", grok.DecodeConfig)
	RegisterOutputConfigDecodeFunc("stdout", stdout.DecodeConfig)
	RegisterOutputConfigDecodeFunc("lancer", lancerlogstream.DecodeConfig)
	RegisterOutputConfigDecodeFunc("lancergrpc", lancergrpc.DecodeConfig)
}

type Pipeline struct {
	c          *Config
	ctx        context.Context
	cancel     context.CancelFunc
	configPath string
	configMd5  string
}

type Config struct {
	Input     ConfigItem            `toml:"input"`
	Processor map[string]ConfigItem `toml:"processor"`
	Output    map[string]ConfigItem `toml:"output"`
}

type ConfigItem struct {
	Name   string         `toml:"type"`
	Config toml.Primitive `toml:"config"`
}

func (pipe *Pipeline) Stop() {
	pipe.cancel()
}

type configDecodeFunc = func(md toml.MetaData, primValue toml.Primitive) (c interface{}, err error)

func RegisterInputConfigDecodeFunc(name string, f configDecodeFunc) {
	inputConfigDecodeFactory[name] = f
}

func RegisterProcessorConfigDecodeFunc(name string, f configDecodeFunc) {
	processorConfigDecodeFactory[name] = f
}

func GetInputConfigDecodeFunc(name string) (f configDecodeFunc, err error) {
	if f, exist := inputConfigDecodeFactory[name]; exist {
		return f, nil
	}
	return nil, fmt.Errorf("InputConfigDecodeFunc for %s not exist", name)
}

func GetProcessorConfigDecodeFunc(name string) (f configDecodeFunc, err error) {
	if f, exist := processorConfigDecodeFactory[name]; exist {
		return f, nil
	}
	return nil, fmt.Errorf("ProcessorConfigDecodeFunc for %s not exist", name)
}

func DecodeInputConfig(name string, md toml.MetaData, primValue toml.Primitive) (c interface{}, err error) {
	dFunc, err := GetInputConfigDecodeFunc(name)
	if err != nil {
		return nil, err
	}
	return dFunc(md, primValue)
}

func DecodeProcessorConfig(name string, md toml.MetaData, primValue toml.Primitive) (c interface{}, err error) {
	dFunc, err := GetProcessorConfigDecodeFunc(name)
	if err != nil {
		return nil, err
	}
	return dFunc(md, primValue)
}

func RegisterOutputConfigDecodeFunc(name string, f configDecodeFunc) {
	outputConfigDecodeFactory[name] = f
}

func GetOutputConfigDecodeFunc(name string) (f configDecodeFunc, err error) {
	if f, exist := outputConfigDecodeFactory[name]; exist {
		return f, nil
	}
	return nil, fmt.Errorf("OutputConfigDecodeFunc for %s not exist", name)
}

func DecodeOutputConfig(name string, md toml.MetaData, primValue toml.Primitive) (c interface{}, err error) {
	dFunc, err := GetOutputConfigDecodeFunc(name)
	if err != nil {
		return nil, err
	}
	return dFunc(md, primValue)
}

func (p *Pipeline) logError(err error) {
	configPath := p.ctx.Value("configPath")
	log.Error("failed to run pipeline for %s: %s", configPath, err)
}
