package config

/**
 * @Author: zze
 * @Date: 2022/5/25 10:56
 * @Desc: Nacos 配置
 */
import (
	"demo/common/utils"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	zeroConf "github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zze326/zero-contrib/zrpc/registry/nacos"
	"log"
	"sync"
)

var (
	configClient config_client.IConfigClient
	nacosOnce    sync.Once
)

type Nacos struct {
	Addr        string
	Port        uint64
	Group       string
	DataID      string
	ExtDataIDs  []string `json:",optional"`
	NamespaceID string
}

func (conf *Nacos) InitConfigClient() (err error) {
	nacosOnce.Do(func() {
		configClient, err = clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig: &constant.ClientConfig{TimeoutMs: 5000, NamespaceId: conf.NamespaceID},
				ServerConfigs: []constant.ServerConfig{
					{IpAddr: conf.Addr, Port: conf.Port},
				},
			},
		)
	})
	return
}

func (conf *Nacos) GetConfig() (string, error) {
	var configMap = make(map[interface{}]interface{})
	mainConfig, err := configClient.GetConfig(vo.ConfigParam{DataId: conf.DataID, Group: conf.Group})
	if err != nil {
		return "", err
	}

	mainMap, err := utils.UnmarshalYamlToMap(mainConfig)
	if err != nil {
		return "", err
	}

	var extMap = make(map[interface{}]interface{})
	for _, dataID := range conf.ExtDataIDs {
		extConfig, err := configClient.GetConfig(vo.ConfigParam{DataId: dataID, Group: conf.Group})
		if err != nil {
			return "", err
		}

		tmpExtMap, err := utils.UnmarshalYamlToMap(extConfig)
		if err != nil {
			return "", err
		}

		extMap = utils.MergeMap(extMap, tmpExtMap)
	}

	configMap = utils.MergeMap(configMap, extMap)
	configMap = utils.MergeMap(configMap, mainMap)

	yamlString, err := utils.MarshalObjectToYamlString(configMap)
	if err != nil {
		return "", err
	}

	return yamlString, nil
}

func (conf *Nacos) Listen(onChange func(string, string, string, string)) error {
	return configClient.ListenConfig(vo.ConfigParam{
		DataId:   conf.DataID,
		Group:    conf.Group,
		OnChange: onChange,
	})
}

func (conf *Nacos) NewZrpcClient(serverName, clientName string) zrpc.Client {
	var target string
	target = fmt.Sprintf("nacos://%s:%d/%s?timeout=%s&namespace_id=%s&group_name=%s&app_name=%s", conf.Addr, conf.Port, serverName, "3s", conf.NamespaceID, conf.Group, clientName)
	return zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: target,
	})
}

func MustLoad(nacosConfigFilePath string, v interface{}) *Nacos {
	var (
		err    error
		config string
	)

	var nacosConfig Nacos
	zeroConf.MustLoad(nacosConfigFilePath, &nacosConfig, zeroConf.UseEnv())
	err = nacosConfig.InitConfigClient()
	if err != nil {
		log.Fatalf("init config client error: %v", err)
	}

	config, err = nacosConfig.GetConfig()
	if err != nil {
		log.Fatalf("get config error: %v", err)
	}

	err = zeroConf.LoadConfigFromYamlBytes([]byte(config), v)
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}
	return &nacosConfig
}

//func MustLoadWithListen(nacosConfigFilePath string, v interface{}) *Nacos {
//	var (
//		err         error
//		nacosConfig *Nacos
//	)
//	nacosConfig = MustLoad(nacosConfigFilePath, v)
//	err = nacosConfig.Listen(func(namespace, group, dataId, data string) {
//		err = zeroConf.LoadConfigFromYamlBytes([]byte(data), v)
//		if err != nil {
//			log.Printf("load config error: %v", err)
//		}
//	})
//	if err != nil {
//		log.Fatalf("listen config error: %v", err)
//	}
//	return nacosConfig
//}

func MustRegister(nacosConfig *Nacos, rpcConfig *zrpc.RpcServerConf) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nacosConfig.Addr, nacosConfig.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         nacosConfig.NamespaceID, // namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "info",
	}

	opts := nacos.NewNacosConfig(rpcConfig.Name, rpcConfig.ListenOn, sc, cc)
	err := nacos.RegisterService(opts)
	if err != nil {
		log.Fatalf("register service failed: %s", err)
	}
}
