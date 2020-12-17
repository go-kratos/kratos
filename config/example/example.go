package example

type Option func(*exampleOptions) error

type exampleOptions struct {
	password string
	timeout  int64
}

var defaultExampleOptions = exampleOptions{
	password: "dangerous",
	timeout:  100,
}

func ApplyOptions(c *ExampleConfig) []Option {
	opts := make([]Option, 0)
	if c.Password != nil {
		opts = append(opts, func(o *exampleOptions) error {
			o.password = c.Password.Value
			return nil
		})
	}
	if c.Timeout != nil {
		opts = append(opts, func(o *exampleOptions) error {
			o.timeout = c.Timeout.Value
			return nil
		})
	}
	return opts
}

type ExampleClient struct {
	addr string
	opts exampleOptions
}

func New(addr string, options ...Option) (*ExampleClient, error) {
	e := &ExampleClient{
		addr: addr,
	}
	opts := defaultExampleOptions
	for _, o := range options {
		err := o(&opts)
		if err != nil {
			panic(err)
		}
	}
	e.opts = opts
	return e, nil
}
