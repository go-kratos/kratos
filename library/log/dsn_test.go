package log

import (
	"flag"
	"testing"
	"time"

	xtime "go-common/library/time"

	"github.com/stretchr/testify/assert"
)

func TestUnixParseDSN(t *testing.T) {
	dsn := "unix:///data/log/lancer.sock?chan=1024&timeout=5s"
	flag.Parse()
	flag.Set("log.stdout", "true")
	as := parseDSN(dsn)
	assert.Equal(t, "unix", as.Proto)
	assert.Equal(t, "/data/log/lancer.sock", as.Addr)
	assert.Equal(t, 1024, as.Chan)
	assert.Equal(t, xtime.Duration(5*time.Second), as.Timeout)
}
