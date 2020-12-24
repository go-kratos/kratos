package config

type config struct {
	resolvers []Resolver
}

func (c *config) Value(key string) (Value, bool) {
	for _, r := range c.resolvers {
		v, ok := r.Value(key)
		if ok {
			return v, true
		}
	}
	return nil, false
}
