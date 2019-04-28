package main

import (
	"bytes"
	"flag"
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"go-common/app/tool/cache/common"
)

var (
	encode    = flag.String("encode", "", "encode type: json/pb/raw/gob/gzip")
	mcType    = flag.String("type", "", "type: get/set/del/replace/only_add")
	key       = flag.String("key", "", "key name method")
	expire    = flag.String("expire", "", "expire time code")
	batchSize = flag.Int("batch", 0, "batch size")
	batchErr  = flag.String("batch_err", "break", "batch err to contine or break")
	maxGroup  = flag.Int("max_group", 0, "max group size")

	mcValidTypes   = []string{"set", "replace", "del", "get", "only_add"}
	mcValidPrefix  = []string{"set", "replace", "del", "get", "cache", "add"}
	optionNamesMap = map[string]bool{"batch": true, "max_group": true, "encode": true, "type": true, "key": true, "expire": true, "batch_err": true}
	simpleTypes    = []string{"int", "int8", "int16", "int32", "int64", "float32", "float64", "uint", "uint8", "uint16", "uint32", "uint64", "bool", "string", "[]byte"}
	lenTypes       = []string{"[]", "map"}
)

const (
	_interfaceName = "_mc"
	_multiTpl      = 1
	_singleTpl     = 2
	_noneTpl       = 3
	_typeGet       = "get"
	_typeSet       = "set"
	_typeDel       = "del"
	_typeReplace   = "replace"
	_typeAdd       = "only_add"
)

func resetFlag() {
	*encode = ""
	*mcType = ""
	*batchSize = 0
	*maxGroup = 0
	*batchErr = "break"
}

// options options
type options struct {
	name        string
	keyType     string
	ValueType   string
	template    int
	SimpleValue bool
	// int float 类型
	GetSimpleValue bool
	// string, []byte类型
	GetDirectValue     bool
	ConvertValue2Bytes string
	ConvertBytes2Value string
	GoValue            bool
	ImportPackage      string
	importPackages     []string
	Args               string
	PkgName            string
	ExtraArgsType      string
	ExtraArgs          string
	MCType             string
	KeyMethod          string
	ExpireCode         string
	Encode             string
	UseMemcached       bool
	InitValue          bool
	OriginValueType    string
	UseStrConv         bool
	Comment            string
	GroupSize          int
	MaxGroup           int
	EnableBatch        bool
	BatchErrBreak      bool
	LenType            bool
	PointType          bool
}

