package client

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-client/conf"
	"github.com/waittttting/cRPC-common/http"
	"github.com/waittttting/cRPC-common/model"
	"github.com/waittttting/cRPC-common/tcp"
)

type RpcClient struct {
	config *conf.LocalConf
	cc *cloudConfig
	ControlCenterConn *tcp.Connection
	// 业务服务长连接 key =
	ServerClient map[string][]*tcp.Connection
}

func NewRpcClient(config *conf.LocalConf) *RpcClient {
	return &RpcClient{
		config: config,
		ServerClient: make(map[string][]*tcp.Connection),
	}
}

func (rc *RpcClient) Start() {
	// 发送 http 请求获取 config file
	cc := rc.getServerConfig()
	rc.cc = cc
	// 连接控制中心
	conn := tcp.CreateSocket(cc.ControlCenterHost)
	rc.ControlCenterConn = conn
	// TODO: 鉴权？内部服务是否需要鉴权
	// 发送注册信息
	p := tcp.NewRegisterMessage(rc.config.Client.ServerName, rc.config.Client.ServerVersion)
	// 注册到控制中心
	err := conn.Send(p)
	if err != nil {
		logrus.Fatalf("send server message error [%v]", err)
	}
	// 获取要交互的服务以及相关配置信息
}

// 存储在配置中心的云端配置
type cloudConfig struct {
	ControlCenterHost string `json:"control_center_host"`
}

func (rc *RpcClient) getServerConfig() *cloudConfig {

	params := map[string]string{
		"server_name" : rc.config.Client.ServerName,
		"server_version" : rc.config.Client.ServerVersion}
	result, err := http.Post(rc.config.Client.ConfigCenterHost + "/get/config", params)
	if err != nil {
		logrus.Fatalf("load server config err [%v]", err)
	}
	var cc cloudConfig
	var sc model.ServerConfig
	err = json.Unmarshal([]byte(result), &sc)
	if err != nil {
		logrus.Fatalf("marshal to cloudConfig err [%v]", err)
	}
	err = json.Unmarshal([]byte(sc.Config), &cc)
	if err != nil {
		logrus.Fatalf("marshal to model.ServerConfig err [%v]", err)
	}
	return &cc
}
