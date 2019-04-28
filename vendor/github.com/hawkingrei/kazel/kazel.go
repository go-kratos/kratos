/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"

	bzl "github.com/bazelbuild/buildtools/build"
	"github.com/golang/glog"
)

const (
	vendorPath     = "vendor/"
	automanagedTag = "automanaged"
	manualTag      = "manual"
)

var (
	root      = flag.String("root", ".", "root of go source")
	dryRun    = flag.Bool("dry-run", false, "run in dry mode")
	printDiff = flag.Bool("print-diff", false, "print diff to stdout")
	validate  = flag.Bool("validate", false, "run in dry mode and exit nonzero if any BUILD files need to be updated")
	cfgPath   = flag.String("cfg-path", ".kazelcfg.json", "path to kazel config (relative paths interpreted relative to -repo.")
	iswrote   = false
)

func main() {
	flag.Parse()
	flag.Set("alsologtostderr", "true")
	if *root == "" {
		glog.Fatalf("-root argument is required")
	}
	if *validate {
		*dryRun = true
	}
	v, err := newVendorer(*root, *cfgPath, *dryRun)
	if err != nil {
		glog.Fatalf("unable to build vendorer: %v", err)
	}
	if err = os.Chdir(v.root); err != nil {
		glog.Fatalf("cannot chdir into root %q: %v", v.root, err)
	}

	if v.cfg.ManageGoRules {
		if err = v.walkVendor(); err != nil {
			glog.Fatalf("err walking vendor: %v", err)
		}
		if err = v.walkRepo(); err != nil {
			glog.Fatalf("err walking repo: %v", err)
		}
	}
	if err = v.walkGenerated(); err != nil {
		glog.Fatalf("err walking generated: %v", err)
	}
	if _, err = v.walkSource("."); err != nil {
		glog.Fatalf("err walking source: %v", err)
	}
	written := 0
	if written, err = v.reconcileAllRules(); err != nil {
		glog.Fatalf("err reconciling rules: %v", err)
	}
	if *validate && written > 0 {
		fmt.Fprintf(os.Stderr, "\n%d BUILD files not up-to-date.\n", written)
		os.Exit(1)
	}
	if iswrote {
		fmt.Fprintf(os.Stderr, "\nPlease re-run git-add\n")
		os.Exit(1)
	}
}

// Vendorer collects context, configuration, and cache while walking the tree.
type Vendorer struct {
	ctx          *build.Context
	icache       map[icacheKey]icacheVal
	skippedPaths []*regexp.Regexp
	dryRun       bool
	root         string
	cfg          *Cfg
	newRules     map[string][]*bzl.Rule // package path -> list of rules to add or update
	managedAttrs []string
}

func newVendorer(root, cfgPath string, dryRun bool) (*Vendorer, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path: %v", err)
	}
	if !filepath.IsAbs(cfgPath) {
		cfgPath = filepath.Join(absRoot, cfgPath)
	}
	cfg, err := ReadCfg(cfgPath)
	if err != nil {
		return nil, err
	}

	v := Vendorer{
		ctx:          context(),
		dryRun:       dryRun,
		root:         absRoot,
		icache:       map[icacheKey]icacheVal{},
		cfg:          cfg,
		newRules:     make(map[string][]*bzl.Rule),
		managedAttrs: []string{"srcs", "deps", "importpath", "compilers"},
	}

	for _, sp := range cfg.SkippedPaths {
		r, err := regexp.Compile(sp)
		if err != nil {
			return nil, err
		}
		v.skippedPaths = append(v.skippedPaths, r)
	}
	for _, builtinSkip := range []string{
		"^\\.git",
		"^bazel-*",
	} {
		v.skippedPaths = append(v.skippedPaths, regexp.MustCompile(builtinSkip))
	}

	return &v, nil

}

type icacheKey struct {
	path, srcDir string
}

type icacheVal struct {
	pkg *build.Package
	err error
}

func (v *Vendorer) importPkg(path string, srcDir string) (*build.Package, error) {
	k := icacheKey{path: path, srcDir: srcDir}
	if val, ok := v.icache[k]; ok {
		return val.pkg, val.err
	}

	// cache miss
	pkg, err := v.ctx.Import(path, srcDir, build.ImportComment)
	v.icache[k] = icacheVal{pkg: pkg, err: err}
	return pkg, err
}

func writeHeaders(file *bzl.File) {
	pkgRule := bzl.Rule{
		Call: &bzl.CallExpr{
			X: &bzl.LiteralExpr{Token: "package"},
		},
	}
	pkgRule.SetAttr("default_visibility", asExpr([]string{"//visibility:public"}))

	file.Stmt = append(file.Stmt,
		[]bzl.Expr{
			pkgRule.Call,
			&bzl.CallExpr{
				X: &bzl.LiteralExpr{Token: "load"},
				List: asExpr([]string{
					"@io_bazel_rules_go//go:def.bzl",
				}).(*bzl.ListExpr).List,
			},
		}...,
	)
}

func writeRules(file *bzl.File, rules []*bzl.Rule) {
	for _, rule := range rules {
		file.Stmt = append(file.Stmt, rule.Call)
	}
}

