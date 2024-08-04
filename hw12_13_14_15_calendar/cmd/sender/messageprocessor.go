package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/streadway/amqp"
)

type MessageProcessor struct {
	Logger         *logger.Logger
	MessageChannel <-chan amqp.Delivery
	SendStorage    *sqlstorage.SendStorage
}

func NewMessageProcessor(
	logger *logger.Logger,
	messageChannel <-chan amqp.Delivery,
	sendStorage *sqlstorage.SendStorage,
) *MessageProcessor {
	return &MessageProcessor{
		Logger:         logger,
		MessageChannel: messageChannel,
		SendStorage:    sendStorage,
	}
}

func (m *MessageProcessor) ProcessMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			m.Logger.Info("context cancelled")
			return
		case d := <-m.MessageChannel:
			m.Logger.Info(fmt.Sprintf("received a message: %s", d.Body))

			var evt storage.Event
			err := json.Unmarshal(d.Body, &evt)
			if err != nil {
				m.Logger.Info(fmt.Sprintf("error decoding JSON: %s", err))
				continue
			}

			err = m.SendStorage.Add(fmt.Sprintf(
				"Notification by event #%s, %s: %s",
				evt.ID,
				evt.Title,
				evt.Description,
			))
			if err != nil {
				m.Logger.Error(fmt.Errorf("fail send notification: %w", err))
			} else {
				m.Logger.Info("notification success sent")
			}

			if err := d.Ack(false); err != nil {
				m.Logger.Error(fmt.Errorf("error acknowledging message: %w", err))
			} else {
				m.Logger.Info("acknowledged message")
			}
		}
	}
}
