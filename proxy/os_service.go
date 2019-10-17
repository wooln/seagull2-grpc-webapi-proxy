package proxy

import (	
	"flag"		
	"github.com/golang/glog"
	"github.com/kardianos/service" //装包成操作系统的服务运行	
  )

//go程序包装成运行的部分
type program struct{
	ProxyConfig GrpcWebApiProxyConfig
	RegisterActions []RegisterAction 
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
	return ProxyGrpc2WebApi(p.RegisterActions, p.ProxyConfig)
}