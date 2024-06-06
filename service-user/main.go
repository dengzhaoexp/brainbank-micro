package main

import (
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"strings"
	userConfig "user/config"
	"user/core"
	"user/pkg/utils/idmaker"
	userLogger "user/pkg/utils/logger"
	"user/repositry/cache"
	"user/repositry/dao"
	"user/userService"
)

func main() {
	// 先加载基础配置
	loading()

	// 读取etcd和该服务相关配置
	etcdConfig := userConfig.Config.Etcd
	etcdAddrs := strings.Join([]string{etcdConfig.Host, ":", etcdConfig.Port}, "")
	serviceConfig := userConfig.Config.Server
	serviceAddrs := strings.Join([]string{serviceConfig.Host, ":", serviceConfig.HttpPort}, "")

	// etcd配置
	etcdReg := etcd.NewRegistry(
		registry.Addrs(etcdAddrs)) // 指定注册中心

	// 创建微服务实例
	microService := micro.NewService(
		micro.Name(serviceConfig.Name),
		micro.Address(serviceAddrs),
		micro.Registry(etcdReg))

	//初始化结构命令行参数
	microService.Init()

	// 服务注册
	_ = userService.RegisterUserServiceHandler(microService.Server(), new(core.UserService))

	// 启动微服务
	_ = microService.Run()
}

func loading() {
	// 加载日志
	userLogger.InitLog()
	// 加载配置
	userConfig.InitConfig()
	// 加载mysql数据库
	dao.InitMysql()
	// 加载redis
	cache.InitRedis()
	// 初始化node
	idmaker.InitSnowflakeNode(1)
}
