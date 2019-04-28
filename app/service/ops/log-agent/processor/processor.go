package processor

import (
	"fmt"
	"context"

	"go-common/app/service/ops/log-agent/event"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/output"
)

// Factory is used to register functions creating new output instances.
type Factory = func(cxt context.Context, config interface{}, input <-chan *event.ProcessorEvent) (chan *event.ProcessorEvent, error)

var registry = make(map[string]Factory)

func Register(name string, factory Factory) error {
	log.Info("Registering  processor factory")
	if name == "" {
		return fmt.Errorf("Error registering processor: name cannot be empty")
	}
	if factory == nil {
		return fmt.Errorf("Error registering processor '%v': factory cannot be empty", name)
	}
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering processor '%v': already registered", name)
	}

	registry[name] = factory
	log.Info("Successfully registered processor: '%v'", name)

	return nil
}

func GetFactory(name string) (Factory, error) {
	if _, exists := registry[name]; !exists {
		return nil, fmt.Errorf("Error creating processor. No such processor type exist: '%v'", name)
	}
	return registry[name], nil
}

func WriteToOutput(ctx context.Context, dest string, input <-chan *event.ProcessorEvent) (err error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-input:
				if dest != "" {
					e.Destination = dest
				}
				outputChan, err := output.GetOutputChan(e.Destination)
				if err != nil {
					log.Error("failed to get output chan:%s; discard log", err)
					continue
				}
				outputChan <- e
			}
		}
	}()
	return nil
}
