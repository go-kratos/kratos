package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	flagDep     = flag.String("dep", "main,live,openplatform,ep", "department list , split by comma")
	flagPrefix  = flag.String("prefix", `business`, "prefix path")
	flagService = flag.String("serivce", "interface,job,admin,service", "service type")
	// flagWhite prefix下允许的dir名称
	flagWhite = flag.String("white", "", "white subpath from prefix , split by comma")
)

const (
	codeSuccess = 0
	codeFail    = 1
)

func main() {
	flag.Parse()
	var (
		depList      []string
		serviceList  []string
		filePathList []string
		whiteDirList []string
	)
	filePathList = flag.Args()
	if len(filePathList) <= 0 {
		fmt.Println("No file to check")
		os.Exit(codeSuccess)
	}

	depList = strings.Split(*flagDep, ",")
	serviceList = strings.Split(*flagService, ",")
	for _, wd := range strings.Split(*flagWhite, ",") {
		if wd != "" {
			whiteDirList = append(whiteDirList, strings.Join([]string{*flagPrefix, wd}, "/"))
		}
	}
	code := check(filePathList, serviceList, depList, whiteDirList)
	os.Exit(code)
}

func check(filePathList []string, serviceTypeList []string, depList []string, whiteDirList []string) (code int) {
	var (
		regDep      = strings.Join(depList, "|")
		serviceType = strings.Join(serviceTypeList, "|")
		regStr      = fmt.Sprintf(`%s/(%s)/(%s)`, *flagPrefix, serviceType, regDep)
		reg         *regexp.Regexp
		flag        = true
		failedFiles []string
		err         error
	)
	regStr = strings.Replace(regStr, "/", `\/`, -1)
	if reg, err = regexp.Compile(regStr); err != nil {
		err = errors.Wrapf(err, "regexp : %s", regStr)
		fmt.Printf("%+v\n", err)
		code = codeFail
		return
	}
	for _, p := range filePathList {
		if strings.HasPrefix(p, *flagPrefix) {
			if whiteCheck(whiteDirList, p) {
				continue
			}
			if !reg.MatchString(p) {
				failedFiles = append(failedFiles, p)
				flag = false
				break
			}
		}
	}
	if !flag {
		fmt.Println("invalid files : ")
		for _, f := range failedFiles {
			fmt.Printf("\t%s\n", f)
		}
		code = codeFail
	} else {
		code = codeSuccess
	}
	return
}

func whiteCheck(whiteDirList []string, path string) bool {
	for _, wd := range whiteDirList {
		if strings.HasPrefix(path, wd) {
			return true
		}
	}
	return false
}
