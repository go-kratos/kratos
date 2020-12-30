package config

import "expvar"

func setVariable(v expvar.Var, val Value) error {
	switch vv := v.(type) {
	case *expvar.Int:
		intVal, err := val.Int()
		if err != nil {
			return err
		}
		vv.Set(intVal)
	case *expvar.Float:
		floatVal, err := val.Float()
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
