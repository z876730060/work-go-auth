package cloud

import (
	"context"
	"log/slog"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/z876730060/auth/internal/service"
	"github.com/z876730060/auth/pkg/feign"
)

func init() {
	RegisterManagerInstance.AddRegister(&Nacos{
		log: slog.Default().With("cloud", "nacos"),
	})
}

type Nacos struct {
	client  naming_client.INamingClient
	context context.Context
	cancel  func()
	log     *slog.Logger
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
		n.log.Error("nacos connect error", "err", err)
		return err
	}
	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          cfg.Application.IP,
		Port:        uint64(cfg.Application.Port),
		ServiceName: cfg.Application.Name,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		GroupName:   cfg.Cloud.Nacos.Group, // 默认值DEFAULT_GROUP
	})
	if err != nil {
		n.log.Error("nacos register error", "err", err)
		return err
	}
	n.client = client
	n.context, n.cancel = context.WithCancel(context.Background())
	go n.discovery(cfg)
	n.log.Info("nacos register success")
	return nil
}

func (n *Nacos) Unregister(cfg service.Config) error {
	if !cfg.Cloud.Nacos.Enable {
		return nil
	}

	n.log.Info("nacos unregister")
	n.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          cfg.Application.IP,
		Port:        uint64(cfg.Cloud.Nacos.Port),
		ServiceName: cfg.Application.Name,
		GroupName:   cfg.Cloud.Nacos.Group, // 默认值DEFAULT_GROUP
		Ephemeral:   true,
	})
	n.cancel()
	n.client.CloseClient()
	n.client = nil
	return nil
}

func (n *Nacos) discovery(cfg service.Config) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	n.feignDiscovery(cfg)
	for {
		select {
		case <-n.context.Done():
			n.log.Info("nacos discovery context done")
			return
		case <-ticker.C:
			n.feignDiscovery(cfg)
		}
	}
}

func (n *Nacos) feignDiscovery(cfg service.Config) {
	serviceSlice := make([]string, 0)
	pageNo := 1
	for {
		// 发现服务
		serviceList, err := n.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			NameSpace: cfg.Cloud.Nacos.Namespace, // 默认值public
			GroupName: cfg.Cloud.Nacos.Group,     // 默认值DEFAULT_GROUP
			PageNo:    uint32(pageNo),
			PageSize:  100,
		})
		if err != nil {
			n.log.Error("nacos discovery error", "err", err)
			continue
		}
		serviceSlice = append(serviceSlice, serviceList.Doms...)
		if serviceList.Count <= int64(pageNo*100) {
			break
		}
		pageNo++
	}
	feign.FeignClients.Clear()
	for _, name := range serviceSlice {
		if !feign.FeignClients.NeedDiscover(name) {
			continue
		}
		serviceInfo, err := n.client.GetService(vo.GetServiceParam{
			ServiceName: name,
			GroupName:   cfg.Cloud.Nacos.Group, // 默认值DEFAULT_GROUP
		})
		if err != nil {
			n.log.Error("nacos discovery error", "err", err)
			continue
		}
		if len(serviceInfo.Hosts) == 0 {
			continue
		} else {
			instance := feign.DiscoverInstance{
				Scheme:      "http",
				Ip:          serviceInfo.Hosts[0].Ip,
				Port:        serviceInfo.Hosts[0].Port,
				ServiceName: name,
			}
			feign.FeignClients.Register(instance)
			n.log.Debug("nacos discovery success", "instance", instance)
		}
	}
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
