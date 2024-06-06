package main

import (
	apiConfig "api-gateway/config"
	apiLogger "api-gateway/pkg/utils/logger"
	"api-gateway/service/chatService"
	"api-gateway/service/fileService"
	"api-gateway/service/userService"
	"api-gateway/weblib"
	"api-gateway/wrapper"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"strings"
	"time"
)

func main() {
	// 加载配置文件
	loading()

	// 读取etcd和该服务相关配置
	etcdConfig := apiConfig.Config.Etcd
	etcdAddrs := strings.Join([]string{etcdConfig.Host, ":", etcdConfig.Port}, "")
	serviceConfig := apiConfig.Config.Server
	serviceAddrs := strings.Join([]string{serviceConfig.Host, ":", serviceConfig.HttpPort}, "")

	// etcd的配置
	etcdReg := etcd.NewRegistry(
		registry.Addrs(etcdAddrs))

	// 用户服务及其调用实例
	userMicroService := micro.NewService(
		micro.Name("userService.Client"),
		micro.WrapClient(wrapper.NewUserWrapper))
	service01 := userService.NewUserService("rpcUserService", userMicroService.Client())

	// 文件服务及其调用实例
	fileMicroService := micro.NewService(
		micro.Name("fileService.Client"),
		micro.WrapClient(wrapper.NewFileWrapper))
	service02 := fileService.NewFileService("rpcFileService", fileMicroService.Client())

	// chat服务及其调用实例
	chatMicroService := micro.NewService(
		micro.Name("chatService.Client"),
		micro.WrapClient(wrapper.NewChatWrapper),
		micro.Client(client.NewClient(client.RequestTimeout(120*time.Second), client.Retries(0))),
	)
	service03 := chatService.NewConversationService("rpcChatService", chatMicroService.Client())

	// 创建微服务实例，使用gin暴露http接口并注册到etcd中
	server := web.NewService(
		web.Name(serviceConfig.Name),
		web.Address(serviceAddrs),
		// 将服务调用实例使用gin处理
		web.Handler(weblib.NewRouter(service01, service02, service03)),
		web.Registry(etcdReg),
		web.Metadata(map[string]string{"proto": "http"}))

	// 初始化
	_ = server.Init()
	// 运行
	_ = server.Run()
}

func loading() {
	// 实例化日志
	apiLogger.InitLog()
	// 加载配置
	apiConfig.InitConfig()
}
