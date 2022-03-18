module github.com/go-kratos/kratos/contrib/registry/kubernetes/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.2.1
	github.com/json-iterator/go v1.1.12
	k8s.io/api v0.23.1
	k8s.io/apimachinery v0.23.1
	k8s.io/client-go v0.23.1
)

replace github.com/go-kratos/kratos/v2 => ../../../
