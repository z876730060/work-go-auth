package job

import (
	"github.com/hyperjiang/xxljob"
	"github.com/z876730060/auth/internal/service"
)

var (
	executor      *xxljob.Executor
	jobHandlerMap = make(map[string]xxljob.JobHandler)
)

func Register(cfg service.Config) {
	// 注册定时任务
	if !cfg.CronJob.Enable {
		return
	}
	executor = xxljob.NewExecutor(
		xxljob.WithAppName(cfg.CronJob.AppName),
		xxljob.WithHost(cfg.CronJob.Host),
		xxljob.WithPort(cfg.CronJob.Port),
		xxljob.WithLogDir(cfg.CronJob.LogDir),
	)
	for name, job := range jobHandlerMap {
		executor.AddJobHandler(name, job)
	}
	executor.Start()
}

func Unregister() {
	if executor == nil {
		return
	}
	executor.Stop()
}
