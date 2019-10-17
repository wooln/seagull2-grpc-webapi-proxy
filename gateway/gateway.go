package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
	"errors"  
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(ctx context.Context, registerActions []RegisterAction, opts  GrpcWebApiProxyConfig) (http.Handler, error) {

	mux := gwruntime.NewServeMux(opts.Mux...)
	option := []grpc.DialOption{grpc.WithInsecure()}

	for _, actionItem := range registerActions {
		
		endpint := opts.GrpcEndpointMapping[actionItem.EndpointKey];
		if(endpint == ""){
			return nil, errors.New("未找到key为"+actionItem.EndpointKey+"的GrpcEndpointMapping配置")
		}
		action :=actionItem.Action
		if err := action(ctx, mux, endpint, option); err != nil {			
			return nil, err
		}
    }

	// for _, f := range []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error {
	// 	// // examplepb.RegisterEchoServiceHandler,
	// 	// // examplepb.RegisterStreamServiceHandler,
	// 	// // examplepb.RegisterABitOfEverythingServiceHandler,
	// 	// // examplepb.RegisterFlowCombinationHandler,
	// 	// // examplepb.RegisterNonStandardServiceHandler,
	// 	// // examplepb.RegisterResponseBodyServiceHandler,
	// } {
	// 	if err := f(ctx, mux, conn); err != nil {
	// 		return nil, err
	// 	}
	// }

	return mux, nil
}

func dial(ctx context.Context, network, addr string) (*grpc.ClientConn, error) {
	switch network {
	case "tcp":
		return dialTCP(ctx, addr)
	case "unix":
		return dialUnix(ctx, addr)
	default:
		return nil, fmt.Errorf("unsupported network type %q", network)
	}
}

// dialTCP creates a client connection via TCP.
// "addr" must be a valid TCP address with a port number.
func dialTCP(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, grpc.WithInsecure())
}

// dialUnix creates a client connection via a unix domain socket.
// "addr" must be a valid path to the socket.
func dialUnix(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	d := func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("unix", addr, timeout)
	}
	return grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithDialer(d))
}
