package proxy

import (
	"os"
	"path/filepath"  
	"encoding/json"
	"github.com/wooln/seagull2-grpc-webapi-proxy/gateway"
  )

//配置文件相关
// Config is the runner app config structure.
type GrpcProxyOSServiceConfig struct {
	OSServiceConfig OSServiceConfig	
	ProxyConfig gateway.GrpcWebApiProxyConfig	
}

type OSServiceConfig struct {
	Name, DisplayName, Description string	
	Stderr, Stdout string
}


//获取可执行文件的目录，因为以服务运行时候需要计算的
func GetExePath() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, _ := filepath.Split(fullexecpath)	
	return dir, nil
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

func getConfig(path string) (*GrpcProxyOSServiceConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &GrpcProxyOSServiceConfig{}

	r := json.NewDecoder(f)
	err = r.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func GetConfigByDefaultPath() (*GrpcProxyOSServiceConfig, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &GrpcProxyOSServiceConfig{}

	r := json.NewDecoder(f)
	err = r.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