func (v *Vendorer) resolve(ipath string) Label {
	if ipath == v.cfg.GoPrefix {
		return Label{
			tag: "go_default_library",
		}
	} else if strings.HasPrefix(ipath, v.cfg.GoPrefix) {
		return Label{
			pkg: strings.TrimPrefix(ipath, v.cfg.GoPrefix+"/"),
			tag: "go_default_library",
		}
	}
	if v.cfg.VendorMultipleBuildFiles {
		return Label{
			pkg: "vendor/" + ipath,
			tag: "go_default_library",
		}
	}
	return Label{
		pkg: "vendor",
		tag: ipath,
	}
}

func (v *Vendorer) walk(root string, f func(path, ipath string, pkg *build.Package, conffile, proto []string) error) error {
	skipVendor := true
	if root == vendorPath {
		skipVendor = false
	}
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if skipVendor && strings.HasPrefix(path, vendorPath) {
			return filepath.SkipDir
		}
		for _, r := range v.skippedPaths {
			if r.MatchString(path) {
				return filepath.SkipDir
			}
		}
		if _, err = os.Stat(filepath.Join(path, ".skip_kazel")); !os.IsNotExist(err) {
			return filepath.SkipDir
		}
		ipath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		protofiles := v.getAllProto(path)
		conffiles := v.getAllConf(path)
		pkg, err := v.importPkg(".", filepath.Join(v.root, path))
		if err != nil {
			if _, ok := err.(*build.NoGoError); err != nil && ok {
				return nil
			}
			return err
		}
		return f(path, ipath, pkg, conffiles, protofiles)
	})
}

func (v *Vendorer) walkRepo() error {
	for _, root := range v.cfg.SrcDirs {
		if err := v.walk(root, v.updatePkg); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vendorer) getAllProto(path string) []string {
	var protofiles []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		glog.Fatalf("getallproto fail to readdir")
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".proto") {
			if f.Mode() != os.ModeDir {
				protofiles = append(protofiles, f.Name())
			}
		}
	}
	return protofiles
}

func (v *Vendorer) getAllConf(path string) []string {
	var conffiles []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		glog.Fatalf("getallconf fail to readdir")
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".toml") || strings.Contains(f.Name(), ".yaml") {
			if f.Mode() != os.ModeDir {
				conffiles = append(conffiles, f.Name())
			}
		}
	}
	return conffiles
}

func (v *Vendorer) updateSinglePkg(path string) error {
	pkg, err := v.importPkg(".", "./"+path)
	if err != nil {
		if _, ok := err.(*build.NoGoError); err != nil && ok {
			return nil
		}
		return err
	}
	protofiles := v.getAllProto(path)
	conffiles := v.getAllConf(path)
	return v.updatePkg(path, "", pkg, conffiles, protofiles)
}

type ruleType int

// The RuleType* constants enumerate the bazel rules supported by this tool.
const (
	RuleTypeGoBinary ruleType = iota
	RuleTypeGoLibrary
	RuleTypeGoTest
	RuleTypeGoXTest
	RuleTypeCGoGenrule
	RuleTypeFileGroup
	RuleTypeOpenAPILibrary
	RuleTypeProtoLibrary
	RuleTypeGoProtoLibrary
)

// RuleKind converts a value of the RuleType* enum into the BUILD string.
func (rt ruleType) RuleKind() string {
	switch rt {
	case RuleTypeGoBinary:
		return "go_binary"
	case RuleTypeGoLibrary:
		return "go_library"
	case RuleTypeGoTest:
		return "go_test"
	case RuleTypeGoXTest:
		return "go_test"
	case RuleTypeCGoGenrule:
		return "cgo_genrule"
	case RuleTypeFileGroup:
		return "filegroup"
	case RuleTypeOpenAPILibrary:
		return "openapi_library"
	case RuleTypeProtoLibrary:
		return "proto_library"
	case RuleTypeGoProtoLibrary:
		return "go_proto_library"
	}
	panic("unreachable")
}

// NamerFunc is a function that returns the appropriate name for the rule for the provided RuleType.
type NamerFunc func(ruleType) string

func (v *Vendorer) updatePkg(path, _ string, pkg *build.Package, conffile, protofile []string) error {

	srcNameMap := func(srcs ...[]string) *bzl.ListExpr {
		return asExpr(merge(srcs...)).(*bzl.ListExpr)
	}
	goFileNotProto := []string{}
	for _, v := range pkg.GoFiles {
		if !strings.Contains(v, ".pb.go") {
			goFileNotProto = append(goFileNotProto, v)
		}
	}
	srcs := srcNameMap(goFileNotProto, pkg.SFiles)
	cgoSrcs := srcNameMap(pkg.CgoFiles, pkg.CFiles, pkg.CXXFiles, pkg.HFiles)
	testSrcs := srcNameMap(pkg.TestGoFiles)
	xtestSrcs := srcNameMap(pkg.XTestGoFiles)
	pf := protoFileInfo(v.cfg.GoPrefix, path, protofile)

	v.addRules(path, v.emit(path, srcs, cgoSrcs, testSrcs, xtestSrcs, pf, pkg, conffile, func(rt ruleType) string {
		switch rt {
		case RuleTypeGoBinary:
			return filepath.Base(pkg.Dir)
		case RuleTypeGoLibrary:
			return "go_default_library"
		case RuleTypeGoTest:
			return "go_default_test"
		case RuleTypeGoXTest:
			return "go_default_xtest"
		case RuleTypeCGoGenrule:
			return "cgo_codegen"
		case RuleTypeProtoLibrary:
			return pf.packageName + "_proto"
		case RuleTypeGoProtoLibrary:
			return pf.packageName + "_go_proto"
		}
		panic("unreachable")
	}))

	return nil
}

