package controller

import (
	"errors"
	"multiple-k8s-informer/queue"
	"multiple-k8s-informer/resource"
	"multiple-k8s-informer/store"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

// Cluster 集群对象所需的信息
type Cluster struct {
	List        []ResourceAndNamespace `json:"list" yaml:"list"`
	ConfigPath  string                 `json:"configPath" yaml:"configPath"` // kube config文件
	Insecure    bool                   `json:"insecure" yaml:"insecure"`     // 是否跳过证书认证
	ClusterName string                 `json:"clusterName" yaml:"clusterName"`
}

// MultiClusterInformer 多集群informer的接口对象
type MultiClusterInformer interface {
	// Run 执行多集群的informer的方法
	Run()
	// Stop 停止informer
	Stop()
	// AddEventHandler 加入回调handler
	AddEventHandler(handler HandleFunc)
	// HandleObject 调用handler处理资源对象
	HandleObject(object queue.QueueObject) error
	// Queue 队列接口对象
	queue.Queue
	// Store 本地缓存接口对象
	store.Store
}

var _ MultiClusterInformer = &Controller{}

type InformerList []cache.Controller

func (informerList InformerList) Run(stopCh chan struct{}) {
	for _, informer := range informerList {
		go informer.Run(stopCh)

		if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
			return
		}

	}
}

func (c *Cluster) NewClient() (*kubernetes.Clientset, error) {

	if c.ConfigPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", c.ConfigPath)
		if err != nil {
			return nil, err
		}
		config.Insecure = c.Insecure
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		return clientset, nil
	}

	return nil, errors.New("无法找到集群client端")
}

type Controller struct {
	clients []*kubernetes.Clientset
	queue.Queue
	store.Store
	StopCh     chan struct{}
	Informers  InformerList
	HandleFunc HandleFunc
}

func (c *Controller) Run() {
	klog.Info("run controller...")
	defer c.Queue.Close()
	c.Informers.Run(c.StopCh)
	<-c.StopCh

}

// HandleObject 自定义回调方法
func (c *Controller) HandleObject(obj queue.QueueObject) error {
	if c.HandleFunc != nil {
		err := c.HandleFunc(obj)
		return err
	}
	return nil
}

// Stop 停止
func (c *Controller) Stop() {
	c.StopCh <- struct{}{}
}

type HandleFunc func(object queue.QueueObject) error

func (c *Controller) AddEventHandler(handler HandleFunc) {
	c.HandleFunc = handler
}

func InitHandleFunc(resourceName, clusterName string, worker queue.Queue) cache.ResourceEventHandlerFuncs {
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queueObj := queue.QueueObject{ClusterName: clusterName, ResourceType: resourceName, Event: resource.EventAdd, Obj: key, CreateAt: time.Now()}
				worker.Push(queueObj)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queueObj := queue.QueueObject{ClusterName: clusterName, ResourceType: resourceName, Event: resource.EventUpdate, Obj: key, CreateAt: time.Now()}
				worker.Push(queueObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queueObj := queue.QueueObject{ClusterName: clusterName, ResourceType: resourceName, Event: resource.EventDelete, Obj: key, CreateAt: time.Now()}
				worker.Push(queueObj)
			}
		},
	}

	return handler
}
