package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/service/rabbitmq"
)

// ProcessEvents обрабатывает события.
func ProcessEvents(
	ctx context.Context,
	ticker *time.Ticker,
	evtService *service.Event,
	logg *logger.Logger,
	rabbit *rabbitmq.Rabbit,
	cancel context.CancelFunc,
) {
	for {
		select {
		case <-ticker.C:
			// Получаем события для напоминания
			events, err := evtService.GetEventRemind(time.Now())
			if err != nil {
				logg.Error(fmt.Errorf("load event for remid: %w", err))
				continue
			}

			// Проверяем, есть ли события для отправки
			if len(events) == 0 {
				logg.Info("no messages for send")
				continue
			}

			// Логируем каждое событие
			for _, event := range events {
				logg.Info(fmt.Sprintf("process event: %s", event.ID))
			}

			// Публикуем сообщения
			err = rabbit.PublishMessages(events)
			if err != nil {
				logg.Error(err)
				cancel()
				return
			}

		case <-ctx.Done():
			// Логируем, что контекст отменен
			logg.Info("context cancelled")
			return
		}
	}
}