func (v *Vendorer) emit(path string, srcs, cgoSrcs, testSrcs, xtestSrcs *bzl.ListExpr, protoSrcs ProtoInfo, pkg *build.Package, conffile []string, namer NamerFunc) []*bzl.Rule {
	var goLibAttrs = make(Attrs)
	var rules []*bzl.Rule
	embedlist := []string{}
	if len(protoSrcs.src) > 0 {
		protoRuleAttrs := make(Attrs)

		protoRuleAttrs.SetList("srcs", asExpr(protoSrcs.src).(*bzl.ListExpr))
		protoRuleAttrs.SetList("deps", asExpr(protoMap(path, protoSrcs.imports)).(*bzl.ListExpr))

		rules = append(rules, newRule(RuleTypeProtoLibrary, namer, protoRuleAttrs))
		goProtoRuleAttrs := make(Attrs)
		if protoSrcs.isGogo {
			if protoSrcs.hasServices {
				goProtoRuleAttrs.SetList("compilers", asExpr([]string{"@io_bazel_rules_go//proto:gogofast_grpc"}).(*bzl.ListExpr))
			} else {
				goProtoRuleAttrs.SetList("compilers", asExpr([]string{"@io_bazel_rules_go//proto:gogofast_proto"}).(*bzl.ListExpr))
			}
		} else {
			if protoSrcs.hasServices {
				goProtoRuleAttrs.SetList("compilers", asExpr([]string{"@io_bazel_rules_go//proto:go_grpc"}).(*bzl.ListExpr))
			} else {
				goProtoRuleAttrs.SetList("compilers", asExpr([]string{"@io_bazel_rules_go//proto:go_proto"}).(*bzl.ListExpr))
			}
		}

		protovalue := ":" + protoSrcs.packageName + "_proto"
		goProtoRuleAttrs.Set("proto", asExpr(protovalue))
		goProtoRuleAttrs.Set("importpath", asExpr(protoSrcs.importPath))
		goProtoRuleAttrs.SetList("deps", asExpr(goProtoMap(path, protoSrcs.imports)).(*bzl.ListExpr))
		rules = append(rules, newRule(RuleTypeGoProtoLibrary, namer, goProtoRuleAttrs))

		embedlist = append(embedlist, protoSrcs.packageName+"_go_proto")
	}

	deps := v.extractDeps(depMapping(pkg.Imports))
	if len(srcs.List) >= 0 {
		if len(cgoSrcs.List) != 0 {
			goLibAttrs.SetList("srcs", &bzl.ListExpr{List: addExpr(srcs.List, cgoSrcs.List)})
			goLibAttrs.SetList("clinkopts", asExpr([]string{"-lz", "-lm", "-lpthread", "-ldl"}).(*bzl.ListExpr))
			goLibAttrs.Set("cgo", &bzl.LiteralExpr{Token: "True"})
		} else {
			goLibAttrs.Set("srcs", srcs)
		}
		if strings.Contains(path, "vendor") {
			goLibAttrs.Set("importpath", asExpr(strings.Replace(path, "vendor/", "", -1)))
		} else {
			goLibAttrs.Set("importpath", asExpr(filepath.Join(v.cfg.GoPrefix, path)))
		}
		goLibAttrs.SetList("visibility", asExpr([]string{"//visibility:public"}).(*bzl.ListExpr))

	} else if len(cgoSrcs.List) == 0 {
		return nil
	}
	if len(conffile) > 0 {
		goLibAttrs.SetList("data", asExpr(conffile).(*bzl.ListExpr))
	}
	if len(deps.List) > 0 {
		goLibAttrs.SetList("deps", deps)
	}

	if pkg.IsCommand() {
		rules = append(rules, newRule(RuleTypeGoBinary, namer, map[string]bzl.Expr{
			"embed": asExpr([]string{":" + namer(RuleTypeGoLibrary)}),
		}))
	}

	addGoDefaultLibrary := len(cgoSrcs.List) > 0 || len(srcs.List) > 0 || len(protoSrcs.src) == 1 || len(conffile) > 0

	if len(testSrcs.List) != 0 {
		testRuleAttrs := make(Attrs)

		testRuleAttrs.SetList("srcs", testSrcs)
		testRuleAttrs.SetList("deps", v.extractDeps(depMapping(pkg.TestImports)))
		//testRuleAttrs.Set("rundir", asExpr("."))
		//testRuleAttrs.Set("importmap", asExpr(filepath.Join(v.cfg.GoPrefix, path)))
		//testRuleAttrs.Set("importpath", asExpr(filepath.Join(v.cfg.GoPrefix, path)))
		if addGoDefaultLibrary {
			testRuleAttrs.SetList("embed", asExpr([]string{":" + namer(RuleTypeGoLibrary)}).(*bzl.ListExpr))

		}
		rules = append(rules, newRule(RuleTypeGoTest, namer, testRuleAttrs))
	}
	if len(embedlist) > 0 {
		goLibAttrs.SetList("embed", asExpr(embedlist).(*bzl.ListExpr))
	}
	if addGoDefaultLibrary || len(embedlist) > 0 {
		rules = append(rules, newRule(RuleTypeGoLibrary, namer, goLibAttrs))
	}

	if len(xtestSrcs.List) != 0 {
		xtestRuleAttrs := make(Attrs)

		xtestRuleAttrs.SetList("srcs", xtestSrcs)
		xtestRuleAttrs.SetList("deps", v.extractDeps(pkg.XTestImports))

		rules = append(rules, newRule(RuleTypeGoXTest, namer, xtestRuleAttrs))
	}

	return rules
}

