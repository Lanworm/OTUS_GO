package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Lanworm/OTUS_GO/final_project/internal/logger"
	"github.com/Lanworm/OTUS_GO/final_project/internal/server/http"
	"github.com/Lanworm/OTUS_GO/final_project/internal/server/http/httphandler"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lanworm/OTUS_GO/final_project/internal/config"
	"github.com/Lanworm/OTUS_GO/final_project/pkg/shortcuts"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./build/local/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	configs, err := config.NewConfig(configFile)
	shortcuts.FatalIfErr(err)

	logg, err := logger.New(configs.Logger.Level, os.Stdout)
	shortcuts.FatalIfErr(err)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	httpServer := http.NewHTTPServer(logg, configs.Server.HTTP)
	handlerHTTP := httphandler.NewHandler(logg)
	httpServer.RegisterRoutes(handlerHTTP)
	go func() {
		logg.ServerLog(fmt.Sprintf("http server started on: http://%s", configs.Server.HTTP.GetFullAddress()))
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}
	}()
	logg.ServerLog("server is running...")

	<-ctx.Done()

	timeOutCtx, timeCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer timeCancel()

	if err := httpServer.Stop(timeOutCtx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
}
