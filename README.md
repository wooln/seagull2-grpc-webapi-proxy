# Summary
用于把grpc该代理成http webapi的类库

# Usage

1. import本类库
2. 准备和引入[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway/)生成的go文件
3. 准备配置文件
```
{
	"Name": "服务名称",
	"DisplayName": "服务显示名",
	"Description": "服务描述",
	
    "WebAPIPort":"如:8081, 代称成webapi的端口, 不要缺少前面的冒号",
	"GrpcEndpointMapping": {
		"Greeter" : "localhost:8080",
		"GreeterNew" : "localhost:8080"
	},
	
	"Stderr": "C:\\builder_err.log",
	"Stdout": "C:\\builder_out.log"
}
```
5. 调用Proxy, 传入[]RegisterAction类型的参数, EndpointKey为配置文件中的GrpcEndpointMapping自动的Key
```
package main

import (	
	proxyLib "github.com/wooln/seagull2-grpc-webapi-proxy"
	gw "Foo_Contracts"
)

func main()  {
	actions := []proxyLib.RegisterAction {
		proxyLib.RegisterAction{
			Action: gw.RegisterGreeterHandlerFromEndpoint,
			EndpointKey: "Greeter",
		},		
		proxyLib.RegisterAction{
			Action: gw.RegisterGreeterNewHandlerFromEndpoint,
			EndpointKey: "GreeterNew",
		},		
	}
	proxyLib.Proxy(actions)
}
```
