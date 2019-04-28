package dao

import (
	"flag"
	"go-common/app/service/main/antispam/conf"
	"os"
	"testing"
)

var (
	d     *Dao
	kwi   *KeywordDaoImpl
	regdi *RegexpDaoImpl
	rdi   *RuleDaoImpl
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.antispam-service")
		flag.Set("conf_token", "e0de72afaf4946ca836e9b7b459b833b")
		flag.Set("tree_id", "11041")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(""); err != nil {
		panic(err)
	}
	Init(conf.Conf)
	d = New(conf.Conf)
	kwi = NewKeywordDao()
	regdi = NewRegexpDao()
	rdi = NewRuleDao()
	m.Run()
	os.Exit(0)
}
