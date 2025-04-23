package kuberegistry

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	namespace  = "default"
	deployName = "hello-deployment"
	podName    = "hello"
)

var deployment = appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name: deployName,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": podName,
			},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": podName,
				},
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:  "nginx",
						Image: "nginx:alpine",
						Ports: []apiv1.ContainerPort{
							{
								Name:          "http",
								Protocol:      apiv1.ProtocolTCP,
								ContainerPort: 80,
							},
						},
						Command: []string{
							"nginx",
							"-g",
							"daemon off;",
						},
					},
				},
			},
		},
	},
}

func getClientSet() (*kubernetes.Clientset, error) {
	restConfig, err := rest.InClusterConfig()
	home := homedir.HomeDir()

	if err != nil {
		kubeconfig := filepath.Join(home, ".kube", "config")
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func int32Ptr(i int32) *int32 { return &i }

func TestSetEnv(t *testing.T) {
	os.Setenv("HOSTNAME", podName)
	if os.Getenv("HOSTNAME") != podName {
		t.Fatal("error")
	}
}

func TestRegistry(t *testing.T) {
	currentNamespace = "default"

	clientSet, err := getClientSet()
	if err != nil {
		t.Fatal(err)
	}

	r := NewRegistry(clientSet, currentNamespace)
	r.Start()

	svrHello := &registry.ServiceInstance{
		ID:        "1",
		Name:      "hello",
		Version:   "v1.0.0",
		Endpoints: []string{"http://127.0.0.1:80"},
	}
	_, err = clientSet.AppsV1().Deployments(namespace).Create(context.Background(), &deployment, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}

	watch, err := r.Watch(context.Background(), svrHello.Name)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = watch.Stop()
	}()

	go func() {
		for {
			res, err1 := watch.Next()
			if err1 != nil {
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	pod, err := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=hello",
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(pod.Items) < 1 {
		t.Fatal("fetch resource error")
	}

	os.Setenv("HOSTNAME", pod.Items[0].Name)

	// Always remember delete test resource
	defer func() {
		_ = clientSet.AppsV1().Deployments(namespace).Delete(context.Background(), deployName, metav1.DeleteOptions{})
	}()

	if err = r.Register(context.Background(), svrHello); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	res, err := r.GetService(context.Background(), svrHello.Name)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 && res[0].Name != svrHello.Name {
		t.Fatal(err)
	}

	if err1 := r.Deregister(context.Background(), svrHello); err1 != nil {
		t.Fatal(err1)
	}

	time.Sleep(time.Second)

	res, err = r.GetService(context.Background(), svrHello.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Fatal("not expected empty")
	}
}
