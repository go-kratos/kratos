module github.com/go-kratos/kratos/contrib/registry/kubernetes/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.1.2
	github.com/json-iterator/go v1.1.11
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
)

replace github.com/go-kratos/kratos/v2 => ../../../
