package config

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	resolveTable = []byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x53, 0x0, 0x53, 0x2e, 0x0, 0x44,
		0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x44, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x4d, 0x4d, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x4d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x4d, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x0, 0x0, 0x0, 0x0, 0x4d, 0x0, 0x0, 0x0, 0x0,
		0x4d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	}
	resolveMap = map[string]interface{}{
		"+.INF": math.Inf(1), "+.Inf": math.Inf(1), "+.inf": math.Inf(1), "-.INF": math.Inf(-1), "-.Inf": math.Inf(-1), "-.inf": math.Inf(-1),
		".INF": math.Inf(1), ".Inf": math.Inf(1), ".NAN": math.NaN(), ".NaN": math.NaN(), ".inf": math.Inf(1), ".nan": math.NaN(), "<<": "<<",
		"FALSE": false, "False": false, "NULL": nil, "Null": nil, "TRUE": true, "True": true, "false": false, "null": nil, "true": true, "~": nil,
	}
	floatregexp = regexp.MustCompile(`^[-+]?(\.[0-9]+|[0-9]+(\.[0-9]*)?)([eE][-+]?[0-9]+)?$`)
)

// Decoder is config decoder.
type Decoder func(*KeyValue, map[string]interface{}) error

// Resolver resolve placeholder in config.
type Resolver func(map[string]interface{}) error

// Option is config option.
type Option func(*options)

type options struct {
	sources  []Source
	decoder  Decoder
	resolver Resolver
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

// WithResolver with config resolver.
func WithResolver(r Resolver) Option {
	return func(o *options) {
		o.resolver = r
	}
}

// WithLogger with config logger.
// Deprecated: use global logger instead.
func WithLogger(l log.Logger) Option {
	return func(o *options) {}
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

// defaultResolver resolve placeholder in map value,
// placeholder format in ${key:default}.
func defaultResolver(input map[string]interface{}) error {
	mapper := func(name string) string {
		args := strings.SplitN(strings.TrimSpace(name), ":", 2) //nolint:gomnd
		if v, has := readValue(input, args[0]); has {
			s, _ := v.String()
			return s
		} else if len(args) > 1 { // default value
			return args[1]
		}
		return ""
	}

	var resolve func(map[string]interface{}) error
	resolve = func(sub map[string]interface{}) error {
		for k, v := range sub {
			switch vt := v.(type) {
			case string:
				sub[k] = expand(vt, mapper)
			case map[string]interface{}:
				if err := resolve(vt); err != nil {
					return err
				}
			case []interface{}:
				for i, iface := range vt {
					switch it := iface.(type) {
					case string:
						vt[i] = expand(it, mapper)
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

func expand(s string, mapper func(string) string) interface{} {
	r := regexp.MustCompile(`\${(.*?)}`)
	re := r.FindAllStringSubmatch(s, -1)
	for _, i := range re {
		if len(i) == 2 { //nolint:gomnd
			s = strings.ReplaceAll(s, i[0], mapper(i[1]))
		}
	}

	return parser(s)
}

func parser(in string) interface{} {
	if in != "" {
		// The prefix is enough of a hint about what it might be.
		hint := resolveTable[in[0]]
		if hint != 0 {
			// checked the map
			if v, ok := resolveMap[in]; ok {
				return v
			}

			switch hint {
			case 'M':
				// We've already checked the map above.
			case '.':
				// Not in the map, so maybe a normal float.
				floatv, err := strconv.ParseFloat(in, 64)
				if err == nil {
					return floatv
				}

			case 'D', 'S':
				// Int, float

				plain := strings.Replace(in, "_", "", -1)
				intv, err := strconv.ParseInt(plain, 0, 64)
				if err == nil {
					if intv == int64(int(intv)) {
						return int(intv)
					} else { //nolint:revive
						return intv
					}
				}
				uintv, err := strconv.ParseUint(plain, 0, 64)
				if err == nil {
					return uintv
				}
				if floatregexp.MatchString(plain) {
					floatv, err := strconv.ParseFloat(plain, 64)
					if err == nil {
						return floatv
					}
				}
				if strings.HasPrefix(plain, "0b") {
					intv, err := strconv.ParseInt(plain[2:], 2, 64)
					if err == nil {
						if intv == int64(int(intv)) {
							return int(intv)
						} else { //nolint:revive
							return intv
						}
					}
					uintv, err := strconv.ParseUint(plain[2:], 2, 64)
					if err == nil {
						return uintv
					}
				} else if strings.HasPrefix(plain, "-0b") {
					intv, err := strconv.ParseInt("-"+plain[3:], 2, 64)
					if err == nil {
						if true || intv == int64(int(intv)) {
							return int(intv)
						} else { //nolint:revive
							return intv
						}
					}
				}
				if strings.HasPrefix(plain, "0o") {
					intv, err := strconv.ParseInt(plain[2:], 8, 64)
					if err == nil {
						if intv == int64(int(intv)) {
							return int(intv)
						} else { //nolint:revive
							return intv
						}
					}
					uintv, err := strconv.ParseUint(plain[2:], 8, 64)
					if err == nil {
						return uintv
					}
				} else if strings.HasPrefix(plain, "-0o") {
					intv, err := strconv.ParseInt("-"+plain[3:], 8, 64)
					if err == nil {
						if true || intv == int64(int(intv)) {
							return int(intv)
						} else { //nolint:revive
							return intv
						}
					}
				}
			default:
				panic("internal error: missing handler for resolver table: " + string(rune(hint)) + " (with " + in + ")")
			}
		}
	}

	return in
}
