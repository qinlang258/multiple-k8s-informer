package controller

import (
	"context"
	"multiple-k8s-informer/queue"
	"multiple-k8s-informer/resource"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

// ResourceAndNamespace 资源与namespace
type ResourceAndNamespace struct {
	RType     string `json:"rType" yaml:"rType"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

// 创建 "k8s.io/api/core/v1"的核心包
func (r *ResourceAndNamespace) CreateCoreV1IndexInformer(client *kubernetes.Clientset, worker queue.Queue, clusterName string) (indexer cache.Indexer, informer cache.Controller) {

	restClient := client.CoreV1().RESTClient()
	lw := cache.NewListWatchFromClient(restClient, r.RType, r.Namespace, fields.Everything())

	switch r.RType {
	case resource.Services:
		indexer, informer = cache.NewIndexerInformer(lw, &v1.Service{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
	case resource.Pods:
		indexer, informer = cache.NewIndexerInformer(lw, &v1.Pod{}, 0, InitHandleFunc(resource.Pods, clusterName, worker), cache.Indexers{})
	case resource.ConfigMaps:
		indexer, informer = cache.NewIndexerInformer(lw, &v1.ConfigMap{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
	case resource.Secrets:
		indexer, informer = cache.NewIndexerInformer(lw, &v1.Secret{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
	case resource.Events:
		indexer, informer = cache.NewIndexerInformer(lw, &v1.Event{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
	}
	return
}

// appsv1 "k8s.io/api/apps/v1" 构造informer需要的资源
func (r *ResourceAndNamespace) CreateAppsV1IndexInformer(client *kubernetes.Clientset, worker queue.Queue, clusterName string) (indexer cache.Indexer, informer cache.Controller) {
	restClient := client.AppsV1().RESTClient()
	lw := cache.NewListWatchFromClient(restClient, r.RType, r.Namespace, fields.Everything())

	switch r.RType {
	case resource.Deployments:
		indexer, informer = cache.NewIndexerInformer(lw, &appsv1.Deployment{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
	case resource.Statefulsets:
		indexer, informer = cache.NewIndexerInformer(lw, &appsv1.StatefulSet{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
	case resource.Daemonsets:
		indexer, informer = cache.NewIndexerInformer(lw, &appsv1.DaemonSet{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
	}
	return
}

func (r *ResourceAndNamespace) CreateAllCoreV1IndexInformer(client *kubernetes.Clientset, worker queue.Queue, clusterName string) (indexerList []cache.Indexer, informerList []cache.Controller) {
	ctx := context.TODO()

	nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})

	if err != nil {
		klog.Error(ctx, err.Error())
	}

	for _, values := range nsList.Items {
		restClient := client.AppsV1().RESTClient()
		lw := cache.NewListWatchFromClient(restClient, r.RType, values.GetName(), fields.Everything())
		switch r.RType {
		case resource.Services:
			indexer, informer := cache.NewIndexerInformer(lw, &v1.Service{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		case resource.Pods:
			indexer, informer := cache.NewIndexerInformer(lw, &v1.Pod{}, 0, InitHandleFunc(resource.Pods, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		case resource.ConfigMaps:
			indexer, informer := cache.NewIndexerInformer(lw, &v1.ConfigMap{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		case resource.Secrets:
			indexer, informer := cache.NewIndexerInformer(lw, &v1.Secret{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		case resource.Events:
			indexer, informer := cache.NewIndexerInformer(lw, &v1.Event{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		}
	}
	return
}

func (r *ResourceAndNamespace) CreateAllAppsV1IndexInformer(client *kubernetes.Clientset, worker queue.Queue, clusterName string) (indexerList []cache.Indexer, informerList []cache.Controller) {
	ctx := context.TODO()

	nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})

	if err != nil {
		klog.Error(ctx, err.Error())
	}

	for _, values := range nsList.Items {
		restClient := client.AppsV1().RESTClient()
		lw := cache.NewListWatchFromClient(restClient, r.RType, values.GetName(), fields.Everything())
		switch r.RType {
		case resource.Deployments:
			indexer, informer := cache.NewIndexerInformer(lw, &appsv1.Deployment{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		case resource.Statefulsets:
			indexer, informer := cache.NewIndexerInformer(lw, &appsv1.StatefulSet{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		case resource.Daemonsets:
			indexer, informer := cache.NewIndexerInformer(lw, &appsv1.DaemonSet{}, 0, InitHandleFunc(resource.Services, clusterName, worker), cache.Indexers{})
			indexerList = append(indexerList, indexer)
			informerList = append(informerList, informer)
		}
	}

	return
}
