package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/app"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/grpc/grpchandler"
	internalhttp "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http/httphandler"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/pkg/shortcuts"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/pkg/storageconf"
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

	ctx := context.Background()

	appConfig, err := config.NewConfig(configFile)
	shortcuts.FatalIfErr(err)

	logg, err := logger.New(appConfig.Logger.Level, os.Stdout)
	shortcuts.FatalIfErr(err)

	eventStorage, _ := storageconf.NewStorage(ctx, *appConfig, logg)

	_ = app.New(logg, eventStorage)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	httpServer := internalhttp.NewHTTPServer(logg, appConfig.Server.HTTP)
	evtService := service.NewEventService(logg, eventStorage)

	handlerHTTP := httphandler.NewHandler(logg, evtService)
	httpServer.RegisterRoutes(handlerHTTP)
	go func() {
		logg.ServerLog(fmt.Sprintf("http server started on: %s", appConfig.Server.HTTP.GetFullAddress()))
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}
	}()

	handlerGrpc := grpchandler.NewHandler(logg, evtService)
	grpcServer := internalgrpc.NewRPCServer(
		handlerGrpc,
		logg,
		appConfig.Server.GRPC,
	)
	go func() {
		logg.ServerLog(fmt.Sprintf("grpc server started on: %s", appConfig.Server.GRPC.GetFullAddress()))
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			return
		}
	}()

	logg.ServerLog("calendar is running...")

	<-ctx.Done()

	timeOutCtx, timeCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer timeCancel()

	if err := httpServer.Stop(timeOutCtx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}

	if err := grpcServer.Stop(timeOutCtx); err != nil {
		logg.Error("failed to stop grpc server: " + err.Error())
	}
}
