package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddrDSN(t *testing.T) {
	t.Run("test parse addr dsn", func(t *testing.T) {
		addr := parseDSNAddr("test:test@tcp(172.16.0.148:3306)/test?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8")
		assert.Equal(t, "172.16.0.148:3306", addr)
	})
	t.Run("test password has @", func(t *testing.T) {
		addr := parseDSNAddr("root:root@dev@tcp(1.2.3.4:3306)/abc?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8")
		assert.Equal(t, "1.2.3.4:3306", addr)
	})
}
