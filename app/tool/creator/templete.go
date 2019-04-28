package main

var (
	tpPackage     = "package %s\n\n"
	tpImport      = "import (\n\t%s\n)\n\n"
	tpVar         = "var (\n\t%s\n\t\t)\n"
	tpInterface   = "type %sInterface interface {\n%s}\n"
	tpIntfcFunc   = "%s(%s) %s\n"
	tpMonkeyFunc  = "// Mock%s .\nfunc Mock%s(%s %s,%s) (guard *monkey.PatchGuard) {\n\treturn monkey.PatchInstanceMethod(reflect.TypeOf(%s), \"%s\", func(_ %s, %s) (%s) {\n\t\treturn %s\n\t})\n}\n\n"
	tpTestReset   = "\n\t\tconvCtx.Reset(func() {%s\n\t\t})"
	tpTestFunc    = "func Test%s%s(t *testing.T){\n\tconvey.Convey(\"%s\", t, func(convCtx convey.C){\n\t\t%s\tconvCtx.Convey(\"When everything goes positive\", func(convCtx convey.C){\n\t\t\t%s\n\t\t\t})\n\t\t})%s\n\t})\n}\n\n"
	tpTestDaoMain = `func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "")
		flag.Set("conf_token", "")
		flag.Set("tree_id", "")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}else{
		flag.Set("conf", "%s")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
`
	tpTestServiceMain = `func TestMain(m *testing.M){
		flag.Set("conf", "%s")
	 	flag.Parse()
	 	if err := conf.Init(); err != nil {
		 	panic(err)
	 	}
		s = New(conf.Conf)
	 	os.Exit(m.Run())
 }
`
	moha = `Generation success!

            $$$$$$$$$$$$$$$       $$$$$$$$$$$$$$$%.
         &$               &$    =$                $         
         $                 $$$$$$=                @&        
      B=$$                 $$$$$$                  $$&=     
      $$$$      +1s        $$&-$$       +1s        $$$      
      $$$$                 $-   $                  $$$      
         $                 $    $                  .B        
         $                 @    $                 .=        
          $                $     $                %         
           $            =$         =$            @&          
            #$$$$$$$$$%              -@$$$$$$$$B            
                                                            
                      莫生        莫生
                       气         气

                       代码辣鸡非我意,
                       自己动手分田地;
                       你若气死谁如意?
                       谈笑风生活长命.

// Release 1.1.3. Powered by 主站质保团队`
)
