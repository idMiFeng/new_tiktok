package video

import (
	"api/global"
	"api/pb/video"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"video_service/config"
	"video_service/logger"
)

// InitRpcVideoClient 初始化用户服务的RPC客户端
func InitRpcVideoClient() {
	consulAddress := config.Conf.ConsulConfig.Addr // 替换为您的 Consul 地址
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log.Error("Failed to create Consul client:", err)
	}

	// 使用 Consul 客户端获取服务实例
	serviceEntries, _, err := consulClient.Health().Service("video_service", "", true, nil)
	if err != nil {
		logger.Log.Error("Failed to query service from Consul:", err)
	}

	// 选择一个服务实例并创建 gRPC 连接
	var target string
	if len(serviceEntries) > 0 {
		// 选择第一个实例
		service := serviceEntries[0]
		target = fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	} else {
		logger.Log.Error("No available service instances found")
	}

	videoConn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal(err)
	}
	videoClient := video.NewVideoServiceClient(videoConn)
	global.VideoClient = videoClient
}
