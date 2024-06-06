package main

import (
	"chat/chatService"
	"chat/config"
	"chat/core"
	iLogger "chat/pkg/logger"
	"chat/repositry/dao"
	"chat/repositry/es"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"strings"
)

func main() {
	// 加载基础配置
	loading()

	// 加载etcd和该服务相关配置
	etcdConfig := config.Config.Etcd
	etcdAddrs := strings.Join([]string{etcdConfig.Host, ":", etcdConfig.Port}, "")
	serviceConfig := config.Config.Server
	serviceAddrs := strings.Join([]string{serviceConfig.Host, ":", serviceConfig.HttpPort}, "")

	// etcd配置
	etcdReg := etcd.NewRegistry(
		registry.Addrs(etcdAddrs)) // 指定注册中心

	// 创建对话微服务
	microService := micro.NewService(
		micro.Name(serviceConfig.Name),
		micro.Address(serviceAddrs),
		micro.Registry(etcdReg),
	)

	// 初始化结构命令行参数
	microService.Init()

	// 服务注册
	_ = chatService.RegisterConversationServiceHandler(microService.Server(), new(core.ConversionServer))

	// 启动微服务
	_ = microService.Run()
}

func loading() {
	// 加载日志配置
	iLogger.InitLog()
	// 加载基础配置
	config.InitConfig()
	// 加载mysql配置
	dao.InitMysql()
	// 加载es配置
	es.InitElastic()
}
