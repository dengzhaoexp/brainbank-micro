package wrapper

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"go-micro.dev/v4/client"
)

type userWrapper struct {
	client.Client
}

func (w *userWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
	cmdName := req.Service() + "." + req.Endpoint()
	config := hystrix.CommandConfig{
		Timeout:                30000,
		RequestVolumeThreshold: 20,   // 熔断器请求阈值，意思是有20个请求才能进行错误百分比计算
		ErrorPercentThreshold:  50,   // 错误百分比，当错误超过百分比时，直接进行降级处理，直至熔断器再次开启，默认50%
		SleepWindow:            5000, // 过多长时间，熔断器再次检测是否开启，单位毫秒ms（默认5秒）
	}
	hystrix.ConfigureCommand(cmdName, config)

	return hystrix.Do(cmdName, func() error {
		return w.Client.Call(ctx, req, resp)
	}, func(err error) error {
		return err
	})
}

func NewUserWrapper(c client.Client) client.Client {
	return &userWrapper{c}
}
