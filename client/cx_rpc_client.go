package client

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-client/conf"
	"github.com/waittttting/cRPC-common"
)

type RpcClient struct {
	config *conf.LocalConf
}

func NewRpcClient(config *conf.LocalConf) *RpcClient {
	return &RpcClient{
		config: config,
	}
}

func (rc *RpcClient) Start() {
	// 发送 http 请求获取 config file
	cc := rc.getServerConfig()
	print(cc.ControlCenterHost)
	// 连接 控制中心
}

func (rc *RpcClient) getServerConfig() *cloudConfig {

	params := map[string]string{
		"server_name" : rc.config.Client.ServerName,
		"server_version" : rc.config.Client.ServerVersion}
	result, err := common.Post(rc.config.Client.ConfigCenterHost + "/get/config", params)
	if err != nil {
		logrus.Fatalf("load server config err [%v]", err)
	}
	var cc cloudConfig
	err = json.Unmarshal([]byte(result), &cc)
	if err != nil {
		logrus.Fatalf("marshal to cloudConfig err [%v]", err)
	}
	return &cc
}

// 存储在配置中心的配置
type cloudConfig struct {
	ControlCenterHost string `json:"control_center_host"`
	ControlCenterPort string `json:"control_center_port"`
}