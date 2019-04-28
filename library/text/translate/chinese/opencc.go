package chinese

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-common/library/log"
)

var (
	defaultConversion = "s2twp"
)

// Group holds a sequence of dicts
type Group struct {
	Files []string
	Dicts []*dict
}

func (g *Group) String() string {
	return fmt.Sprintf("%+v", g.Files)
}

// OpenCC contains the converter
type openCC struct {
	Conversion  string
	Description string
	DictGroup   []*Group
}

var conversions = map[string]*openCC{
	"s2twp": {Conversion: s2twp},
	// "hk2s": {Conversion: hk2s}, "s2hk": {Conversion: s2hk}, "s2t": {Conversion: s2t},
	// "s2tw": {Conversion: s2tw}, "t2hk": {Conversion: t2hk},
	// "t2s": {Conversion: t2s}, "t2tw": {Conversion: t2tw},
	// "tw2s": {Conversion: tw2s}, "tw2sp": {Conversion: tw2sp},
}

// Init construct an instance of OpenCC.
func Init() {
	for k, v := range conversions {
		if err := v.dict(k); err != nil {
			panic(err)
		}
	}
}

// Converts .
func Converts(ctx context.Context, in ...string) (out map[string]string) {
	var err error
	out = make(map[string]string, len(in))
	for _, v := range in {
		if out[v], err = convert(v, defaultConversion); err != nil {
			log.Error("convert(%s),err:%+v", in, err)
			out[v] = v
		}
	}
	return
}

// Convert string from Simplified Chinese to Traditional Chinese .
func Convert(ctx context.Context, in string) (out string) {
	var err error
	if out, err = convert(in, defaultConversion); err != nil {
		log.Error("convert(%s),err:%+v", in, err)
	}
	return
}

func (cc *openCC) dict(conversion string) error {
	var m interface{}
	json.Unmarshal([]byte(cc.Conversion), &m)
	config := m.(map[string]interface{})
	cc.Description = config["name"].(string)
	dictChain, ok := config["conversion_chain"].([]interface{})
	if !ok {
		return fmt.Errorf("format %+v not correct", config)
	}
	for _, v := range dictChain {
		d, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("should be map inside conversion_chain")
		}
		dictMap, ok := d["dict"]
		if !ok {
			return fmt.Errorf("should have dict inside conversion_chain")
		}
		if dict, ok := dictMap.(map[string]interface{}); ok {
			group, err := cc.group(dict)
			if err != nil {
				return err
			}
			cc.DictGroup = append(cc.DictGroup, group)
		}
	}
	return nil
}

func (cc *openCC) group(d map[string]interface{}) (*Group, error) {
	typ, ok := d["type"].(string)
	if !ok {
		return nil, fmt.Errorf("type should be string")
	}
	res := &Group{}
	switch typ {
	case "group":
		dicts, ok := d["dicts"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("dicts field invalid")
		}
		for _, dict := range dicts {
			d, ok := dict.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("dicts items invalid")
			}
			group, err := cc.group(d)
			if err != nil {
				return nil, err
			}
			res.Files = append(res.Files, group.Files...)
			res.Dicts = append(res.Dicts, group.Dicts...)
		}
	case "txt":
		file, ok := d["file"]
		if !ok {
			return nil, fmt.Errorf("no file field found")
		}
		daDict, err := buildFromFile(file.(string))
		if err != nil {
			return nil, err
		}
		res.Files = append(res.Files, file.(string))
		res.Dicts = append(res.Dicts, daDict)
	default:
		return nil, fmt.Errorf("type should be txt or group")
	}
	return res, nil
}

// convert string from Simplified Chinese to Traditional Chinese or vice versa
func convert(in, conversion string) (string, error) {
	if conversion == "" {
		conversion = defaultConversion
	}
	for _, group := range conversions[conversion].DictGroup {
		r := []rune(in)
		var tokens []string
		for i := 0; i < len(r); {
			s := r[i:]
			var token string
			max := 0
			for _, dict := range group.Dicts {
				ret, err := dict.prefixMatch(string(s))
				if err != nil {
					return "", err
				}
				if len(ret) > 0 {
					o := ""
					for k, v := range ret {
						if len(k) > max {
							max = len(k)
							token = v[0]
							o = k
						}
					}
					i += len([]rune(o))
					break
				}
			}
			if max == 0 { //no match
				token = string(r[i])
				i++
			}
			tokens = append(tokens, token)
		}
		in = strings.Join(tokens, "")
	}
	return in, nil
}
