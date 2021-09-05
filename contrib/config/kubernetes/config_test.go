package kubernetes

import (
	"path/filepath"
	"testing"

	"k8s.io/client-go/util/homedir"
)

func TestSource(t *testing.T) {
	home := homedir.HomeDir()
	s := NewSource(
		Namespace("mesh"),
		LabelSelector(""),
		KubeConfig(filepath.Join(home, ".kube", "config")),
	)
	kvs, err := s.Load()
	if err != nil {
		t.Error(err)
	}
	for _, v := range kvs {
		t.Log(v)
	}
}
