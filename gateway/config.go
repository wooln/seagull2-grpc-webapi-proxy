package gateway

import (
	"context"	
	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
  )

type GrpcWebApiProxyConfig struct {
	WebAPIPort string 
	GrpcEndpointMapping map[string]string
	CheckEndpoint bool
	DocDir string
	Mux []runtime.ServeMuxOption
}

//单个要注册的
type RegisterAction struct {
	//注册方法指引
	Action func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	//grpc配置映射的key，通过key照Endpint的数值
	EndpointKey string
}