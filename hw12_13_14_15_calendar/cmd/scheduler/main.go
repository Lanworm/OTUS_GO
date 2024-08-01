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
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service/rabbitmq"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/database"
	sqlstorage "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/pkg/shortcuts"
	"github.com/streadway/amqp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./build/local/scheduler-config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config2.NewSchedulerConfig(configFile)
	shortcuts.FatalIfErr(err)

	logg, err := logger.New(config.Logger.Level, os.Stdout)
	shortcuts.FatalIfErr(err)

	duration, err := time.ParseDuration(config.Interval)
	shortcuts.FatalIfErr(err)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	db := database.New(
		&config.Database,
		logg,
	)
	defer db.Close(ctx)

	err = db.Connect(ctx)
	shortcuts.FatalIfErr(err)

	eventStorage := sqlstorage.NewEventStorage(db, logg)
	evtService := service.NewEventService(logg, eventStorage)

	t := time.Now()
	dropCount, err := evtService.DropOldEvents(t.Year())
	if err != nil {
		logg.Error(err)
	}

	logg.Info(fmt.Sprintf("dropped old items coutnt: %d", dropCount))
	rc, err := amqp.Dial(config.Rabbit.Dsn)
	shortcuts.FatalIfErr(err)

	defer rc.Close()

	rcch, err := rc.Channel()
	shortcuts.FatalIfErr(err)
	defer rcch.Close()

	rabbit := rabbitmq.NewRabbit(rc, rcch, &config.Rabbit)
	err = rabbit.InitialQueue()
	shortcuts.FatalIfErr(err)

	ticker := time.NewTicker(duration)
	logg.Info(fmt.Sprintf("scheduler starterd with duration: %s", config.Interval))
forloop:
	for {
		select {
		case <-ticker.C:
			events, err := evtService.GetEventRemind(time.Now())
			if err != nil {
				logg.Error(fmt.Errorf("load event for remid: %w", err))
				continue
			}

			if len(events) > 0 {
				for _, event := range events {
					logg.Info(fmt.Sprintf("process event: %s", event.ID))
				}
				err := rabbit.PublishMessages(events)
				if err != nil {
					logg.Error(err)
					cancel()
					break forloop
				}
			} else {
				logg.Info("no messages for send")
			}
		case <-ctx.Done():
			logg.Info("context cancelled")
			break forloop
		}
	}
}
