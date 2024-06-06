package wrapper

import (
	"context"
	"go-micro.dev/v4/client"
	"time"
)

type ChatWrapper struct {
	client.Client
}

func (cw *ChatWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
	//cmdName := req.Service() + "." + req.Endpoint()
	//config := hystrix.CommandConfig{
	//	Timeout:                30000,
	//	RequestVolumeThreshold: 20,   // 熔断器请求阈值，意思是有20个请求才能进行错误百分比计算
	//	ErrorPercentThreshold:  50,   // 错误百分比，当错误超过百分比时，直接进行降级处理，直至熔断器再次开启，默认为50%
	//	SleepWindow:            5000, // 熔断器再次检测开启经过的时间，单位毫秒ms，默认为5s
	//}
	//hystrix.ConfigureCommand(cmdName, config)
	//return hystrix.Do(cmdName, func() error {
	//
	//	return cw.Client.Call(ctx, req, resp, opts...)
	//}, func(err error) error {
	//	return err
	//})

	timeoutCTx, _ := context.WithTimeout(ctx, 120*time.Second)
	return cw.Client.Call(timeoutCTx, req, resp, client.WithRetries(0), client.WithRequestTimeout(120*time.Second))
}

func NewChatWrapper(c client.Client) client.Client {
	return &ChatWrapper{Client: c}
}
