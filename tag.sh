#!/bin/bash
tag=$1

modules=(
"" "cmd/kratos/" "cmd/protoc-gen-errors/" "cmd/protoc-gen-go-http/"
"contrib/config/apollo/" "contrib/config/consul/" "contrib/config/etcd/" "contrib/config/kubernetes/" "contrib/config/nacos/"
"contrib/encoding/magpack/" "contrib/log/aliyun/" "contrib/log/zap/" "contrib/log/fluent/" "contrib/metrics/datadog/"
"contrib/metrics/prometheus/" "contrib/opensergo/" "contrib/registry/consul/" "contrib/registry/discovery/" "contrib/registry/etcd/"
"contrib/registry/eureka/" "contrib/registry/kubernetes/" "contrib/registry/nacos/" "contrib/registry/polaris/" "contrib/registy/zookeeper/"
)

for(( i=0;i<${#modules[@]};i++)) do
  git tag "${modules[i]}${tag}"
done