func (v *Vendorer) addRules(pkgPath string, rules []*bzl.Rule) {
	cleanPath := filepath.Clean(pkgPath)
	v.newRules[cleanPath] = append(v.newRules[cleanPath], rules...)
}

func (v *Vendorer) walkVendor() error {
	var rules []*bzl.Rule
	updateFunc := func(path, ipath string, pkg *build.Package, conffile, proto []string) error {
		srcNameMap := func(srcs ...[]string) *bzl.ListExpr {
			return asExpr(
				apply(
					merge(srcs...),
					mapper(func(s string) string {
						return strings.TrimPrefix(filepath.Join(path, s), "vendor/")
					}),
				),
			).(*bzl.ListExpr)
		}

		srcs := srcNameMap(pkg.GoFiles, pkg.SFiles)
		cgoSrcs := srcNameMap(pkg.CgoFiles, pkg.CFiles, pkg.CXXFiles, pkg.HFiles)
		testSrcs := srcNameMap(pkg.TestGoFiles)
		xtestSrcs := srcNameMap(pkg.XTestGoFiles)
		pf := protoFileInfo(v.cfg.GoPrefix, path, proto)
		tagBase := v.resolve(ipath).tag

		rules = append(rules, v.emit(path, srcs, cgoSrcs, testSrcs, xtestSrcs, pf, pkg, []string{}, func(rt ruleType) string {
			switch rt {
			case RuleTypeGoBinary:
				return tagBase + "_bin"
			case RuleTypeGoLibrary:
				return tagBase
			case RuleTypeGoTest:
				return tagBase + "_test"
			case RuleTypeGoXTest:
				return tagBase + "_xtest"
			case RuleTypeCGoGenrule:
				return tagBase + "_cgo"
			case RuleTypeProtoLibrary:
				return pf.packageName + "_proto"
			case RuleTypeGoProtoLibrary:
				return pf.packageName + "_go_proto"
			}
			panic("unreachable")
		})...)

		return nil
	}
	if v.cfg.VendorMultipleBuildFiles {
		updateFunc = v.updatePkg
	}
	if err := v.walk(vendorPath, updateFunc); err != nil {
		return err
	}
	v.addRules(vendorPath, rules)

	return nil
}

func (v *Vendorer) extractDeps(deps []string) *bzl.ListExpr {
	return asExpr(
		depMapping(apply(
			merge(deps),
			filterer(func(s string) bool {
				pkg, err := v.importPkg(s, v.root)
				if err != nil {
					if strings.Contains(err.Error(), `cannot find package "C"`) ||
						// added in go1.7
						strings.Contains(err.Error(), `cannot find package "context"`) ||
						strings.Contains(err.Error(), `cannot find package "net/http/httptrace"`) {
						return false
					}
					fmt.Fprintf(os.Stderr, "extract err: %v\n", err)
					return false
				}
				if pkg.Goroot {
					return false
				}
				return true
			}),
			mapper(func(s string) string {
				return v.resolve(s).String()
			}),
		)),
	).(*bzl.ListExpr)
}

