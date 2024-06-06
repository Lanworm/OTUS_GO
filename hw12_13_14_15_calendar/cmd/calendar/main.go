package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/app"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/server/http"
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

	config, err := config.NewConfig(configFile)
	shortcuts.FatalIfErr(err)

	logg, err := logger.New(config.Logger.Level, os.Stdout)
	shortcuts.FatalIfErr(err)

	eventStorage, _ := storageconf.NewStorage(ctx, *config, logg)

	calendar := app.New(logg, eventStorage)

	server := internalhttp.NewServer(logg, calendar, config.Server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