func parse(s *common.Source) (opts []*options) {
	f := s.F
	fset := s.Fset
	src := s.Src
	c := f.Scope.Lookup(_interfaceName)
	if (c == nil) || (c.Kind != ast.Typ) {
		log.Fatalln("无法找到缓存声明")
	}
	lines := strings.Split(src, "\n")
	lists := c.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
	for _, list := range lists {
		opt := options{Args: s.GetDef(_interfaceName), UseMemcached: true, importPackages: s.Packages(list)}
		opt.name = list.Names[0].Name
		opt.KeyMethod = "key" + opt.name
		opt.ExpireCode = "d.mc" + opt.name + "Expire"
		// get comment
		line := fset.Position(list.Pos()).Line - 3
		if len(lines)-1 >= line {
			comment := lines[line]
			opt.Comment = common.RegexpReplace(`\s+//(?P<name>.+)`, comment, "$name")
			opt.Comment = strings.TrimSpace(opt.Comment)
		}
		// get options
		line = fset.Position(list.Pos()).Line - 2
		comment := lines[line]
		os.Args = []string{os.Args[0]}
		if regexp.MustCompile(`\s+//\s*mc:.+`).Match([]byte(comment)) {
			args := strings.Split(common.RegexpReplace(`//\s*mc:(?P<arg>.+)`, comment, "$arg"), " ")
			for _, arg := range args {
				arg = strings.TrimSpace(arg)
				if arg != "" {
					// validate option name
					argName := common.RegexpReplace(`-(?P<name>[\w_-]+)=.+`, arg, "$name")
					if !optionNamesMap[argName] {
						log.Fatalf("选项:%s 不存在 请检查拼写\n", argName)
					}
					os.Args = append(os.Args, arg)
				}
			}
		}
		resetFlag()
		flag.Parse()
		if *mcType != "" {
			opt.MCType = *mcType
		}
		if *key != "" {
			opt.KeyMethod = *key
		}
		if *expire != "" {
			opt.ExpireCode = *expire
		}
		opt.EnableBatch = (*batchSize != 0) && (*maxGroup != 0)
		opt.BatchErrBreak = *batchErr == "break"
		opt.GroupSize = *batchSize
		opt.MaxGroup = *maxGroup
		// get type from prefix
		if opt.MCType == "" {
			for _, t := range mcValidPrefix {
				if strings.HasPrefix(strings.ToLower(opt.name), t) {
					if t == "add" {
						t = _typeSet
					}
					opt.MCType = t
					break
				}
			}
			if opt.MCType == "" {
				log.Fatalln(opt.name + "请指定方法类型(type=get/set/del...)")
			}
		}
		if opt.MCType == "cache" {
			opt.MCType = _typeGet
		}
		params := list.Type.(*ast.FuncType).Params.List
		if len(params) == 0 {
			log.Fatalln(opt.name + "参数不足")
		}
		if s.ExprString(params[0].Type) != "context.Context" {
			log.Fatalln(opt.name + "第一个参数必须为context")
		}
		for _, param := range params {
			if len(param.Names) > 1 {
				log.Fatalln(opt.name + "不支持省略类型")
			}
		}
		// get template
		if len(params) == 1 {
			opt.template = _noneTpl
		} else if (len(params) == 2) && (opt.MCType == _typeSet || opt.MCType == _typeAdd || opt.MCType == _typeReplace) {
			if _, ok := params[1].Type.(*ast.MapType); ok {
				opt.template = _multiTpl
			} else {
				opt.template = _noneTpl
			}
		} else {
			if _, ok := params[1].Type.(*ast.ArrayType); ok {
				opt.template = _multiTpl
			} else {
				opt.template = _singleTpl
			}
		}
		// extra args
		if len(params) > 2 {
			args := []string{""}
			allArgs := []string{""}
			var pos = 2
			if (opt.MCType == _typeAdd) || (opt.MCType == _typeSet) || (opt.MCType == _typeReplace) {
				pos = 3
			}
			for _, pa := range params[pos:] {
				paType := s.ExprString(pa.Type)
				if len(pa.Names) == 0 {
					args = append(args, paType)
					allArgs = append(allArgs, paType)
					continue
				}
				var names []string
				for _, name := range pa.Names {
					names = append(names, name.Name)
				}
				allArgs = append(allArgs, strings.Join(names, ",")+" "+paType)
				args = append(args, strings.Join(names, ","))
			}
			if len(args) > 1 {
				opt.ExtraArgs = strings.Join(args, ",")
				opt.ExtraArgsType = strings.Join(allArgs, ",")
			}
		}
		// get k v from results
		results := list.Type.(*ast.FuncType).Results.List
		if s.ExprString(results[len(results)-1].Type) != "error" {
			log.Fatalln("最后返回值参数需为error")
		}
		for _, res := range results {
			if len(res.Names) > 1 {
				log.Fatalln(opt.name + "返回值不支持省略类型")
			}
		}
		if opt.MCType == _typeGet {
			if len(results) != 2 {
				log.Fatalln("参数个数不对")
			}
		}
		// get key type and value type
		if (opt.MCType == _typeAdd) || (opt.MCType == _typeSet) || (opt.MCType == _typeReplace) {
			if opt.template == _multiTpl {
				p, ok := params[1].Type.(*ast.MapType)
				if !ok {
					log.Fatalf("%s: 参数类型错误 批量设置数据时类型需为map类型\n", opt.name)
				}
				opt.keyType = s.ExprString(p.Key)
				opt.ValueType = s.ExprString(p.Value)
			} else if opt.template == _singleTpl {
				opt.keyType = s.ExprString(params[1].Type)
				opt.ValueType = s.ExprString(params[2].Type)
			} else {
				opt.ValueType = s.ExprString(params[1].Type)
			}
		}
		if opt.MCType == _typeGet {
			if opt.template == _multiTpl {
				if p, ok := results[0].Type.(*ast.MapType); ok {
					opt.keyType = s.ExprString(p.Key)
					opt.ValueType = s.ExprString(p.Value)
				} else {
					log.Fatalf("%s: 返回值类型错误 批量获取数据时返回值需为map类型\n", opt.name)
				}
			} else if opt.template == _singleTpl {
				opt.keyType = s.ExprString(params[1].Type)
				opt.ValueType = s.ExprString(results[0].Type)
			} else {
				opt.ValueType = s.ExprString(results[0].Type)
			}
		}
		if opt.MCType == _typeDel {
			if opt.template == _multiTpl {
				p, ok := params[1].Type.(*ast.ArrayType)
				if !ok {
					log.Fatalf("%s: 类型错误 参数需为[]类型\n", opt.name)
				}
				opt.keyType = s.ExprString(p.Elt)
			} else if opt.template == _singleTpl {
				opt.keyType = s.ExprString(params[1].Type)
			}
		}
		for _, t := range simpleTypes {
			if t == opt.ValueType {
				opt.SimpleValue = true
				opt.GetSimpleValue = true
				opt.ConvertValue2Bytes = convertValue2Bytes(t)
				opt.ConvertBytes2Value = convertBytes2Value(t)
				break
			}
		}
		if opt.ValueType == "string" {
			opt.LenType = true
		} else {
			for _, t := range lenTypes {
				if strings.HasPrefix(opt.ValueType, t) {
					opt.LenType = true
					break
				}
			}
		}
		if opt.SimpleValue && (opt.ValueType == "[]byte" || opt.ValueType == "string") {
			opt.GetSimpleValue = false
			opt.GetDirectValue = true
		}
		if opt.MCType == _typeGet && opt.template == _multiTpl {
			opt.UseMemcached = false
		}
		if strings.HasPrefix(opt.ValueType, "*") {
			opt.InitValue = true
			opt.PointType = true
			opt.OriginValueType = strings.Replace(opt.ValueType, "*", "", 1)
		} else {
			opt.OriginValueType = opt.ValueType
		}
		if *encode != "" {
			var flags []string
			for _, f := range strings.Split(*encode, "|") {
				switch f {
				case "gob":
					flags = append(flags, "memcache.FlagGOB")
				case "json":
					flags = append(flags, "memcache.FlagJSON")
				case "raw":
					flags = append(flags, "memcache.FlagRAW")
				case "pb":
					flags = append(flags, "memcache.FlagProtobuf")
				case "gzip":
					flags = append(flags, "memcache.FlagGzip")
				default:
					log.Fatalf("%s: encode类型无效\n", opt.name)
				}
			}
			opt.Encode = strings.Join(flags, " | ")
		} else {
			if opt.SimpleValue {
				opt.Encode = "memcache.FlagRAW"
			} else {
				opt.Encode = "memcache.FlagJSON"
			}
		}
		opt.Check()
		opts = append(opts, &opt)
	}
	return
}

