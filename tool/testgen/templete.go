package main

var (
	tpPackage     = "package %s\n\n"
	tpImport      = "import (\n\t%s\n)\n\n"
	tpVar         = "var (\n\t%s\n)\n"
	tpInterface   = "type %sInterface interface {\n%s}\n"
	tpIntfcFunc   = "%s(%s) %s\n"
	tpMonkeyFunc  = "// Mock%s .\nfunc Mock%s(%s %s,%s) (guard *monkey.PatchGuard) {\n\treturn monkey.PatchInstanceMethod(reflect.TypeOf(%s), \"%s\", func(_ %s, %s) (%s) {\n\t\treturn %s\n\t})\n}\n\n"
	tpTestReset   = "\n\t\tReset(func() {%s\n\t\t})"
	tpTestFunc    = "func Test%s%s(t *testing.T){%s\n\tConvey(\"%s\", t, func(){\n\t\t%s\tConvey(\"When everything goes positive\", func(){\n\t\t\t%s\n\t\t\t})\n\t\t})%s\n\t})\n}\n\n"
	tpTestDaoMain = `func TestMain(m *testing.M) {
	flag.Set("conf", "%s")
	flag.Parse()
	%s
	os.Exit(m.Run())
}
`
	tpTestServiceMain = `func TestMain(m *testing.M){
	flag.Set("conf", "%s")
	flag.Parse()
	%s
	os.Exit(m.Run())
}
`
	tpTestMainNew = `if err := paladin.Init(); err != nil {
		panic(err)
	}
	%s`
	tpTestMainOld = `if err := conf.Init(); err != nil {
		panic(err)
	}
	%s`
	print = `Generation success!                                            
                          莫生气
                       代码辣鸡非我意,
                       自己动手分田地;
                       你若气死谁如意?
                       谈笑风生活长命.
// Release 1.2.3. Powered by Kratos`
)
