package selector

import (
	"reflect"
)

func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Ptr:
		elem := v.Elem()
		if !elem.IsValid() {
			return true
		}
		return IsNil(elem.Interface())
	default:
		return false
	}
}
