package generator

import (
	"fmt"
	"strings"
)

type ParamsBase struct {
	ImportPrefix string            // String to prefix to imported package file names.
	ImportMap    map[string]string // Mapping from .proto file name to import path.
	//Tpl          bool              // generate service implementation template
	ExplicitHTTP bool // Only generate for method that add http option
}

type GeneratorParamsInterface interface {
	GetBase() *ParamsBase
	SetParam(key string, value string) error
}

type BasicParam struct{ ParamsBase }

func (b *BasicParam) GetBase() *ParamsBase {
	return &b.ParamsBase
}
func (b *BasicParam) SetParam(key string, value string) error {
	return nil
}

func ParseGeneratorParams(parameter string, result GeneratorParamsInterface) error {
	ps := make(map[string]string)
	for _, p := range strings.Split(parameter, ",") {
		if p == "" {
			continue
		}
		i := strings.Index(p, "=")
		if i < 0 {
			return fmt.Errorf("invalid parameter %q: expected format of parameter to be k=v", p)
		}
		k := p[0:i]
		v := p[i+1:]
		if v == "" {
			return fmt.Errorf("invalid parameter %q: expected format of parameter to be k=v", k)
		}
		ps[k] = v
	}

	if result.GetBase().ImportMap == nil {
		result.GetBase().ImportMap = map[string]string{}
	}
	for k, v := range ps {
		switch {
		case k == "explicit_http":
			if v == "true" || v == "1" {
				result.GetBase().ExplicitHTTP = true
			}
		case k == "import_prefix":
			result.GetBase().ImportPrefix = v
			// Support import map 'M' prefix per https://github.com/golang/protobuf/blob/6fb5325/protoc-gen-go/generator/generator.go#L497.
		case len(k) > 0 && k[0] == 'M':
			result.GetBase().ImportMap[k[1:]] = v // 1 is the length of 'M'.
		case len(k) > 0 && strings.HasPrefix(k, "go_import_mapping@"):
			result.GetBase().ImportMap[k[18:]] = v // 18 is the length of 'go_import_mapping@'.
		default:
			e := result.SetParam(k, v)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
