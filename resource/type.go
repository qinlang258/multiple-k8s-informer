package resource

// 支持资源对象类型
const (
	All        = "all"
	Pods       = "pods"
	Services   = "services"
	ConfigMaps = "configmaps"
	Secrets    = "secrets"
	Events     = "events"
)

const (
	Deployments  = "deployments"
	Statefulsets = "statefulsets"
	Daemonsets   = "daemonsets"
)

// 事件类型
const (
	EventAdd    = "add"
	EventUpdate = "update"
	EventDelete = "delete"
)
