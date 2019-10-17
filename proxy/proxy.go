package proxy

import (
	"context"  // Use "golang.org/x/net/context" for Golang version <= 1.6
	"net/http"
	"log"  
	"errors"  
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
  )

func ProxyGrpc2WebApi(registerActions []RegisterAction, config GrpcWebApiProxyConfig) error {	
	port := config.WebAPIPort

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	//循环把要注册的注册上。	
	for _, actionItem := range registerActions {
		endpint := config.GrpcEndpointMapping[actionItem.EndpointKey];
		if(endpint == ""){
			return errors.New("未找到key为"+actionItem.EndpointKey+"的GrpcEndpointMapping配置")
		}
		err := actionItem.Action(ctx, mux, endpint, opts)
		if err != nil {
			return err
		}
	}

	log.Printf("Greeter grpc gateway server listening on port " + port);
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(port, mux)
}