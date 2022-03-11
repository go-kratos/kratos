module github.com/SeeMusic/kratos/contrib/config/kubernetes/v2

go 1.16

require (
	github.com/SeeMusic/kratos/v2 v2.2.0
	k8s.io/api v0.23.3
	k8s.io/apimachinery v0.23.3
	k8s.io/client-go v0.23.3
)

replace github.com/SeeMusic/kratos/v2 => ../../../