func (option *options) Check() {
	var valid bool
	for _, x := range mcValidTypes {
		if x == option.MCType {
			valid = true
			break
		}
	}
	if !valid {
		log.Fatalf("%s: 类型错误 不支持%s类型\n", option.name, option.MCType)
	}
	if (option.MCType != _typeDel) && !option.SimpleValue && !strings.Contains(option.ValueType, "*") && !strings.Contains(option.ValueType, "[]") && !strings.Contains(option.ValueType, "map") {
		log.Fatalf("%s: 值类型只能为基本类型/slice/map/指针类型\n", option.name)
	}
}

func genHeader(opts []*options) (src string) {
	option := options{PkgName: os.Getenv("GOPACKAGE"), UseMemcached: false}
	var packages []string
	packagesMap := map[string]bool{`"context"`: true}
	for _, opt := range opts {
		if len(opt.importPackages) > 0 {
			for _, pkg := range opt.importPackages {
				if !packagesMap[pkg] {
					packages = append(packages, pkg)
					packagesMap[pkg] = true
				}
			}
		}
		if opt.Args != "" {
			option.Args = opt.Args
		}
		if opt.UseMemcached {
			option.UseMemcached = true
		}
		if opt.SimpleValue && !opt.GetDirectValue {
			option.UseStrConv = true
		}
		if opt.EnableBatch {
			option.EnableBatch = true
		}
	}
	option.ImportPackage = strings.Join(packages, "\n")
	src = _headerTemplate
	t := template.Must(template.New("header").Parse(src))
	var buffer bytes.Buffer
	err := t.Execute(&buffer, option)
	if err != nil {
		log.Fatalf("execute template: %s", err)
	}
	// Format the output.
	src = strings.Replace(buffer.String(), "\t", "", -1)
	src = regexp.MustCompile("\n+").ReplaceAllString(src, "\n")
	src = strings.Replace(src, "NEWLINE", "", -1)
	src = strings.Replace(src, "ARGS", option.Args, -1)
	return
}

