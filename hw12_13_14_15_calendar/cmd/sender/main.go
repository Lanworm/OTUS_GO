package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	config "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
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

	config, err := config.NewSchedulerConfig(configFile)
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
	// Создание обработчика
	h := NewMessageProcessor(logg, messageChannel, sendStorage)

	// Запуск обработки сообщений
	h.ProcessMessages(ctx)
}
