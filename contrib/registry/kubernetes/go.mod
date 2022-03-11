module github.com/SeeMusic/kratos/contrib/registry/kubernetes/v2

go 1.16

require (
	github.com/SeeMusic/kratos/v2 v2.2.0
	github.com/json-iterator/go v1.1.12
	k8s.io/api v0.23.1
	k8s.io/apimachinery v0.23.1
	k8s.io/client-go v0.23.1
)

replace github.com/SeeMusic/kratos/v2 => ../../../
