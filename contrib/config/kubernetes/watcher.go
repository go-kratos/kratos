package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	k       *kube
	watcher watch.Interface
}

func newWatcher(k *kube) (config.Watcher, error) {
	w, err := k.client.CoreV1().ConfigMaps(k.opts.Namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: k.opts.LabelSelector,
		FieldSelector: k.opts.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return &watcher{
		k:       k,
		watcher: w,
	}, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
ResultChan:
	ch := <-w.watcher.ResultChan()
	if ch.Object == nil {
		// 重新获取watcher
		k8sWatcher, err := w.k.client.CoreV1().ConfigMaps(w.k.opts.Namespace).Watch(context.Background(), metav1.ListOptions{
			LabelSelector: w.k.opts.LabelSelector,
			FieldSelector: w.k.opts.FieldSelector,
		})
		if err != nil {
			return nil, err
		}
		w.watcher = k8sWatcher
		goto ResultChan
	}
	cm, ok := ch.Object.(*v1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("kubernetes Object not ConfigMap")
	}
	if ch.Type == "DELETED" {
		return nil, fmt.Errorf("kubernetes configmap delete %s", cm.Name)
	}
	return w.k.configMap(*cm), nil
}

func (w *watcher) Stop() error {
	w.watcher.Stop()
	return nil
}
