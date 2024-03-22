package metrics

import "os"

func EnableOTELExemplar() error {
	return os.Setenv("OTEL_GO_X_EXEMPLAR", "true")
}
