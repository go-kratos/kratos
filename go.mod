module github.com/bilibili/kratos

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/aristanetworks/goarista v0.0.0-20190409234242-46f4bc7b73ef // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cznic/b v0.0.0-20181122101859-a26611c4d92d // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/cznic/strutil v0.0.0-20181122101858-275e90344537 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dgryski/go-farm v0.0.0-20190323231341-8198c7b169ec
	github.com/fatih/color v1.7.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.2.0
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/golang/mock v1.3.1 // indirect
	github.com/golang/protobuf v1.3.1
	github.com/google/btree v1.0.0 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.4 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/montanaflynn/stats v0.5.0
	github.com/openzipkin/zipkin-go v0.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.2
	github.com/remyoudompheng/bigfft v0.0.0-20190321074620-2f0d2b0e0001 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/sirupsen/logrus v1.4.1
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/tsuna/gohbase v0.0.0-20190201102810-d3184c1526df
	github.com/urfave/cli v1.20.0
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/genproto v0.0.0-20190701230453-710ae3a149df
	google.golang.org/grpc v1.22.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.26.0
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.26.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190605123033-f99c8df09eb5
	golang.org/x/lint => golang.org/x/lint v0.0.0-20190511005446-959b441ac422
	golang.org/x/net => golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.0.0-20190614102709-0f29369cfe45
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190427212804-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190602015325-4c4f7f33c9ed
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190606050223-4d9ae51c2468
	google.golang.org/appengine => github.com/golang/appengine v1.6.1
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190701230453-710ae3a149df
	google.golang.org/grpc => github.com/grpc/grpc-go v1.22.0
)
