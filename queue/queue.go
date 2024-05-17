package queue

import (
	"errors"

	"time"

	"k8s.io/client-go/util/workqueue"
)

// QueueObject 入队对象
// 用来包装经由informer收到的资源对象

// resouce.QueueObject 结构体内容
type QueueObject struct {
	ClusterName  string      // 集群名称	集群名字
	Event        string      // 事件对象	ADD/DELETE/UPDATE
	ResourceType string      // 资源类型	pods/events/deployments
	Key          string      // <namespace>/<name>
	Obj          interface{} // runtime.Object	资源对象
	CreateAt     time.Time   // 创建时间，也可以记录更新次数 与 更新时间
}

type Queue interface {
	// Push 将监听到的资源放入queue中
	Push(QueueObject)
	// Pop 拿出队列
	Pop() (QueueObject, error)
	// ReQueue 重新放入队列，次数可配置
	ReQueue(QueueObject) error
	// Finish 完成入列操作
	Finish(QueueObject)
	// Close 关闭所有informer
	Close()
	// SetReMaxReQueueTime 设置最大重新入列次数
	SetReMaxReQueueTime(int)
}

type Wq struct {
	workqueue.RateLimitingInterface
	MaxReQueueTime int
}

var _ Queue = &Wq{}

func NewQueue(maxReQueueTime int) *Wq {
	return &Wq{
		workqueue.NewRateLimitingQueue(workqueue.DefaultItemBasedRateLimiter()),
		maxReQueueTime,
	}
}

func (q *Wq) SetReMaxReQueueTime(maxReQueueTime int) {
	q.MaxReQueueTime = maxReQueueTime
}

func (q *Wq) Push(obj QueueObject) {
	q.AddRateLimited(obj)
}

func (q *Wq) Pop() (QueueObject, error) {
	item, shutdown := q.Get()
	if shutdown {
		return QueueObject{}, errors.New("Controller has been stoped. ")
	}
	return item.(QueueObject), nil
}

func (q *Wq) ReQueue(obj QueueObject) (err error) {
	// 这里根据配置文件的尝试次数来推送至队列

	if q.NumRequeues(obj) < q.MaxReQueueTime {
		q.AddRateLimited(obj)
		return nil
	}

	q.Forget(obj)
	q.Done(obj)
	return errors.New("This object has been requeued for many times, but still fails. ")
}

func (q *Wq) Finish(obj QueueObject) {
	q.Forget(obj)
	q.Done(obj)
}

func (q *Wq) Close() {
	q.ShutDown()
}
