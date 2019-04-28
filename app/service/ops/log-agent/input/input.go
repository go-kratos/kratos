package input

import (
	"fmt"
	"context"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/event"
)

type Input interface {
	Run() (err error)
	Stop()
	Ctx() (ctx context.Context)
}

// Factory is used to register functions creating new Input instances.
type Factory = func(ctx context.Context, config interface{}, connector chan<- *event.ProcessorEvent) (Input, error)

var registry = make(map[string]Factory)

func Register(name string, factory Factory) error {
	log.Info("Registering input factory")
	if name == "" {
		return fmt.Errorf("Error registering input: name cannot be empty")
	}
	if factory == nil {
		return fmt.Errorf("Error registering input '%v': factory cannot be empty", name)
	}
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering input '%v': already registered", name)
	}

	registry[name] = factory
	log.Info("Successfully registered input")

	return nil
}

func GetFactory(name string) (Factory, error) {
	if _, exists := registry[name]; !exists {
		return nil, fmt.Errorf("Error creating input. No such input type exist: '%v'", name)
	}
	return registry[name], nil
}
