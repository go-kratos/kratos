package paladin

import "time"

// Bool return bool value.
func Bool(v *Value, def bool) bool {
	b, err := v.Bool()
	if err != nil {
		return def
	}
	return b
}

// Int return int value.
func Int(v *Value, def int) int {
	i, err := v.Int()
	if err != nil {
		return def
	}
	return i
}

// Int32 return int32 value.
func Int32(v *Value, def int32) int32 {
	i, err := v.Int32()
	if err != nil {
		return def
	}
	return i
}

// Int64 return int64 value.
func Int64(v *Value, def int64) int64 {
	i, err := v.Int64()
	if err != nil {
		return def
	}
	return i
}

// Float32 return float32 value.
func Float32(v *Value, def float32) float32 {
	f, err := v.Float32()
	if err != nil {
		return def
	}
	return f
}

// Float64 return float32 value.
func Float64(v *Value, def float64) float64 {
	f, err := v.Float64()
	if err != nil {
		return def
	}
	return f
}

// String return string value.
func String(v *Value, def string) string {
	s, err := v.String()
	if err != nil {
		return def
	}
	return s
}

// Duration parses a duration string. A duration string is a possibly signed sequence of decimal numbers
// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func Duration(v *Value, def time.Duration) time.Duration {
	dur, err := v.Duration()
	if err != nil {
		return def
	}
	return dur
}