func (v *Vendorer) reconcileAllRules() (int, error) {
	var paths []string
	for path := range v.newRules {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	written := 0
	for _, path := range paths {
		w, err := ReconcileRules(path, v.newRules[path], v.managedAttrs, v.dryRun, v.cfg.ManageGoRules)
		if w {
			written++
		}
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

// Attrs collects the attributes for a rule.
type Attrs map[string]bzl.Expr

// Set sets the named attribute to the provided bazel expression.
func (a Attrs) Set(name string, expr bzl.Expr) {
	a[name] = expr
}

// SetList sets the named attribute to the provided bazel expression list.
func (a Attrs) SetList(name string, expr *bzl.ListExpr) {
	if len(expr.List) == 0 {
		return
	}
	a[name] = expr
}

// Label defines a bazel label.
type Label struct {
	pkg, tag string
}

func (l Label) String() string {
	return fmt.Sprintf("//%v:%v", l.pkg, l.tag)
}

func asExpr(e interface{}) bzl.Expr {
	rv := reflect.ValueOf(e)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &bzl.LiteralExpr{Token: fmt.Sprintf("%d", e)}
	case reflect.Float32, reflect.Float64:
		return &bzl.LiteralExpr{Token: fmt.Sprintf("%f", e)}
	case reflect.String:
		return &bzl.StringExpr{Value: e.(string)}
	case reflect.Slice, reflect.Array:
		var list []bzl.Expr
		for i := 0; i < rv.Len(); i++ {
			list = append(list, asExpr(rv.Index(i).Interface()))
		}
		return &bzl.ListExpr{List: list}
	default:
		glog.Fatalf("Uh oh")
		return nil
	}
}

type sed func(s []string) []string

func mapString(in []string, f func(string) string) []string {
	var out []string
	for _, s := range in {
		out = append(out, f(s))
	}
	return out
}

func mapper(f func(string) string) sed {
	return func(in []string) []string {
		return mapString(in, f)
	}
}

func filterString(in []string, f func(string) bool) []string {
	var out []string
	for _, s := range in {
		if f(s) {
			out = append(out, s)
		}
	}
	return out
}

func filterer(f func(string) bool) sed {
	return func(in []string) []string {
		return filterString(in, f)
	}
}

func apply(stream []string, seds ...sed) []string {
	for _, sed := range seds {
		stream = sed(stream)
	}
	return stream
}

func merge(streams ...[]string) []string {
	var out []string
	for _, stream := range streams {
		out = append(out, stream...)
	}
	return out
}

func newRule(rt ruleType, namer NamerFunc, attrs map[string]bzl.Expr) *bzl.Rule {
	rule := &bzl.Rule{
		Call: &bzl.CallExpr{
			X: &bzl.LiteralExpr{Token: rt.RuleKind()},
		},
	}
	rule.SetAttr("name", asExpr(namer(rt)))
	for k, v := range attrs {
		rule.SetAttr(k, v)
	}
	rule.SetAttr("tags", asExpr([]string{automanagedTag}))
	return rule
}

// findBuildFile determines the name of a preexisting BUILD file, returning
// a default if no such file exists.
func findBuildFile(pkgPath string) (bool, string) {
	options := []string{"BUILD.bazel", "BUILD"}
	for _, b := range options {
		path := filepath.Join(pkgPath, b)
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return true, path
		}
	}
	return false, filepath.Join(pkgPath, "BUILD")
}

// ReconcileRules reconciles, simplifies, and writes the rules for the specified package, adding
// additional dependency rules as needed.
func ReconcileRules(pkgPath string, rules []*bzl.Rule, managedAttrs []string, dryRun bool, manageGoRules bool) (bool, error) {
	goProtoLibrary := []string{}
	_, path := findBuildFile(pkgPath)
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		f := &bzl.File{}
		writeHeaders(f)
		if manageGoRules {
			reconcileLoad(path, f, rules)
		}
		writeRules(f, rules)
		return writeFile(path, f, false, dryRun)
	} else if err != nil {
		return false, err
	}
	if info.IsDir() {
		return false, fmt.Errorf("%q cannot be a directory", path)
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	f, err := bzl.Parse(path, b)
	if err != nil {
		return false, err
	}
	oldRules := make(map[string]*bzl.Rule)
	for _, r := range f.Rules("") {
		if r.Kind() == "proto_library" {
			if r.Attr("tags") == nil {
				r.SetAttr("tags", asExpr([]string{automanagedTag}))
			}
		}
		if r.Kind() == "go_proto_library" {
			if r.Attr("tags") == nil {
				r.SetAttr("tags", asExpr([]string{automanagedTag}))
			}
			goProtoLibrary = append(goProtoLibrary, ":"+r.Name())
		}
		if (r.Kind() == "go_library" || r.Kind() == "go_test") && (len(rules) == 3 || len(rules) == 4) && !strings.Contains(pkgPath, "vendor") {
			if r.Attr("tags") == nil {
				r.SetAttr("tags", asExpr([]string{automanagedTag}))
			}
		}
		if r.Kind() == "go_library" || r.Kind() == "go_test" {
			if listExpr, ok := r.Attr("deps").(*bzl.ListExpr); ok {
				olddeps := []string{}
				for _, v := range listExpr.List {
					olddeps = append(olddeps, v.(*bzl.StringExpr).Value)
				}
				newdeps := depMapping(olddeps)
				r.SetAttr("deps", asExpr(newdeps))
			}
		}
		oldRules[r.Name()] = r
	}
	if len(goProtoLibrary) > 0 && goProtoLibrary != nil {
		r, ok := oldRules["go_default_library"]
		if ok {
			r.SetAttr("embed", asExpr(goProtoLibrary))
			oldRules["go_default_library"] = r
		}
	}

	for _, r := range rules {
		o, ok := oldRules[r.Name()]
		if !ok {
			f.Stmt = append(f.Stmt, r.Call)
			continue
		}
		if !RuleIsManaged(o, manageGoRules) {
			continue
		}
		reconcileAttr := func(o, n *bzl.Rule, name string) {
			if e := n.Attr(name); e != nil {
				o.SetAttr(name, e)
			} else {
				o.DelAttr(name)
			}
		}
		for _, attr := range managedAttrs {
			reconcileAttr(o, r, attr)
		}
		delete(oldRules, r.Name())
	}

	for _, r := range oldRules {
		if !RuleIsManaged(r, manageGoRules) {
			continue
		}
		f.DelRules(r.Kind(), r.Name())
	}
	if manageGoRules {
		reconcileLoad(path, f, f.Rules(""))
	}

	return writeFile(path, f, true, dryRun)
}

func reconcileLoad(path string, f *bzl.File, rules []*bzl.Rule) {

	contains := func(s []string, e string) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	usedRuleKindsMap := map[string][]string{}
	for _, r := range rules {
		// Select only the Go rules we need to import, excluding builtins like filegroup.
		// TODO: make less fragile
		switch r.Kind() {
		case "go_prefix", "go_library", "go_binary", "go_test", "cgo_genrule", "cgo_library":
			if !contains(usedRuleKindsMap["@io_bazel_rules_go//go:def.bzl"], r.Kind()) {
				usedRuleKindsMap["@io_bazel_rules_go//go:def.bzl"] = append(usedRuleKindsMap["@io_bazel_rules_go//go:def.bzl"], r.Kind())
			}
		case "gazelle":
			if !contains(usedRuleKindsMap["@bazel_gazelle//:def.bzl"], r.Kind()) {
				usedRuleKindsMap["@bazel_gazelle//:def.bzl"] = append(usedRuleKindsMap["@bazel_gazelle//:def.bzl"], r.Kind())
			}
		case "go_proto_library":
			if !contains(usedRuleKindsMap["@io_bazel_rules_go//proto:def.bzl"], r.Kind()) {
				usedRuleKindsMap["@io_bazel_rules_go//proto:def.bzl"] = append(usedRuleKindsMap["@io_bazel_rules_go//proto:def.bzl"], r.Kind())
			}
		}
	}
	usedRuleKindsList := []string{}
	for k := range usedRuleKindsMap {
		usedRuleKindsList = append(usedRuleKindsList, k)
	}
	sort.Strings(usedRuleKindsList)

	for _, r := range f.Rules("load") {
		args := bzl.Strings(&bzl.ListExpr{List: r.Call.List})
		if len(args) == 0 {
			continue
		}
		if !contains(usedRuleKindsList, args[0]) {
			continue
		}
		if len(usedRuleKindsMap[args[0]]) == 0 {
			if r.Name() != "" {
				f.DelRules(r.Kind(), r.Name())
			}
			continue
		}
		r.Call.List = asExpr(append(
			[]string{args[0]}, usedRuleKindsMap[args[0]]...,
		)).(*bzl.ListExpr).List
		delete(usedRuleKindsMap, args[0])
	}
	for k, v := range usedRuleKindsMap {
		rule :=
			&bzl.CallExpr{
				X: &bzl.LiteralExpr{Token: "load"},
			}
		rule.List = asExpr(append(
			[]string{k}, v...,
		)).(*bzl.ListExpr).List
		f.Stmt = append([]bzl.Expr{rule}, f.Stmt...)
	}
}

// RuleIsManaged returns whether the provided rule is managed by this tool,
// based on the tags set on the rule.
func RuleIsManaged(r *bzl.Rule, manageGoRules bool) bool {
	var automanaged bool
	if !manageGoRules && (strings.HasPrefix(r.Kind(), "go_") || strings.HasPrefix(r.Kind(), "cgo_")) {
		return false
	}
	for _, tag := range r.AttrStrings("tags") {
		if tag == automanagedTag {
			automanaged = true
			break
		}
	}
	return automanaged
}

func writeFile(path string, f *bzl.File, exists, dryRun bool) (bool, error) {
	var info bzl.RewriteInfo
	bzl.Rewrite(f, &info)
	out := bzl.Format(f)
	//if strings.Contains(path, "vendor") {
	//	return false, nil
	//}
	if exists {
		orig, err := ioutil.ReadFile(path)
		if err != nil {
			return false, err
		}
		if bytes.Compare(orig, out) == 0 {
			return false, nil
		}
		if *printDiff {
			Diff(orig, out)
		}
	}
	if dryRun {
		fmt.Fprintf(os.Stderr, "DRY-RUN: wrote %q\n", path)
		return true, nil
	}
	werr := ioutil.WriteFile(path, out, 0644)
	if werr == nil {
		fmt.Fprintf(os.Stderr, "wrote %q\n", path)
		iswrote = true
	}
	return werr == nil, werr
}

func context() *build.Context {
	return &build.Context{
		GOARCH:      "amd64",
		GOOS:        "linux",
		GOROOT:      build.Default.GOROOT,
		GOPATH:      build.Default.GOPATH,
		ReleaseTags: []string{"go1.1", "go1.2", "go1.3", "go1.4", "go1.5", "go1.6", "go1.7", "go1.8", "go1.9", "go1.10"},
		Compiler:    runtime.Compiler,
		CgoEnabled:  true,
	}
}

func walk(root string, walkFn filepath.WalkFunc) error {
	return nil
}

func depMapping(dep []string) []string {
	result := []string{}
	mapping := map[string]string{
		"//vendor/github.com/golang/protobuf/proto:go_default_library":                    "@com_github_golang_protobuf//proto:go_default_library",
		"//vendor/github.com/golang/protobuf/ptypes/any:go_default_library":               "@io_bazel_rules_go//proto/wkt:any_go_proto",
		"//vendor/github.com/golang/protobuf/jsonpb:go_default_library":                   "@com_github_golang_protobuf//jsonpb:go_default_library",
		"//vendor/github.com/golang/protobuf/protoc-gen-go/plugin:go_default_library":     "@com_github_golang_protobuf//protoc-gen-go/plugin:go_default_library",
		"//vendor/github.com/golang/protobuf/protoc-gen-go/descriptor:go_default_library": "@com_github_golang_protobuf//protoc-gen-go/descriptor:go_default_library",
		"//vendor/github.com/golang/protobuf/ptypes:go_default_library":                   "@com_github_golang_protobuf//ptypes:go_default_library_gen",
		"//vendor/github.com/golang/protobuf/ptypes/empty:go_default_library":             "@io_bazel_rules_go//proto/wkt:empty_go_proto",

		"//vendor/github.com/gogo/protobuf/gogoproto:go_default_library":       "@com_github_gogo_protobuf//gogoproto:go_default_library",
		"//vendor/github.com/gogo/protobuf/proto:go_default_library":           "@com_github_gogo_protobuf//proto:go_default_library",
		"//vendor/github.com/gogo/protobuf/protoc-gen-gogo:go_default_library": "@com_github_gogo_protobuf//protoc-gen-gogo:go_default_library",
		"//vendor/github.com/gogo/protobuf/sortkeys:go_default_library":        "@com_github_gogo_protobuf//sortkeys:go_default_library",
		"//vendor/github.com/gogo/protobuf/types:go_default_library":           "@com_github_gogo_protobuf//types:go_default_library",
		"//vendor/github.com/gogo/protobuf/jsonpb:go_default_library":          "@com_github_gogo_protobuf//jsonpb:go_default_library",

		"//vendor/google.golang.org/grpc/codes:go_default_library":                "@org_golang_google_grpc//codes:go_default_library",
		"//vendor/google.golang.org/grpc/credentials:go_default_library":          "@org_golang_google_grpc//credentials:go_default_library",
		"//vendor/google.golang.org/grpc/metadata:go_default_library":             "@org_golang_google_grpc//metadata:go_default_library",
		"//vendor/google.golang.org/grpc/peer:go_default_library":                 "@org_golang_google_grpc//peer:go_default_library",
		"//vendor/google.golang.org/grpc/status:go_default_library":               "@org_golang_google_grpc//status:go_default_library",
		"//vendor/google.golang.org/grpc/resolver:go_default_library":             "@org_golang_google_grpc//resolver:go_default_library",
		"//vendor/google.golang.org/grpc/balancer:go_default_library":             "@org_golang_google_grpc//balancer:go_default_library",
		"//vendor/google.golang.org/grpc/balancer/base:go_default_library":        "@org_golang_google_grpc//balancer/base:go_default_library",
		"//vendor/google.golang.org/grpc/connectivity:go_default_library":         "@org_golang_google_grpc//connectivity:go_default_library",
		"//vendor/google.golang.org/grpc:go_default_library":                      "@org_golang_google_grpc//:go_default_library",
		"//vendor/google.golang.org/grpc/grpclog:go_default_library":              "@org_golang_google_grpc//grpclog:go_default_library",
		"//vendor/google.golang.org/grpc/interop:go_default_library":              "@org_golang_google_grpc//interop:go_default_library",
		"//vendor/google.golang.org/grpc/interop/grpc_testing:go_default_library": "@org_golang_google_grpc//interop/grpc_testing:go_default_library",
		"//vendor/google.golang.org/grpc/stress/grpc_testing:go_default_library":  "@org_golang_google_grpc//stress/grpc_testing:go_default_library",
		"//vendor/google.golang.org/grpc/reflection:go_default_library":           "@org_golang_google_grpc//reflection:go_default_library",
		"//vendor/google.golang.org/grpc/testdata:go_default_library":             "@org_golang_google_grpc//testdata:go_default_library",
		"//vendor/google.golang.org/grpc/interop/server:go_default_library":       "@org_golang_google_grpc//interop/server:go_default_library",
		"//vendor/google.golang.org/grpc/interop/client:go_default_library":       "@org_golang_google_grpc//interop/client:go_default_library",
		"//vendor/google.golang.org/grpc/interop/http2:go_default_library":        "@org_golang_google_grpc//interop/http2:go_default_library",
		"//vendor/google.golang.org/grpc/stress/client:go_default_library":        "@org_golang_google_grpc//stress/client:go_default_library",
		"//vendor/google.golang.org/grpc/keepalive:go_default_library":            "@org_golang_google_grpc//keepalive:go_default_library",
		"//vendor/google.golang.org/grpc/encoding/gzip:go_default_library":        "@org_golang_google_grpc//encoding/gzip:go_default_library",
		"//vendor/google.golang.org/grpc/stats:go_default_library":                "@org_golang_google_grpc//stats:go_default_library",
		"//vendor/google.golang.org/grpc/tap:go_default_library":                  "@org_golang_google_grpc//tap:go_default_library",
		"//vendor/google.golang.org/grpc/encoding:go_default_library":             "@org_golang_google_grpc//encoding:go_default_library",

		"//vendor/google.golang.org/genproto/googleapis/rpc/status:go_default_library": "@org_golang_google_genproto//googleapis/rpc/status:go_default_library",

		"//vendor/golang.org/x/net/context:go_default_library":         "@org_golang_x_net//context:go_default_library",
		"//vendor/golang.org/x/net/http2:go_default_library":           "@org_golang_x_net//http2:go_default_library",
		"//vendor/golang.org/x/net/proxy:go_default_library":           "@org_golang_x_net//proxy:go_default_library",
		"//vendor/golang.org/x/net/html:go_default_library":            "@org_golang_x_net//html:go_default_library",
		"//vendor/golang.org/x/net/html/atom:go_default_library":       "@org_golang_x_net//html/atom:go_default_library",
		"//vendor/golang.org/x/net/http2/hpack:go_default_library":     "@org_golang_x_net//http2/hpack:go_default_library",
		"//vendor/golang.org/x/net/context/ctxhttp:go_default_library": "@org_golang_x_net//context/ctxhttp:go_default_library",
		"//vendor/golang.org/x/net/ipv4:go_default_library":            "@org_golang_x_net//ipv4:go_default_library",
		"//vendor/golang.org/x/net/ipv6:go_default_library":            "@org_golang_x_net//ipv6:go_default_library",
		"//vendor/golang.org/x/net/trace:go_default_library":           "@org_golang_x_net//trace:go_default_library",
		"//vendor/golang.org/x/net/websocket:go_default_library":       "@org_golang_x_net//websocket:go_default_library",
	}
	for _, v := range dep {
		mapdep, ok := mapping[v]
		if ok {
			result = append(result, mapdep)
		} else {
			result = append(result, v)
		}
	}
	return result
}

func protoMap(path string, dep []string) []string {
	result := []string{}

	removeMap := map[string]struct{}{
		"//library/time:go_default_library": struct{}{},
	}
	mapping := map[string]string{
		"github.com/gogo/protobuf/gogoproto/gogo.proto": "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
		"google/protobuf/any.proto":                     "@com_google_protobuf//:any_proto",
		"google/api/annotations.proto":                  "@go_googleapis//google/api:annotations_proto",
		"google/protobuf/descriptor.proto":              "@com_google_protobuf//:descriptor_proto",
		"google/protobuf/empty.proto":                   "@com_google_protobuf//:empty_proto",
	}
	for _, v := range dep {
		if _, ok := removeMap[v]; ok {
			continue
		}
		mapdep, ok := mapping[v]
		if ok {
			result = append(result, mapdep)
		} else {
			if custom := customgoproto(path, v); custom != "" {
				result = append(result, custom)
			}
		}
	}
	return result
}

func goProtoMap(path string, dep []string) []string {
	result := []string{}
	mapping := map[string]string{
		// gogo
		"github.com/gogo/protobuf/gogoproto/gogo.proto": "@com_github_gogo_protobuf//gogoproto:go_default_library",
		// googleapis
		"google/api/annotations.proto": "@go_googleapis//google/api:annotations_go_proto",
		"google/rpc/errdetails.proto":  "@go_googleapis//google/rpc:errdetails_go_proto",
		"google/rpc/code.proto":        "@go_googleapis//google/rpc:code_go_proto",
		"google/rpc/status.proto":      "@go_googleapis//google/rpc:status_go_proto",
		// golang protobuf
		"@com_github_golang_protobuf//ptypes/any:go_default_library": "@io_bazel_rules_go//proto/wkt:any_go_proto",
		// google protobuf
		"google/protobuf/wrappers.proto":   "@io_bazel_rules_go//proto/wkt:wrappers_go_proto",
		"google/protobuf/timestamp.proto":  "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
		"google/protobuf/struct.proto":     "@io_bazel_rules_go//proto/wkt:struct_go_proto",
		"google/protobuf/field.proto":      "@io_bazel_rules_go//proto/wkt:field_mask_go_proto",
		"google/protobuf/empty.proto":      "@io_bazel_rules_go//proto/wkt:empty_go_proto",
		"google/protobuf/duration.proto":   "@io_bazel_rules_go//proto/wkt:duration_go_proto",
		"google/protobuf/compiler.proto":   "@io_bazel_rules_go//proto/wkt:compiler_plugin_go_proto",
		"google/protobuf/descriptor.proto": "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
		"google/protobuf/api.proto":        "@io_bazel_rules_go//proto/wkt:api_go_proto",
		"google/protobuf/type.proto":       "@io_bazel_rules_go//proto/wkt:type_go_proto",
		"google/protobuf/source.proto":     "@io_bazel_rules_go//proto/wkt:source_context_go_proto",
		"google/protobuf/any.proto":        "@io_bazel_rules_go//proto/wkt:any_go_proto",
	}
	for _, v := range dep {
		mapdep, ok := mapping[v]
		if ok {
			result = append(result, mapdep)
		} else {
			if custom := customgoprotolibrary(path, v); custom != "" {
				result = append(result, custom)
			}
		}
	}
	return result
}

func addExpr(x []bzl.Expr, y []bzl.Expr) []bzl.Expr {
	return append(x, y...)
}

func customgoprotolibrary(path, dep string) string {
	if strings.HasPrefix(dep, "library") || strings.HasPrefix(dep, "app") && strings.HasSuffix(dep, ".proto") {
		deplist := strings.Split(dep, "/")
		last := deplist[:len(deplist)-1]
		if strings.Join(last, "/") == path {
			return ""
		}
		last[len(last)-1] = last[len(last)-1] + ":" + last[len(last)-1] + "_go_proto"
		dep = strings.Join(last, "/")
		return "//" + dep
	}
	return dep
}

func customgoproto(path, dep string) string {
	if strings.HasPrefix(dep, "library") || strings.HasPrefix(dep, "app") && strings.HasSuffix(dep, ".proto") {
		deplist := strings.Split(dep, "/")
		last := deplist[:len(deplist)-1]
		if strings.Join(last, "/") == path {
			return ""
		}
		last[len(last)-1] = last[len(last)-1] + ":" + last[len(last)-1] + "_proto"
		dep = strings.Join(last, "/")
		return "//" + dep
	}
	return dep
}
