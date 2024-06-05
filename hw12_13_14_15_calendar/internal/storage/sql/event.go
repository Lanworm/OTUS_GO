package sqlstorage

import (
	"context"
	"fmt"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/database"
)

type EventStorage struct {
	db     database.IDatabase
	logger *logger.Logger
}

func NewEventStorage(
	db database.IDatabase,
	logger *logger.Logger,
) *EventStorage {
	return &EventStorage{
		db:     db,
		logger: logger,
	}
}

func (s *EventStorage) Add(item storage.Event) (string, error) {
	//nolint:lll
	query := `insert into c_event(title, start_datetime, end_datetime, description, user_id, remind_before) values($1, $2, $3, $4, $5, $6) returning id`

	var id string
	err := s.db.DB().QueryRow(
		context.Background(),
		query,
		item.Title,
		item.StartDatetime,
		item.EndDatetime,
		item.Description,
		item.UserID,
		item.RemindBefore,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("add event-storage exec query: %w", err)
	}

	return id, nil
}

func (s *EventStorage) Update(id string, item storage.Event) error {
	//nolint:lll
	query := `update c_event set title=$1, start_datetime=$2, end_datetime=$3, description=$4, remind_before=$5 where id=$6`

	res, err := s.db.DB().Query(
		context.Background(),
		query,
		item.Title,
		item.StartDatetime,
		item.EndDatetime,
		item.Description,
		item.RemindBefore,
		id,
	)
	if err != nil {
		return fmt.Errorf("update event-storage exec query: %w", err)
	}

	err = res.Err()
	if err != nil {
		return fmt.Errorf("update event-storage query: %w", err)
	}

	return nil
}

func (s *EventStorage) Delete(id string) error {
	query := `delete from c_event where id=$1`

	res, err := s.db.DB().Query(
		context.Background(),
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete event-storage exec query: %w", err)
	}

	err = res.Err()
	if err != nil {
		return fmt.Errorf("delete event-storage query: %w", err)
	}

	return nil
}

func (s *EventStorage) FindItem(id string) (storage.Event, error) {
	query := `select id, title, start_datetime, end_datetime, description, user_id, remind_before from c_event where id=$1`

	var event storage.Event
	err := s.db.DB().QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&event.ID,
		&event.Title,
		&event.StartDatetime,
		&event.EndDatetime,
		&event.Description,
		&event.UserID,
		&event.RemindBefore,
	)
	if err != nil {
		return storage.Event{}, fmt.Errorf("find event: %w", err)
	}

	return event, nil
}

func (s *EventStorage) List() ([]storage.Event, error) {
	query := `select id, title, start_datetime, end_datetime, description, user_id, remind_before from c_event`

	res, err := s.db.DB().Query(
		context.Background(),
		query,
	)
	if err != nil {
		return []storage.Event{}, fmt.Errorf("list event-storage exec query: %w", err)
	}

	events := make([]storage.Event, 0, 50)
	for res.Next() {
		var event storage.Event
		err = res.Scan(
			&event.ID,
			&event.Title,
			&event.StartDatetime,
			&event.EndDatetime,
			&event.Description,
			&event.UserID,
			&event.RemindBefore,
		)
		if err != nil {
			return events, fmt.Errorf("list event-storage scan query: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}
