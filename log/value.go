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

func bindValues(pairs []interface{}) []interface{} {
	for i := 1; i < len(pairs); i += 2 {
		if v, ok := pairs[i].(Valuer); ok {
			pairs[i] = v()
		}
	}
	return pairs
}
