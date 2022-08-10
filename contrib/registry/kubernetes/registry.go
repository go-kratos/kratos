// Package kuberegistry registry simply implements the Kubernetes-based Registry
package kuberegistry

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	jsoniter "github.com/json-iterator/go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// Defines the key name of specific fields
// Kratos needs to cooperate with the following fields to run properly on Kubernetes:
// kratos-service-id: define the ID of the service
// kratos-service-app: define the name of the service
// kratos-service-version: define the version of the service
// kratos-service-metadata: define the metadata of the service
// kratos-service-protocols: define the protocols of the service
//
// Example Deployment:
/*
apiVersion: apps/v1
kind: Deployment
metadata:
name: nginx
labels:
  app: nginx
spec:
replicas: 5
selector:
  matchLabels:
    app: nginx
template:
  metadata:
    labels:
      app: nginx
      kratos-service-id: "56991810-c77f-4a95-8190-393efa9c1a61"
      kratos-service-app: "nginx"
      kratos-service-version: "v3.5.0"
    annotations:
      kratos-service-protocols: |
        {"80": "http"}
      kratos-service-metadata: |
        {"region": "sh", "zone": "sh001", "cluster": "pd"}
  spec:
    containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
          - containerPort: 80
*/
const (
	// LabelsKeyServiceID is used to define the ID of the service
	LabelsKeyServiceID = "kratos-service-id"
	// LabelsKeyServiceName is used to define the name of the service
	LabelsKeyServiceName = "kratos-service-app"
	// LabelsKeyServiceVersion is used to define the version of the service
	LabelsKeyServiceVersion = "kratos-service-version"
	// AnnotationsKeyMetadata is used to define the metadata of the service
	AnnotationsKeyMetadata = "kratos-service-metadata"
	// AnnotationsKeyProtocolMap is used to define the protocols of the service
	// Through the value of this field, Kratos can obtain the application layer protocol corresponding to the port
	// Example value: {"80": "http", "8081": "grpc"}
	AnnotationsKeyProtocolMap = "kratos-service-protocols"
)

// The Registry simply implements service discovery based on Kubernetes
// It has not been verified in the production environment and is currently for reference only
type Registry struct {
	clientSet       *kubernetes.Clientset
	informerFactory informers.SharedInformerFactory
	podInformer     cache.SharedIndexInformer
	podLister       listerv1.PodLister

	stopCh chan struct{}
}

// NewRegistry is used to initialize the Registry
func NewRegistry(clientSet *kubernetes.Clientset) *Registry {
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Minute*10)
	podInformer := informerFactory.Core().V1().Pods().Informer()
	podLister := informerFactory.Core().V1().Pods().Lister()
	return &Registry{
		clientSet:       clientSet,
		informerFactory: informerFactory,
		podInformer:     podInformer,
		podLister:       podLister,
		stopCh:          make(chan struct{}),
	}
}

// Register is used to register services
// Note that on Kubernetes, it can only be used to update the id/name/version/metadata/protocols of the current service,
// but it cannot be used to update node.
func (s *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	// GetMetadata
	metadataVal, err := marshal(service.Metadata)
	if err != nil {
		return err
	}

	// Generate ProtocolMap
	protocolMap, err := getProtocolMapByEndpoints(service.Endpoints)
	if err != nil {
		return err
	}
	protocolMapVal, err := marshal(protocolMap)
	if err != nil {
		return err
	}

	patchBytes, err := jsoniter.Marshal(map[string]interface{}{
		"metadata": metav1.ObjectMeta{
			Labels: map[string]string{
				LabelsKeyServiceID:      service.ID,
				LabelsKeyServiceName:    service.Name,
				LabelsKeyServiceVersion: service.Version,
			},
			Annotations: map[string]string{
				AnnotationsKeyMetadata:    metadataVal,
				AnnotationsKeyProtocolMap: protocolMapVal,
			},
		},
	})
	if err != nil {
		return err
	}

	if _, err = s.clientSet.
		CoreV1().
		Pods(GetNamespace()).
		Patch(ctx, GetPodName(), types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{}); err != nil {
		return err
	}
	return nil
}

// Deregister the registration.
func (s *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return s.Register(ctx, &registry.ServiceInstance{
		Metadata: map[string]string{},
	})
}

// GetService return the service instances in memory according to the service name.
func (s *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	pods, err := s.podLister.List(labels.SelectorFromSet(map[string]string{
		LabelsKeyServiceName: name,
	}))
	if err != nil {
		return nil, err
	}
	ret := make([]*registry.ServiceInstance, 0, len(pods))
	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		instance, err := getServiceInstanceFromPod(pod)
		if err != nil {
			return nil, err
		}
		ret = append(ret, instance)
	}
	return ret, nil
}

func (s *Registry) sendLatestInstances(ctx context.Context, name string, announcement chan []*registry.ServiceInstance) {
	instances, err := s.GetService(ctx, name)
	if err != nil {
		panic(err)
	}
	announcement <- instances
}

// Watch creates a watcher according to the service name.
func (s *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	stopCh := make(chan struct{}, 1)
	announcement := make(chan []*registry.ServiceInstance, 1)
	s.podInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			select {
			case <-stopCh:
				return false
			case <-s.stopCh:
				return false
			default:
				pod := obj.(*corev1.Pod)
				val := pod.GetLabels()[LabelsKeyServiceName]
				return val == name
			}
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				s.sendLatestInstances(ctx, name, announcement)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				s.sendLatestInstances(ctx, name, announcement)
			},
			DeleteFunc: func(obj interface{}) {
				s.sendLatestInstances(ctx, name, announcement)
			},
		},
	})
	return NewIterator(announcement, stopCh), nil
}

