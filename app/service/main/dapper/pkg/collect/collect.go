package collect

import (
	"go-common/app/service/main/dapper/pkg/process"
)

// Collecter collect span from different source
type Collecter interface {
	Start() error
	RegisterProcess(p process.Processer)
	Close() error
}
