package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cuijxin/my-n9e/config"
	"github.com/cuijxin/my-n9e/models"
	"github.com/cuijxin/my-n9e/pkg/i18n"
	"github.com/cuijxin/my-n9e/pkg/ilog"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"

	prom_runtime "github.com/prometheus/prometheus/pkg/runtime"
)

var version = "not specified"

var (
	vers *bool
	help *bool
)

func init() {
	vers = flag.Bool("v", false, "display the version.")
	help = flag.Bool("h", false, "print this help.")
	flag.Parse()

	if *vers {
		fmt.Println("version:", version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	runner.Init()
	fmt.Println("runner.cwd:", runner.Cwd)
	fmt.Println("runner.hostname:", runner.Hostname)
	fmt.Println("fd_limits", prom_runtime.FdLimits())
	fmt.Println("vm_limits", prom_runtime.VmLimits())
}

func main() {
	parseConf()

	ilog.Init(config.Config.Logger)
	i18n.Init(config.Config.I18N)

	models.InitMySQL(config.Config.MySQL)

	_, cancelFunc := context.WithCancel(context.Background())

	endingProc(cancelFunc)
}

func parseConf() {
	if err := config.Parse(); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	}
}

func endingProc(cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-c
	fmt.Printf("stop signal caught, stopping... pid=%d\n", os.Getpid())

	// 执行清理工作
	// backend.DatasourceCleanUp()
	cancelFunc()
	logger.Close()

	fmt.Println("process stopped successfully")
}
