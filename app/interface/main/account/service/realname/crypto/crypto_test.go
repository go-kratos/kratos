package crypto

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	rsaBase64 = make([]byte, 1000)
	aesBase64 = make([]byte, 5000000)
)

func TestEncrypt(t *testing.T) {
	Convey("encrypt", t, func() {
		// TODO
	})
}

func BenchmarkBytesByFMT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%04x%s%s", len(rsaBase64), rsaBase64, aesBase64)
	}
}

func BenchmarkBytesByBuffer(b *testing.B) {
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintf(&buf, "%04x", len(rsaBase64))
		buf.Write(rsaBase64)
		buf.Write(aesBase64)
		b.StopTimer()
		buf.Reset()
		b.StartTimer()
		_ = buf.Bytes()
	}
}
