package main

import (
	"api/pb/social"
	"api/utils"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"social_service/config"
	"social_service/internal/dao"
	"social_service/internal/handler"
	"social_service/internal/registry"
	"social_service/logger"
	"syscall"
)

func main() {
	var cfn string
	// 0.从命令行获取可能的conf路径
	// api -conf="./conf/config_qa.yaml"
	// api -conf="./conf/config_online.yaml"
	flag.StringVar(&cfn, "conf", "./config/config.yaml", "指定配置文件路径")
	flag.Parse()
	// 1. 加载配置文件
	err := config.Init(cfn)
	if err != nil {
		panic(err) // 程序启动时加载配置文件失败直接退出
	}
	// 2. 初始化翻译
	_ = utils.InitTrans(utils.DefaultLocale)
	// 3. 初始化数据库连接池
	err = dao.InitMysql(config.Conf.MySQLConfig)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	// 4. 初始化Redis连接池
	dao.InitRedis()
	// 5. 注册服务
	//AutoRegister()
	// 5. 初始化Consul
	err = registry.Init(config.Conf.ConsulConfig.Addr)
	if err != nil {
		panic(err) // 程序启动时初始化注册中心失败直接退出
	}

	// 6. 启动服务
	StartService()

	// 7.注册服务到consul
	err = registry.Reg.RegisterService(config.Conf.Name, config.Conf.Address, nil)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	// 服务退出时要注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit // 正常会hang在此处
	// 退出时注销服务
	serviceId := fmt.Sprintf("%s-%s", config.Conf.Name, config.Conf.Address)
	registry.Reg.Deregister(serviceId)

}

// AutoRegister consul自动注册服务
func AutoRegister() {
	consulAddress := config.Conf.ConsulConfig.Addr
	if consulAddress == "" {
		logger.Log.Error("Consul address is not configured")
		return
	}

	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log.Error("Consul client failed:", err)
		return
	}

	serviceName := config.Conf.Name       // 你的服务名称
	serviceAddress := config.Conf.Address // 你的服务地址

	if serviceName == "" {
		logger.Log.Error("Service name is not configured")
		return
	}

	if serviceAddress == "" {
		logger.Log.Error("Service address is not configured")
		return
	}

	registration := &api.AgentServiceRegistration{
		Name:    serviceName,
		Address: serviceAddress,
		ID:      "user_service-" + serviceAddress, // 你的服务唯一 ID
	}

	logger.Log.Info("Registering service with name:", serviceName)
	logger.Log.Info("Service address:", serviceAddress)
	logger.Log.Info("Service ID:", registration.ID)

	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		logger.Log.Error("Consul service register failed:", err)
		return
	}

	listener, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		logger.Log.Error("Net listen failed:", err)
		return
	}

	server := grpc.NewServer()
	social.RegisterSocialServiceServer(server, handler.NewSocialService())

	logger.Log.Info("Starting gRPC server on address:", serviceAddress)
	err = server.Serve(listener)
	if err != nil {
		logger.Log.Error("gRPC server start failed:", err)
	}
}

// StartService 启动服务
func StartService() {
	lis, err := net.Listen("tcp", config.Conf.Address)
	if err != nil {
		logger.Log.Error("Net listen failed:", err)
		return
	}
	// 创建gRPC服务
	s := grpc.NewServer()
	// 健康检查
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	social.RegisterSocialServiceServer(s, handler.NewSocialService())
	// 商品服务注册RPC服务
	// 启动gRPC服务
	go func() {
		err = s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	logger.Log.Info("gRPC server start ", lis.Addr())
}
