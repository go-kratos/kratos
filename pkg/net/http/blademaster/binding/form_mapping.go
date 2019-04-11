package binding

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// scache struct reflect type cache.
var scache = &cache{
	data: make(map[reflect.Type]*sinfo),
}

type cache struct {
	data  map[reflect.Type]*sinfo
	mutex sync.RWMutex
}

func (c *cache) get(obj reflect.Type) (s *sinfo) {
	var ok bool
	c.mutex.RLock()
	if s, ok = c.data[obj]; !ok {
		c.mutex.RUnlock()
		s = c.set(obj)
		return
	}
	c.mutex.RUnlock()
	return
}

func (c *cache) set(obj reflect.Type) (s *sinfo) {
	s = new(sinfo)
	tp := obj.Elem()
	for i := 0; i < tp.NumField(); i++ {
		fd := new(field)
		fd.tp = tp.Field(i)
		tag := fd.tp.Tag.Get("form")
		fd.name, fd.option = parseTag(tag)
		if defV := fd.tp.Tag.Get("default"); defV != "" {
			dv := reflect.New(fd.tp.Type).Elem()
			setWithProperType(fd.tp.Type.Kind(), []string{defV}, dv, fd.option)
			fd.hasDefault = true
			fd.defaultValue = dv
		}
		s.field = append(s.field, fd)
	}
	c.mutex.Lock()
	c.data[obj] = s
	c.mutex.Unlock()
	return
}

type sinfo struct {
	field []*field
}

type field struct {
	tp     reflect.StructField
	name   string
	option tagOptions

	hasDefault   bool          // if field had default value
	defaultValue reflect.Value // field default value
}

func mapForm(ptr interface{}, form map[string][]string) error {
	sinfo := scache.get(reflect.TypeOf(ptr))
	val := reflect.ValueOf(ptr).Elem()
	for i, fd := range sinfo.field {
		typeField := fd.tp
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := fd.name
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			// if "form" tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			if structFieldKind == reflect.Struct {
				err := mapForm(structField.Addr().Interface(), form)
				if err != nil {
					return err
				}
				continue
			}
		}
		inputValue, exists := form[inputFieldName]
		if !exists {
			// Set the field as default value when the input value is not exist
			if fd.hasDefault {
				structField.Set(fd.defaultValue)
			}
			continue
		}
		// Set the field as default value when the input value is empty
		if fd.hasDefault && inputValue[0] == "" {
			structField.Set(fd.defaultValue)
			continue
		}
		if _, isTime := structField.Interface().(time.Time); isTime {
			if err := setTimeField(inputValue[0], typeField, structField); err != nil {
				return err
			}
			continue
		}
		if err := setWithProperType(typeField.Type.Kind(), inputValue, structField, fd.option); err != nil {
			return err
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val []string, structField reflect.Value, option tagOptions) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val[0], 0, structField)
	case reflect.Int8:
		return setIntField(val[0], 8, structField)
	case reflect.Int16:
		return setIntField(val[0], 16, structField)
	case reflect.Int32:
		return setIntField(val[0], 32, structField)
	case reflect.Int64:
		return setIntField(val[0], 64, structField)
	case reflect.Uint:
		return setUintField(val[0], 0, structField)
	case reflect.Uint8:
		return setUintField(val[0], 8, structField)
	case reflect.Uint16:
		return setUintField(val[0], 16, structField)
	case reflect.Uint32:
		return setUintField(val[0], 32, structField)
	case reflect.Uint64:
		return setUintField(val[0], 64, structField)
	case reflect.Bool:
		return setBoolField(val[0], structField)
	case reflect.Float32:
		return setFloatField(val[0], 32, structField)
	case reflect.Float64:
		return setFloatField(val[0], 64, structField)
	case reflect.String:
		structField.SetString(val[0])
	case reflect.Slice:
		if option.Contains("split") {
			val = strings.Split(val[0], ",")
		}
		filtered := filterEmpty(val)
		switch structField.Type().Elem().Kind() {
		case reflect.Int64:
			valSli := make([]int64, 0, len(filtered))
			for i := 0; i < len(filtered); i++ {
				d, err := strconv.ParseInt(filtered[i], 10, 64)
				if err != nil {
					return err
				}
				valSli = append(valSli, d)
			}
			structField.Set(reflect.ValueOf(valSli))
		case reflect.String:
			valSli := make([]string, 0, len(filtered))
			for i := 0; i < len(filtered); i++ {
				valSli = append(valSli, filtered[i])
			}
			structField.Set(reflect.ValueOf(valSli))
		default:
			sliceOf := structField.Type().Elem().Kind()
			numElems := len(filtered)
			slice := reflect.MakeSlice(structField.Type(), len(filtered), len(filtered))
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, filtered[i:], slice.Index(i), ""); err != nil {
					return err
				}
			}
			structField.Set(slice)
		}
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return errors.WithStack(err)
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return errors.WithStack(err)
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return nil
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return errors.WithStack(err)
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		return errors.New("Blank time format")
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return errors.WithStack(err)
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return errors.WithStack(err)
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func filterEmpty(val []string) []string {
	filtered := make([]string, 0, len(val))
	for _, v := range val {
		if v != "" {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
