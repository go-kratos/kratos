package kubernetes

import (
	"context"
	"log"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/config"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	testKey   = "test_config.json"
	namespace = "default"
	name      = "test"
)

var (
	keyPath    = strings.Join([]string{namespace, name, testKey}, "/")
	objectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels: map[string]string{
			"app": "test",
		},
	}
)

func TestSource(t *testing.T) {
	home := homedir.HomeDir()
	s := NewSource(
		Namespace("default"),
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

func TestConfig(t *testing.T) {
	restConfig, err := rest.InClusterConfig()
	home := homedir.HomeDir()

	options := []Option{
		Namespace(namespace),
		LabelSelector("app=test"),
	}

	if err != nil {
		kubeconfig := filepath.Join(home, ".kube", "config")
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			t.Fatal(err)
		}
		options = append(options, KubeConfig(kubeconfig))
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		t.Fatal(err)
	}

	clientSetConfigMaps := clientSet.CoreV1().ConfigMaps(namespace)

	source := NewSource(options...)
	if _, err = clientSetConfigMaps.Create(context.Background(), &v1.ConfigMap{
		ObjectMeta: objectMeta,
		Data: map[string]string{
			testKey: "test config",
		},
	}, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = clientSetConfigMaps.Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
			t.Error(err)
		}
	}()
	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(kvs) != 1 || kvs[0].Key != keyPath || string(kvs[0].Value) != "test config" {
		t.Fatal("config error")
	}

	w, err := source.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()
	// create also produce an event, discard it
	if _, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	if _, err = clientSetConfigMaps.Update(context.Background(), &v1.ConfigMap{
		ObjectMeta: objectMeta,
		Data: map[string]string{
			testKey: "new config",
		},
	}, metav1.UpdateOptions{}); err != nil {
		t.Error(err)
	}

	if kvs, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	if len(kvs) != 1 || kvs[0].Key != keyPath || string(kvs[0].Value) != "new config" {
		t.Fatal("config error")
	}
}

func TestExtToFormat(t *testing.T) {
	restConfig, err := rest.InClusterConfig()
	home := homedir.HomeDir()

	options := []Option{
		Namespace(namespace),
		LabelSelector("app=test"),
	}

	if err != nil {
		kubeconfig := filepath.Join(home, ".kube", "config")
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			t.Fatal(err)
		}
		options = append(options, KubeConfig(kubeconfig))
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		t.Fatal(err)
	}

	clientSetConfigMaps := clientSet.CoreV1().ConfigMaps(namespace)

	tc := `{"a":1}`
	if _, err = clientSetConfigMaps.Create(context.Background(), &v1.ConfigMap{
		ObjectMeta: objectMeta,
		Data: map[string]string{
			testKey: tc,
		},
	}, metav1.CreateOptions{}); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = clientSetConfigMaps.Delete(context.Background(), name, metav1.DeleteOptions{}); err != nil {
			t.Error(err)
		}
	}()

	source := NewSource(options...)
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(len(kvs), 1) {
		t.Errorf("len(kvs) = %d", len(kvs))
	}
	if !reflect.DeepEqual(keyPath, kvs[0].Key) {
		t.Errorf("kvs[0].Key is %s", kvs[0].Key)
	}
	if !reflect.DeepEqual(tc, string(kvs[0].Value)) {
		t.Errorf("kvs[0].Value is %s", kvs[0].Value)
	}
	if !reflect.DeepEqual("json", kvs[0].Format) {
		t.Errorf("kvs[0].Format is %s", kvs[0].Format)
	}
}
