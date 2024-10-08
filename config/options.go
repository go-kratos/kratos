package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
)

// Decoder is config decoder.
type Decoder func(*KeyValue, map[string]interface{}) error

// Resolver resolve placeholder in config.
type Resolver func(map[string]interface{}) error

// Merge is config merge func.
type Merge func(dst, src interface{}) error

// Option is config option.
type Option func(*options)

type options struct {
	sources  []Source
	decoder  Decoder
	resolver Resolver
	merge    Merge
}

// WithSource with config source.
func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

// WithDecoder with config decoder.
// DefaultDecoder behavior:
// If KeyValue.Format is non-empty, then KeyValue.Value will be deserialized into map[string]interface{}
// and stored in the config cache(map[string]interface{})
// if KeyValue.Format is empty,{KeyValue.Key : KeyValue.Value} will be stored in config cache(map[string]interface{})
func WithDecoder(d Decoder) Option {
	return func(o *options) {
		o.decoder = d
	}
}

// WithResolveActualTypes with config resolver.
// bool input will enable conversion of config to data types
func WithResolveActualTypes(enableConvertToType bool) Option {
	return func(o *options) {
		o.resolver = newActualTypesResolver(enableConvertToType)
	}
}

// WithResolver with config resolver.
func WithResolver(r Resolver) Option {
	return func(o *options) {
		o.resolver = r
	}
}

// WithMergeFunc with config merge func.
func WithMergeFunc(m Merge) Option {
	return func(o *options) {
		o.merge = m
	}
}

// defaultDecoder decode config from source KeyValue
// to target map[string]interface{} using src.Format codec.
func defaultDecoder(src *KeyValue, target map[string]interface{}) error {
	if src.Format == "" {
		// expand key "aaa.bbb" into map[aaa]map[bbb]interface{}
		keys := strings.Split(src.Key, ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				target[k] = src.Value
			} else {
				sub := make(map[string]interface{})
				target[k] = sub
				target = sub
			}
		}
		return nil
	}
	if codec := encoding.GetCodec(src.Format); codec != nil {
		return codec.Unmarshal(src.Value, &target)
	}
	return fmt.Errorf("unsupported key: %s format: %s", src.Key, src.Format)
}

func newActualTypesResolver(enableConvertToType bool) func(map[string]interface{}) error {
	return func(input map[string]interface{}) error {
		mapper := mapper(input)
		return resolver(input, mapper, enableConvertToType)
	}
}

// defaultResolver resolve placeholder in map value,
// placeholder format in ${key:default}.
func defaultResolver(input map[string]interface{}) error {
	mapper := mapper(input)
	return resolver(input, mapper, false)
}

func resolver(input map[string]interface{}, mapper func(name string) string, toType bool) error {
	var resolve func(map[string]interface{}) error
	resolve = func(sub map[string]interface{}) error {
		for k, v := range sub {
			switch vt := v.(type) {
			case string:
				sub[k] = expand(vt, mapper, toType)
			case map[string]interface{}:
				if err := resolve(vt); err != nil {
					return err
				}
			case []interface{}:
				for i, iface := range vt {
					switch it := iface.(type) {
					case string:
						vt[i] = expand(it, mapper, toType)
					case map[string]interface{}:
						if err := resolve(it); err != nil {
							return err
						}
					}
				}
				sub[k] = vt
			}
		}
		return nil
	}
	return resolve(input)
}

func mapper(input map[string]interface{}) func(name string) string {
	mapper := func(name string) string {
		args := strings.SplitN(strings.TrimSpace(name), ":", 2) //nolint:mnd
		if v, has := readValue(input, args[0]); has {
			s, _ := v.String()
			return s
		} else if len(args) > 1 { // default value
			return args[1]
		}
		return ""
	}
	return mapper
}

func convertToType(input string) interface{} {
	// Check if the input is a string with quotes
	if strings.HasPrefix(input, "\"") && strings.HasSuffix(input, "\"") {
		// Trim the quotes and return the string value
		return strings.Trim(input, "\"")
	}

	// Try converting to bool
	if input == "true" || input == "false" {
		b, _ := strconv.ParseBool(input)
		return b
	}

	// Try converting to float64
	if strings.Contains(input, ".") {
		if f, err := strconv.ParseFloat(input, 64); err == nil {
			return f
		}
	}

	// Try converting to int64
	if i, err := strconv.ParseInt(input, 10, 64); err == nil {
		return i
	}

	// Default to string if no other conversion succeeds
	return input
}

func expand(s string, mapping func(string) string, toType bool) interface{} {
	r := regexp.MustCompile(`\${(.*?)}`)
	re := r.FindAllStringSubmatch(s, -1)
	var ct interface{}
	for _, i := range re {
		if len(i) == 2 { //nolint:mnd
			m := mapping(i[1])
			if toType {
				ct = convertToType(m)
				return ct
			}
			s = strings.ReplaceAll(s, i[0], m)
		}
	}
	return s
}
