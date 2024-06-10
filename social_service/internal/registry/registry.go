package registry

import "github.com/hashicorp/consul/api"

// 抽象做的好，后期可以很方便的切换不同的注册中心

// Register 自定义一个注册中心的抽象
type Register interface {
	// 注册
	RegisterService(serviceName string, serviceAddress string, tags []string) error
	// 服务发现
	ListService(serviceName string) (map[string]*api.AgentService, error)
	// 注销
	Deregister(serviceID string) error
}
