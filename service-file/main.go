package main

import (
	fileConfig "file/config"
	"file/core"
	"file/fileService"
	"file/pkg/utils/idmaker"
	fileLogger "file/pkg/utils/logger"
	fileDao "file/repositry/dao"
	"file/repositry/mq"
	"file/repositry/oss"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"strings"
)

func main() {
	// 加载基础配置
	loading()

	// 读取etcd和该服务相关配置
	etcdConfig := fileConfig.Config.Etcd
	etcdAddrs := strings.Join([]string{etcdConfig.Host, ":", etcdConfig.Port}, "")
	serviceConfig := fileConfig.Config.Server
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
	_ = fileService.RegisterFileServiceHandler(microService.Server(), new(core.FileService))

	// 启动微服务
	_ = microService.Run()
}

func loading() {
	// 加载日志
	fileLogger.InitLog()
	// 加载配置文件
	fileConfig.InitConfig()
	// 加载mysql数据库
	fileDao.InitMysql()
	// 加载连接oss
	oss.InitMinioClient()
	// 加载消息队列
	mq.InitRabbitMQ()
	// 加载idmaker
	idmaker.InitSnowflakeNode(2)

}
