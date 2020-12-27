package config

import "expvar"

func setVar(key string, v expvar.Var, val Value) error {
	switch vv := v.(type) {
	case *expvar.Int:
		intVal, err := val.Int64()
		if err != nil {
			return err
		}
		vv.Set(intVal)
	case *expvar.Float:
		floatVal, err := val.Float64()
		if err != nil {
			return err
		}
		vv.Set(floatVal)
	case *expvar.String:
		stringVal, err := val.String()
		if err != nil {
			return err
		}
		vv.Set(stringVal)
	}
	return nil
}