// Start is used to start the Registry
// It is non-blocking
func (s *Registry) Start() {
	s.informerFactory.Start(s.stopCh)
	if !cache.WaitForCacheSync(s.stopCh, s.podInformer.HasSynced) {
		return
	}
}

// Close is used to close the Registry
// After closing, any callbacks generated by Watch will not be executed
func (s *Registry) Close() {
	select {
	case <-s.stopCh:
	default:
		close(s.stopCh)
	}
}

// //////////// K8S Runtime ////////////

// ServiceAccountNamespacePath defines the location of the namespace file
const ServiceAccountNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

var currentNamespace = LoadNamespace()

// LoadNamespace is used to get the current namespace from the file
func LoadNamespace() string {
	data, err := os.ReadFile(ServiceAccountNamespacePath)
	if err != nil {
		return ""
	}
	return string(data)
}

// GetNamespace is used to get the namespace of the Pod where the current container is located
func GetNamespace() string {
	return currentNamespace
}

// GetPodName is used to get the name of the Pod where the current container is located
func GetPodName() string {
	return os.Getenv("HOSTNAME")
}

// //////////// ProtocolMap ////////////

type protocolMap map[string]string

func (m protocolMap) GetProtocol(port int32) string {
	return m[strconv.Itoa(int(port))]
}

// //////////// Iterator ////////////

// Iterator performs the conversion from channel to iterator
// It reads the latest changes from the `chan []*registry.ServiceInstance`
// And the outside can sense the closure of Iterator through stopCh
type Iterator struct {
	ch     chan []*registry.ServiceInstance
	stopCh chan struct{}
}

// NewIterator is used to initialize Iterator
func NewIterator(channel chan []*registry.ServiceInstance, stopCh chan struct{}) *Iterator {
	return &Iterator{
		ch:     channel,
		stopCh: stopCh,
	}
}

// Next will block until ServiceInstance changes
func (iter *Iterator) Next() ([]*registry.ServiceInstance, error) {
	select {
	case instances := <-iter.ch:
		return instances, nil
	case <-iter.stopCh:
		return nil, ErrIteratorClosed
	}
}

// Close is used to close the iterator
func (iter *Iterator) Stop() error {
	select {
	case <-iter.stopCh:
	default:
		close(iter.stopCh)
	}
	return nil
}

// //////////// Helper Func ////////////

func marshal(in interface{}) (string, error) {
	return jsoniter.MarshalToString(in)
}

func unmarshal(data string, in interface{}) error {
	return jsoniter.UnmarshalFromString(data, in)
}

func isEmptyObjectString(s string) bool {
	switch s {
	case "", "{}", "null", "nil", "[]":
		return true
	}
	return false
}

func getProtocolMapByEndpoints(endpoints []string) (protocolMap, error) {
	ret := protocolMap{}
	for _, endpoint := range endpoints {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		ret[u.Port()] = u.Scheme
	}
	return ret, nil
}

func getProtocolMapFromPod(pod *corev1.Pod) (protocolMap, error) {
	protoMap := protocolMap{}
	if s := pod.Annotations[AnnotationsKeyProtocolMap]; !isEmptyObjectString(s) {
		err := unmarshal(s, &protoMap)
		if err != nil {
			return nil, &ErrorHandleResource{Namespace: pod.Namespace, Name: pod.Name, Reason: err}
		}
	}
	return protoMap, nil
}

func getMetadataFromPod(pod *corev1.Pod) (map[string]string, error) {
	metadata := map[string]string{}
	if s := pod.Annotations[AnnotationsKeyMetadata]; !isEmptyObjectString(s) {
		err := unmarshal(s, &metadata)
		if err != nil {
			return nil, &ErrorHandleResource{Namespace: pod.Namespace, Name: pod.Name, Reason: err}
		}
	}
	return metadata, nil
}

func getServiceInstanceFromPod(pod *corev1.Pod) (*registry.ServiceInstance, error) {
	podIP := pod.Status.PodIP
	podLabels := pod.GetLabels()
	// Get Metadata
	metadata, err := getMetadataFromPod(pod)
	if err != nil {
		return nil, err
	}
	// Get Protocols Definition
	protocolMap, err := getProtocolMapFromPod(pod)
	if err != nil {
		return nil, err
	}

	// Get Endpoints
	var endpoints []string
	for _, container := range pod.Spec.Containers {
		for _, cp := range container.Ports {
			port := cp.ContainerPort
			protocol := protocolMap.GetProtocol(port)
			if protocol == "" {
				if cp.Name != "" {
					protocol = strings.Split(cp.Name, "-")[0]
				} else {
					protocol = string(cp.Protocol)
				}
			}
			addr := fmt.Sprintf("%s://%s:%d", protocol, podIP, port)
			endpoints = append(endpoints, addr)
		}
	}
	return &registry.ServiceInstance{
		ID:        podLabels[LabelsKeyServiceID],
		Name:      podLabels[LabelsKeyServiceName],
		Version:   podLabels[LabelsKeyServiceVersion],
		Metadata:  metadata,
		Endpoints: endpoints,
	}, nil
}

// //////////// Error Definition ////////////

// ErrIteratorClosed defines the error that the iterator is closed
var ErrIteratorClosed = errors.New("iterator closed")

// ErrorHandleResource defines the error that cannot handle K8S resources normally
type ErrorHandleResource struct {
	Namespace string
	Name      string
	Reason    error
}

// Error implements the error interface
func (err *ErrorHandleResource) Error() string {
	return fmt.Sprintf("failed to handle resource(namespace=%s, name=%s): %s",
		err.Namespace, err.Name, err.Reason)
}
