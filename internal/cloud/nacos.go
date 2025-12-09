package cloud

import (
	"log/slog"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/z876730060/auth/internal/service"
)

func init() {
	RegisterManagerInstance.AddRegister(&Nacos{})
}

type Nacos struct {
	client naming_client.INamingClient
}

func (n *Nacos) Register(cfg service.Config) error {
	if !cfg.Cloud.Nacos.Enable {
		return nil
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"clientConfig":  getClientConfig(cfg),
		"serverConfigs": getServerConfigs(cfg),
	})
	if err != nil {
		slog.Error("nacos connect error", "err", err)
		return err
	}
	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          cfg.Application.IP,
		Port:        uint64(cfg.Cloud.Nacos.Port),
		ServiceName: cfg.Application.Name,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		GroupName:   cfg.Cloud.Nacos.Group, // 默认值DEFAULT_GROUP
	})
	if err != nil {
		slog.Error("nacos register error", "err", err)
		return err
	}
	n.client = client
	slog.Info("nacos register success")
	return nil
}

func (n *Nacos) Unregister(cfg service.Config) error {
	if !cfg.Cloud.Nacos.Enable {
		return nil
	}

	slog.Info("nacos unregister")
	n.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          cfg.Application.IP,
		Port:        uint64(cfg.Cloud.Nacos.Port),
		ServiceName: cfg.Application.Name,
		GroupName:   cfg.Cloud.Nacos.Group, // 默认值DEFAULT_GROUP
		Ephemeral:   true,
	})
	n.client.CloseClient()
	n.client = nil
	return nil
}

func getClientConfig(cfg service.Config) constant.ClientConfig {
	return constant.ClientConfig{
		NamespaceId:         cfg.Cloud.Nacos.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
	}
}

func getServerConfigs(cfg service.Config) []constant.ServerConfig {
	return []constant.ServerConfig{
		{
			IpAddr:      cfg.Cloud.Nacos.Ip,
			ContextPath: "/nacos",
			Port:        uint64(cfg.Cloud.Nacos.Port),
			Scheme:      "http",
		},
	}
}
