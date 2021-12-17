package hystrix

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	name := req.Service() + "." + req.Endpoint()

	// ------------ 自定义配置开始
	config := hystrix.CommandConfig{
		Timeout:                2000, // 用于设置超时时间，超过该时间没有返回响应，意味着请求失败；
		MaxConcurrentRequests:  10,   // 用于设置同一类型请求的最大并发量，达到最大并发量后，接下来的请求会被拒绝；
		RequestVolumeThreshold: 20,   //用于设置指定时间窗口内让断路器跳闸（开启）的最小请求数；
		SleepWindow:            5000, // 断路器跳闸后，在此时间段内，新的请求都会被拒绝；
		ErrorPercentThreshold:  50,   //请求失败百分比，如果超过这个百分比，则断路器跳闸。
	}

	hystrix.ConfigureCommand(name, config)
	// ------------ 自定义配置结束

	return hystrix.Do(name,
		func() error {
			return c.Client.Call(ctx, req, rsp, opts...)
		}, func(err error) error {
			if err != hystrix.ErrTimeout {
				return err
			}
			return nil
		},
	)
}

func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}
