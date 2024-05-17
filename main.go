package main

import (
	"fmt"
	"multiple-k8s-informer/config"
	"multiple-k8s-informer/controller"
	"multiple-k8s-informer/queue"
	"multiple-k8s-informer/resource"
	"multiple-k8s-informer/store"
	"time"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

func NewMultiClusterInformerFromConfig(path string) (controller.MultiClusterInformer, error) {
	sysConfig, err := config.LoadConfig(path)
	if err != nil {
		klog.Error("load config error: ", err)
		return nil, err
	}

	return NewMultiClusterInformer(sysConfig.MaxReQueueTime, sysConfig.Clusters)

}

func NewMultiClusterInformer(maxReQueueTime int, clusters []controller.Cluster) (controller.MultiClusterInformer, error) {
	core := &controller.Controller{
		Queue:  queue.NewQueue(maxReQueueTime),
		StopCh: make(chan struct{}, 1),
	}

	store := make(store.MapIndexers)
	informers := make(controller.InformerList, 0)

	for _, cluster := range clusters {
		//对每个集群 初始化一个 clientset
		client, err := cluster.NewClient()
		if err != nil {
			return nil, err
		}

		for _, r := range cluster.List {
			if r.Namespace == resource.All {
				//当 namespace为 all的时候单独处理
				var indexerListRes []cache.Indexer
				var informerListRes []cache.Controller

				switch r.RType {
				case resource.Deployments:
					indexerListRes, informerListRes = r.CreateAllAppsV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Statefulsets:
					indexerListRes, informerListRes = r.CreateAllAppsV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Daemonsets:
					indexerListRes, informerListRes = r.CreateAllAppsV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Pods:
					indexerListRes, informerListRes = r.CreateAllCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.ConfigMaps:
					indexerListRes, informerListRes = r.CreateAllCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Secrets:
					indexerListRes, informerListRes = r.CreateAllCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Services:
					indexerListRes, informerListRes = r.CreateAllCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Events:
					indexerListRes, informerListRes = r.CreateAllCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				}

				for k, v := range indexerListRes {
					if v != nil || informerListRes[k] != nil {
						store[r.RType] = append(store[r.RType], v)
						informers = append(informers, informerListRes[k])
					}

				}
			} else {
				var indexer cache.Indexer
				var informer cache.Controller

				switch r.RType {
				case resource.Deployments:
					indexer, informer = r.CreateAppsV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Statefulsets:
					indexer, informer = r.CreateAppsV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Daemonsets:
					indexer, informer = r.CreateAppsV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Pods:
					indexer, informer = r.CreateCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.ConfigMaps:
					indexer, informer = r.CreateCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Secrets:
					indexer, informer = r.CreateCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Services:
					indexer, informer = r.CreateCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				case resource.Events:
					indexer, informer = r.CreateCoreV1IndexInformer(client, core.Queue, cluster.ClusterName)
				}

				// 放入 list中
				store[r.RType] = append(store[r.RType], indexer)
				informers = append(informers, informer)
			}

		}
	}

	core.Informers = informers
	core.Store = store

	return core, nil
}

// process execute your own logic
func process(obj queue.QueueObject) error {

	if obj.Event == resource.EventAdd {
		fmt.Println("目前监听到事件为add的资源对象", obj.ResourceType)
	}

	fmt.Println(time.Now(), obj.Event, obj.ResourceType, obj.Key, obj.ClusterName)
	return nil
}

func main() {
	r, err := NewMultiClusterInformerFromConfig("./config.yaml")
	if err != nil {
		klog.Fatal("multi cluster informer err: ", err)
	}

	// 2. add handler
	r.AddEventHandler(func(object queue.QueueObject) error {
		// only the add event
		if object.Event == resource.EventAdd {
			fmt.Println("目前监听到事件为add的资源对象", object.ResourceType)
		}
		//fmt.Println("目前监听到的资源对象", obj.ResourceType, obj.Event)

		if object.ClusterName == "cluster1" {
			fmt.Println("目前监听到集群为cluster1的资源对象", object.ResourceType)
		}

		if object.ClusterName == "cluster2" {
			fmt.Println("目前监听到集群为cluster2的资源对象", object.ResourceType)
		}

		fmt.Println(time.Now(), object.Event, object.ResourceType, object.Key, object.ClusterName)
		return nil
	})

	// 3. run informer
	go r.Run()
	defer r.Stop()

	// 4. Continuously remove resource objects from the queue
	for {
		obj, _ := r.Pop()
		if err = r.HandleObject(obj); err != nil {
			_ = r.ReQueue(obj) // reQueue
		} else {
			// t's done
			r.Finish(obj)
		}

		// method two：Handle it yourself
		//if err = process(obj); err != nil {
		// _ = r.ReQueue(obj) // 重新入列
		//} else { // 完成就结束
		// r.Finish(obj)
		//}
	}
}
