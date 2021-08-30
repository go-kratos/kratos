package kubernetes

import (
	"log"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	"k8s.io/client-go/util/homedir"
)

// 部署在mesh namespace 下configmap
const yamlData = `database:
  mysql:
    dsn: "root:Test@tcp(mysql.database.svc.cluster.local:3306)/test?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
    active: 20
    idle: 10
    idle_timeout: 3600
  redis:
    addr: "redis-master.redis.svc.cluster.local:6379"
    password: ""
    db: 4`

const yamlApp = `application:
  expire: 3600`

func main() {
	conf := config.New(
		config.WithSource(
			NewSource(
				Namespace("mesh"),
				LabelSelector("app=test"),
				KubeConfig(filepath.Join(homedir.HomeDir(), ".kubernetes", "config")),
			),
		),
	)
	err := conf.Load()
	if err != nil {
		log.Panic(err)
	}
}
