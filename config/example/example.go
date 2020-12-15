package example

import "github.com/golang/protobuf/ptypes/wrappers"

type Option func(e *ExampleClient) error

func (c *ExampleConfig) Options() []Option {
	return []Option{
		ApplyPassword(c.Password),
		ApplyTimeout(c.Timeout),
	}
}

type ExampleClient struct {
	addr     string
	password string
	timeout  int64
}

func ApplyPassword(v *wrappers.StringValue) func(e *ExampleClient) error {
	return func(e *ExampleClient) error {
		if v != nil {
			e.password = v.Value
		} else {
			e.password = "dangerous"
		}
		return nil
	}
}

func ApplyTimeout(v *wrappers.Int64Value) func(e *ExampleClient) error {
	return func(e *ExampleClient) error {
		if v != nil {
			e.timeout = v.Value
		} else {
			e.timeout = 1000
		}
		return nil
	}
}

func New(addr string, options ...Option) (*ExampleClient, error) {
	e := &ExampleClient{
		addr: addr,
	}
	for _, opt := range options {
		err := opt(e)
		if err != nil {
			panic(err)
		}
	}
	return e, nil
}
