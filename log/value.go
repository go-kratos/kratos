package log

// Valuer is returns a log value.
type Valuer func() interface{}

// Value return the function value.
func Value(v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v()
	}
	return v
}
