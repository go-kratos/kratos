package env

import (
	"testing"
)

func Test_watcher_next(t *testing.T) {
	t.Run("next after stop should return err", func(t *testing.T) {
		w, err := NewWatcher()
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}

		_ = w.Stop()
		_, err = w.Next()
		if err == nil {
			t.Error("expect error, actual nil")
		}
	})
}

func Test_watcher_stop(t *testing.T) {
	t.Run("stop multiple times should not panic", func(t *testing.T) {
		w, err := NewWatcher()
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}

		_ = w.Stop()
		_ = w.Stop()
	})
}
