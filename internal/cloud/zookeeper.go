package cloud

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/z876730060/auth/internal/service"
	zkregister "github.com/z876730060/work-zkRegister-cloud"
)

func init() {
	RegisterManagerInstance.AddRegister(NewZookeeper())
}

type Zookeeper struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewZookeeper() *Zookeeper {
	ctx, cancel := context.WithCancel(context.Background())
	return &Zookeeper{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (z *Zookeeper) Register(cfg service.Config) error {
	if !cfg.Cloud.Zookeeper.Enable {
		return nil
	}

	zkConn, _, err := zk.Connect([]string{cfg.Cloud.Zookeeper.Ip + ":" + strconv.FormatUint(cfg.Cloud.Zookeeper.Port, 10)}, time.Second*5)
	if err != nil {
		slog.Error("connect zookeeper failed", "err", err)
		return err
	}

	instance := zkregister.ServiceInfo{
		Name:    cfg.Application.Name,
		Address: cfg.Application.IP,
		Port:    cfg.Application.Port,
	}

	zkregister.RegisterZK(z.ctx, zkConn, slog.Default(), instance)
	slog.Info("zookeeper register success")
	return nil
}

func (z *Zookeeper) Unregister(cfg service.Config) error {
	if !cfg.Cloud.Zookeeper.Enable {
		return nil
	}

	slog.Info("unregister zookeeper")
	z.cancel()
	return nil
}
