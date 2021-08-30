module github.com/go-kratos/kratos/registry/kube/v2

go 1.17

require (
	github.com/go-kratos/kratos/v2 v2.0.5
	github.com/json-iterator/go v1.1.11
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
)

replace github.com/go-kratos/kratos/v2 => ../../../kratos