package gateway

import (
	"context"
	"net/http"	
	"github.com/golang/glog"
	"log"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Endpoint describes a gRPC endpoint
type Endpoint struct {
	Network, Addr string
}

// Options is a set of options to be passed to Run
// 替换为了 GrpcWebApiProxyConfig
type Options struct {
	// Addr is the address to listen
	Addr string

	// GRPCServer defines an endpoint of a gRPC service
	GRPCServer Endpoint

	// SwaggerDir is a path to a directory from which the server
	// serves swagger specs.
	SwaggerDir string

	// Mux is a list of options to be passed to the grpc-gateway multiplexer
	Mux []gwruntime.ServeMuxOption
}

// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, registerActions []RegisterAction, opts GrpcWebApiProxyConfig) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	
	//根据配置决定是否启动时进行enpoint活性检查
	if(opts.CheckEndpoint){
		for key := range opts.GrpcEndpointMapping {			
			addr := opts.GrpcEndpointMapping[key]
			log.Println("检查地址...", addr, key)
			conn, err := dial(ctx, "tcp", addr)
			if err != nil {
				return err
			}
			log.Println("检查地址通过", addr, key)
			go func() {
				<-ctx.Done()
				if err := conn.Close(); err != nil {
					glog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
				}
			}()	
		}		
	}


	mux := http.NewServeMux()
	mux.HandleFunc("/doc/", swaggerServer(opts.DocDir))
	//mux.HandleFunc("/healthz", healthzServer(conn))

	gw, err := newGateway(ctx, registerActions, opts)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	addr := opts.WebAPIPort;
	s := &http.Server{
		Addr:    addr,
		Handler: allowCORS(mux),
	}
	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	glog.Infof("Starting listening at %s", addr)
	log.Println("Starting listening WebApi at", addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
		return err
	}
	return nil
}
