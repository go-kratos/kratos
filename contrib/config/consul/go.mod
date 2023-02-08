module github.com/go-kratos/kratos/contrib/config/consul/v2

go 1.15

require (
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/go-kratos/kratos/v2 v2.4.0
	github.com/google/btree v1.0.0 // indirect
	github.com/hashicorp/consul/api v1.18.0
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-msgpack v0.5.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
