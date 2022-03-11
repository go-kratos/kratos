package kubernetes

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/SeeMusic/kratos/v2/config"
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

func ExampleNewSource() {
	conf := config.New(
		config.WithSource(
			NewSource(
				Namespace("mesh"),
				LabelSelector("app=test"),
				KubeConfig(filepath.Join(homedir.HomeDir(), ".kube", "config")),
			),
		),
	)
	err := conf.Load()
	if err != nil {
		log.Panic(err)
	}
}
