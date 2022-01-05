module github.com/go-kratos/kratos/contrib/config/kubernetes/v2

go 1.16

require (
	github.com/go-kratos/kratos/v2 v2.1.4
	k8s.io/api v0.23.1
	k8s.io/apimachinery v0.23.1
	k8s.io/client-go v0.23.1
)

replace github.com/go-kratos/kratos/v2 => ../../../
