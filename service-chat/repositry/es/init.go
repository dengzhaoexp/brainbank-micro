package es

import (
	"chat/config"
	iLogger "chat/pkg/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"strings"
)

var _elastic *elasticsearch.Client

func InitElastic() {
	// 加载配置
	esConfig := config.Config.Elastic
	url := strings.Join([]string{"http://", esConfig.EsHost, ":", esConfig.EsPort}, "")

	// 创建es的配置
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		Password:  esConfig.EsPassword,
	}

	// 创建连接es的客户端
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		iLogger.LogrusObj.Error("Attempt to init elasticsearch fails:", err)
		panic(err)
	}

	// 尝试ping elasticsearch服务器，保证连接的正常
	_, err = es.Ping()
	if err != nil {
		iLogger.LogrusObj.Error("Attempt to ping elasticsearch server fails:", err)
		return
	}

	// 绑定
	_elastic = es
}
