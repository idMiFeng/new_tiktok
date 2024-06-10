package registry

import (
	"errors"
	"fmt"
	"net"
	"social_service/logger"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type consul struct {
	client *api.Client
}

var Reg Register

// 确保某个结构体实现了对应的接口
var _ Register = (*consul)(nil)

// Init 连接至consul服务，初始化全局的consul对象
func Init(addr string) (err error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	Reg = &consul{c}
	return
}

// getOutboundIP 获取本机的出口IP
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// RegisterService 将gRPC服务注册到consul
func (c *consul) RegisterService(serviceName string, serviceAddress string, tags []string) error {
	if serviceName == "" {
		logger.Log.Error("Service name is not configured")
		return errors.New("service name is not configured")
	}

	if serviceAddress == "" {
		logger.Log.Error("Service address is not configured")
		return errors.New("service address is not configured")
	}

	outIp, err := getOutboundIP()
	if err != nil {
		logger.Log.Error("Failed to get outbound IP:", err)
		return err
	}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", outIp.String(), extractPort(serviceAddress)),
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "20s",
	}
	logger.Log.Info(outIp.String(), extractPort(serviceAddress))
	registration := &api.AgentServiceRegistration{
		Name:    serviceName,
		Address: outIp.String(),
		Port:    extractPort(serviceAddress),
		ID:      "social_service-" + serviceAddress,
		Tags:    tags,
		Check:   check,
	}

	logger.Log.Info("Registering service with name:", serviceName)
	logger.Log.Info("Service address:", serviceAddress)
	logger.Log.Info("Service ID:", registration.ID)

	err = c.client.Agent().ServiceRegister(registration)
	if err != nil {
		logger.Log.Error("Consul service register failed:", err)
		return err
	}

	return nil
}

// ListService 服务发现
func (c *consul) ListService(serviceName string) (map[string]*api.AgentService, error) {
	return c.client.Agent().ServicesWithFilter(fmt.Sprintf("Service==`%s`", serviceName))
}

// Deregister 注销服务
func (c *consul) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

// extractPort 从地址中提取端口
func extractPort(address string) int {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		logger.Log.Error("Invalid service address:", err)
		return 0
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		logger.Log.Error("Invalid port:", err)
		return 0
	}
	return portInt
}
