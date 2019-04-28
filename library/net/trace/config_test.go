package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDSN(t *testing.T) {
	_, err := parseDSN(_traceDSN)
	if err != nil {
		t.Error(err)
	}
}

func TestTraceFromEnvFlag(t *testing.T) {
	_, err := TracerFromEnvFlag()
	if err != nil {
		t.Error(err)
	}
}

func TestInit(t *testing.T) {
	Init(nil)
	_, ok := _tracer.(nooptracer)
	assert.False(t, ok)

	_tracer = nooptracer{}

	Init(&Config{Network: "unixgram", Addr: "unixgram:///var/run/dapper-collect/dapper-collect.sock"})
	_, ok = _tracer.(nooptracer)
	assert.False(t, ok)
}
