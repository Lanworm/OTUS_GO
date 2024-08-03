package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config2 "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service/rabbitmq"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
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

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	rc, err := amqp.Dial(config.Rabbit.Dsn)
	shortcuts.FatalIfErr(err)
	defer rc.Close()

	rcch, err := rc.Channel()
	shortcuts.FatalIfErr(err)
	defer rcch.Close()

	rabbit := rabbitmq.NewRabbit(rc, rcch, &config.Rabbit)
	err = rabbit.InitialQueue()
	shortcuts.FatalIfErr(err)

	messageChannel, err := rabbit.Consume()
	shortcuts.FatalIfErr(err)

	db := database.New(
		&config.Database,
		logg,
	)
	defer db.Close(ctx)

	err = db.Connect(ctx)
	shortcuts.FatalIfErr(err)

	sendStorage := sqlstorage.NewSendStorage(db, logg)

foorloop:
	for {
		select {
		case <-ctx.Done():
			logg.Info("context cancelled")
			break foorloop
		case d := <-messageChannel:
			logg.Info(fmt.Sprintf("received a message: %s", d.Body))

			var evt storage.Event
			err := json.Unmarshal(d.Body, &evt)
			if err != nil {
				logg.Info(fmt.Sprintf("error decoding JSON: %s", err))
			}

			err = sendStorage.Add(fmt.Sprintf(
				"Notification by event #%s, %s: %s",
				evt.ID,
				evt.Title,
				evt.Description,
			))
			if err != nil {
				logg.Error(fmt.Errorf("fail send notification: %w", err))
			} else {
				logg.Info("notification success sent")
			}

			if err := d.Ack(false); err != nil {
				logg.Error(fmt.Errorf("error acknowledging message: %w", err))
			} else {
				logg.Info("acknowledged message")
			}
		}
	}
}
