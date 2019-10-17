package proxy

import (	
	"os"
	"log"
	"path/filepath"  
	"github.com/kardianos/service" //装包成操作系统的服务运行	
	"github.com/wooln/seagull2-grpc-webapi-proxy/gateway"	
  )
var logger service.Logger

//主方法
func Proxy(registerActions []gateway.RegisterAction)  {			
	
	config, err := GetConfigByDefaultPath()
	if err != nil {
		log.Fatal(err)
	}

	if(config.ProxyConfig.DocDir == ""){
		execPaht, _ := GetExePath()
		config.ProxyConfig.DocDir =  filepath.Join(execPaht, "doc")
	}

	osServiceConfig := config.OSServiceConfig

	svcConfig := &service.Config{
		Name:        osServiceConfig.Name,
		DisplayName: osServiceConfig.DisplayName,
		Description: osServiceConfig.Description,
	}

	prg := &program{
		RegisterActions: registerActions,
		ProxyConfig : config.ProxyConfig,		
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