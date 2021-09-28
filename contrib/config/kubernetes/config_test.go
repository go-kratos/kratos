package kubernetes

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"k8s.io/client-go/util/homedir"
)

func TestSource(t *testing.T) {
	home := homedir.HomeDir()
	s := NewSource(
		Namespace("mesh"),
		LabelSelector(""),
		KubeConfig(filepath.Join(home, ".kube", "config")),
	)
	kvs, err := s.Load()
	if err != nil {
		t.Error(err)
	}
	for _, v := range kvs {
		t.Log(v)
	}
}

func ExampleNewSource() {
	/*
			部署在mesh namespace 下configmap
			yamlData 示例
			database:
			  mysql:
				dsn: "root:Test@tcp(mysql.database.svc.cluster.local:3306)/test?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
				active: 20
				idle: 10
				idle_timeout: 3600
			  redis:
				addr: "redis-master.redis.svc.cluster.local:6379"
				password: ""
				db: 4

			yamlApp 示例
			application:
		  		expire: 3600
	*/
	conf := config.New(
		config.WithSource(
			NewSource(
				Namespace("mesh"),
				LabelSelector("app=test"),
				KubeConfig(filepath.Join(homedir.HomeDir(), ".kube", "config")),
			),
		),
	)
	err := conf.Load()
	if err != nil {
		log.Panic(err)
	}
}
