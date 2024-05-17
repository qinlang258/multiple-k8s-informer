package store

import (
	"multiple-k8s-informer/resource"

	"k8s.io/client-go/tools/cache"
)

// Store 本地缓存接口
type Store interface {
	// List 列出所有资源对象
	List(string) []interface{}
	// ListKeys 列出所有资源对象的key
	ListKeys(string) []string
	// GetByKey 输入特定key，返回资源对象
	GetByKey(r string, key string) (items []interface{}, exists bool)
}

// 这个store是 informer的 indexers结构体的Store内容，只提取了List，ListKeys，GetByKey这3个方法，并根据我们想要达到的目的修改了代码
// 因为我们这个是对多个集群进行监控，因此这里多了一层 切片的选择,这里的 string是对 资源类型的map，mapIndexer[Pods]，mapIndexer[Services]例如这样的
// type Store interface {
// 	List() []interface{}
// 	ListKeys() []string
// 	GetByKey(key string) (item interface{}, exists bool, err error)
// }

type MapIndexers map[string][]cache.Indexer

var _ Store = MapIndexers{}

func (mapIndexers MapIndexers) List(resourceName string) (items []interface{}) {
	switch resourceName {
	case resource.All:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.Services:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.Pods:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.ConfigMaps:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.Events:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.Deployments:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.Statefulsets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.Daemonsets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	case resource.Secrets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.List()...)
			}
		}
	}
	return
}

func (mapIndexers MapIndexers) ListKeys(resourceName string) (items []string) {
	switch resourceName {
	case resource.All:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.Services:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.Pods:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.ConfigMaps:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.Events:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.Deployments:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.Statefulsets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.Daemonsets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	case resource.Secrets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				items = append(items, indexer.ListKeys()...)
			}
		}
	}
	return
}

func (mapIndexers MapIndexers) GetByKey(resourceName, key string) (items []interface{}, ok bool) {
	switch resourceName {
	case resource.All:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	case resource.Services:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	case resource.Pods:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	case resource.Deployments:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	case resource.Statefulsets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	case resource.ConfigMaps:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	case resource.Daemonsets:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	case resource.Events:
		for _, mapIndexer := range mapIndexers {
			for _, indexer := range mapIndexer {
				item, exists, err := indexer.GetByKey(key)
				if err != nil {
					continue
				}
				if exists {
					ok = true
					items = append(items, item)
				}
			}
		}
	}

	return

}
