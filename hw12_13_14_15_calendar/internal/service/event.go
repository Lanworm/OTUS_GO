package service

import (
	"fmt"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/enum"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
)

type Event struct {
	logger  *logger.Logger
	storage storage.IStorage
}

func NewEventService(
	logger *logger.Logger,
	storage storage.IStorage,
) *Event {
	return &Event{
		logger:  logger,
		storage: storage,
	}
}

func (e *Event) CreateEvent(evt *storage.Event) (string, error) {
	res, err := e.storage.Add(evt)
	if err != nil {
		return "", fmt.Errorf("create event: %w", err)
	}

	return res, nil
}

func (e *Event) GetEvent(id string) (*storage.Event, error) {
	res, err := e.storage.FindItem(id)
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}

	return res, nil
}

func (e *Event) UpdateEvent(evt *storage.Event) error {
	err := e.storage.Update(evt)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	return nil
}

func (e *Event) DeleteEvent(id string) error {
	err := e.storage.Delete(id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}

	return nil
}

func (e *Event) ListEvent(r enum.RangeDuration) ([]storage.Event, error) {
	timeNow := time.Now().UTC()
	year, month, day := timeNow.Date()

	switch r {
	case enum.DAY:
		startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC)

		list, err := e.storage.ListRange(&startDate, &endDate)
		if err != nil {
			return nil, fmt.Errorf("load event per day list: %w", err)
		}

		return list, nil
	case enum.WEEK:
		startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, month, day+8, 0, 0, 0, 0, time.UTC)

		list, err := e.storage.ListRange(&startDate, &endDate)
		if err != nil {
			return nil, fmt.Errorf("load event per week list: %w", err)
		}

		return list, nil
	case enum.MONTH:
		startDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, month+1, day+1, 0, 0, 0, 0, time.UTC)

		list, err := e.storage.ListRange(&startDate, &endDate)
		if err != nil {
			return nil, fmt.Errorf("load event per month list: %w", err)
		}

		return list, nil
	}

	return nil, fmt.Errorf("unknown range duration: %s", r)
}