func genBody(opts []*options) (res string) {
	for _, option := range opts {
		var src string
		if option.template == _multiTpl {
			switch option.MCType {
			case _typeGet:
				src = _multiGetTemplate
			case _typeSet:
				src = _multiSetTemplate
			case _typeReplace:
				src = _multiReplaceTemplate
			case _typeDel:
				src = _multiDelTemplate
			case _typeAdd:
				src = _multiAddTemplate
			}
		} else if option.template == _singleTpl {
			switch option.MCType {
			case _typeGet:
				src = _singleGetTemplate
			case _typeSet:
				src = _singleSetTemplate
			case _typeReplace:
				src = _singleReplaceTemplate
			case _typeDel:
				src = _singleDelTemplate
			case _typeAdd:
				src = _singleAddTemplate
			}
		} else {
			switch option.MCType {
			case _typeGet:
				src = _noneGetTemplate
			case _typeSet:
				src = _noneSetTemplate
			case _typeReplace:
				src = _noneReplaceTemplate
			case _typeDel:
				src = _noneDelTemplate
			case _typeAdd:
				src = _noneAddTemplate
			}
		}
		src = strings.Replace(src, "KEY", option.keyType, -1)
		src = strings.Replace(src, "NAME", option.name, -1)
		src = strings.Replace(src, "VALUE", option.ValueType, -1)
		src = strings.Replace(src, "GROUPSIZE", strconv.Itoa(option.GroupSize), -1)
		src = strings.Replace(src, "MAXGROUP", strconv.Itoa(option.MaxGroup), -1)
		t := template.Must(template.New("cache").Parse(src))
		var buffer bytes.Buffer
		err := t.Execute(&buffer, option)
		if err != nil {
			log.Fatalf("execute template: %s", err)
		}
		// Format the output.
		src = strings.Replace(buffer.String(), "\t", "", -1)
		src = regexp.MustCompile("\n+").ReplaceAllString(src, "\n")
		res = res + "\n" + src
	}
	return
}

func main() {
	log.SetFlags(0)
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("程序解析失败, err: %+v  请企业微信联系 @wangxu01", err)
		}
	}()
	options := parse(common.NewSource(common.SourceText()))
	header := genHeader(options)
	body := genBody(options)
	code := common.FormatCode(header + "\n" + body)
	// Write to file.
	dir := filepath.Dir(".")
	outputName := filepath.Join(dir, "mc.cache.go")
	err := ioutil.WriteFile(outputName, []byte(code), 0644)
	if err != nil {
		log.Fatalf("写入文件失败: %s", err)
	}
	log.Println("mc.cache.go: 生成成功")
}

func convertValue2Bytes(t string) string {
	switch t {
	case "int", "int8", "int16", "int32", "int64":
		return "[]byte(strconv.FormatInt(int64(val), 10))"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "[]byte(strconv.FormatUInt(val, 10))"
	case "bool":
		return "[]byte(strconv.FormatBool(val))"
	case "float32":
		return "[]byte(strconv.FormatFloat(val, 'E', -1, 32))"
	case "float64":
		return "[]byte(strconv.FormatFloat(val, 'E', -1, 64))"
	case "string":
		return "[]byte(val)"
	case "[]byte":
		return "val"
	}
	return ""
}

func convertBytes2Value(t string) string {
	switch t {
	case "int", "int8", "int16", "int32", "int64":
		return "strconv.ParseInt(v, 10, 64)"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "strconv.ParseUInt(v, 10, 64)"
	case "bool":
		return "strconv.ParseBool(v)"
	case "float32":
		return "float32(strconv.ParseFloat(v, 32))"
	case "float64":
		return "strconv.ParseFloat(v, 64)"
	}
	return ""
}
