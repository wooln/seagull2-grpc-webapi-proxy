package proxylib

import (
	"context"  // Use "golang.org/x/net/context" for Golang version <= 1.6
	"os"
	"flag"
	"net/http"
	"log"  
	"path/filepath"  
	"encoding/json"
	"errors"  
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"github.com/kardianos/service" //装包成操作系统的服务运行	
  )
var logger service.Logger

//配置文件相关
// Config is the runner app config structure.
type Config struct {
	Name, DisplayName, Description string
	WebAPIPort string 
	GrpcEndpointMapping map[string]string
	Stderr, Stdout string
}

func getConfigPath() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, execname := filepath.Split(fullexecpath)
	ext := filepath.Ext(execname)
	name := execname[:len(execname)-len(ext)]

	return filepath.Join(dir, name+".json"), nil
}

func getConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &Config{}

	r := json.NewDecoder(f)
	err = r.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

//go程序包装成运行的部分
type program struct{
	RegisterActions []RegisterAction 
}

//单个要注册的
type RegisterAction struct {
	//注册方法指引
	Action func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	//grpc配置映射的key，通过key照Endpint的数值
	EndpointKey string
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func (p *program) run() {
	// Do work here
	flag.Parse()
	defer glog.Flush() 
	if err := p.proxyGrpc(); err != nil {    
	  glog.Fatal(err)
	}
}

func (p *program) proxyGrpc() error {
	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	config, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	port := config.WebAPIPort

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	//循环把要注册的注册上。	
	for _, actionItem := range p.RegisterActions {
		endpint := config.GrpcEndpointMapping[actionItem.EndpointKey];
		if(endpint == ""){
			return errors.New("未找到key为"+actionItem.EndpointKey+"的GrpcEndpointMapping配置")
		}
		err = actionItem.Action(ctx, mux, endpint, opts)
		if err != nil {
			return err
		}
	}

	log.Printf("Greeter grpc gateway server listening on port " + port);
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(port, mux)
}

//主方法
func Proxy(registerActions []RegisterAction)  {		
	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	config, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	svcConfig := &service.Config{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		Description: config.Description,
	}

	prg := &program{
		RegisterActions: registerActions,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
		
	//根据参数决定什么行动, 通常使用install和uninstall即可
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}	
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}