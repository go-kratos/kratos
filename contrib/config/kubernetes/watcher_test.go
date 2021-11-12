package kubernetes

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func TestKube(t *testing.T) {
	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		t.Error(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Error(err)
	}
	cmWatcher, err := client.CoreV1().ConfigMaps("mesh").Watch(context.Background(), metav1.ListOptions{
		LabelSelector: "app=test",
		// FieldSelector:        "",
	})
	if err != nil {
		t.Error(err)
	}
	go func() {
		time.Sleep(5 * time.Second)
		cmWatcher.Stop()
	}()
	for c := range cmWatcher.ResultChan() {
		if c.Object == nil {
			return
		}
		t.Log(c.Type, c.Object)
	}
}
