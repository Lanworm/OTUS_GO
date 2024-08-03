package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	config2 "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/grpc/grpchandler"
	internalhttp "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http/httphandler"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/database"
	memorystorage "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/pkg/shortcuts"
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

	config, err := config2.NewConfig(configFile)
	shortcuts.FatalIfErr(err)

	logg, err := logger.New(config.Logger.Level, os.Stdout)
	shortcuts.FatalIfErr(err)

	var eventStorage storage.IStorage
	if config.Storage.InDatabase() {
		db := database.New(
			&config.Database,
			logg,
		)
		err := db.Connect(ctx)
		shortcuts.FatalIfErr(err)

		eventStorage = sqlstorage.NewEventStorage(db, logg)
	} else {
		eventStorage = memorystorage.New()
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	httpServer := internalhttp.NewHTTPServer(logg, config.Server.HTTP)
	evtService := service.NewEventService(logg, eventStorage)

	handlerHTTP := httphandler.NewHandler(logg, evtService)
	httpServer.RegisterRoutes(handlerHTTP)
	go func() {
		logg.ServerLog(fmt.Sprintf("http server started on: %s", config.Server.HTTP.GetFullAddress()))
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
		config.Server.GRPC,
	)
	go func() {
		logg.ServerLog(fmt.Sprintf("grpc server started on: %s", config.Server.GRPC.GetFullAddress()))
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
