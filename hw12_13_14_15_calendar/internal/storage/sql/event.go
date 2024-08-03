package sqlstorage

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func (s *EventStorage) Add(item *storage.Event) (string, error) {
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

func (s *EventStorage) Update(item *storage.Event) error {
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
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("update event-storage exec query: %w", err)
	}

	defer res.Close()

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

	defer res.Close()

	err = res.Err()
	if err != nil {
		return fmt.Errorf("delete event-storage query: %w", err)
	}

	return nil
}

func (s *EventStorage) FindItem(id string) (*storage.Event, error) {
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
		return nil, fmt.Errorf("not found event: %w", err)
	}

	return &event, nil
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

	defer res.Close()

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

func (s *EventStorage) ListRange(start, end *time.Time) ([]storage.Event, error) {
	query := strings.Builder{}

	query.WriteString(`select id, title, start_datetime, end_datetime, description, user_id, remind_before from c_event`)

	type whereStruct struct {
		t     time.Time
		field string
		expr  string
	}

	where := make([]whereStruct, 0, 2)

	if start != nil {
		where = append(where, whereStruct{t: *start, expr: ">=", field: "start_datetime"})
	}

	if end != nil {
		where = append(where, whereStruct{t: *end, expr: "<", field: "start_datetime"})
	}

	if len(where) > 0 {
		query.WriteString(" where ")
	}

	args := make([]any, 0, 2)
	for i, w := range where {
		query.WriteString(fmt.Sprintf(" %s %s $%d", w.field, w.expr, i+1))
		if len(where)-1 > i {
			query.WriteString(" AND ")
		}
		args = append(args, w.t.Format(time.RFC3339))
	}

	res, err := s.db.DB().Query(
		context.Background(),
		query.String(),
		args...,
	)
	if err != nil {
		return []storage.Event{}, fmt.Errorf("list event-storage exec query: %w", err)
	}

	defer res.Close()

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

func (s *EventStorage) GetEventRemind(now time.Time) ([]storage.Event, error) {
	query := strings.Builder{}

	query.WriteString(`select id, title, start_datetime, end_datetime, description, user_id, remind_before from c_event where start_datetime - $1 >= make_interval(0, 0, 0, 0, 0, 0, remind_before) and start_datetime - $2 < make_interval(0, 0, 0, 0, 0, 0, remind_before)`) //nolint:lll

	args := make([]any, 0, 2)
	args = append(args, now.Format(time.RFC3339))
	args = append(args, now.Add(time.Minute).Format(time.RFC3339))

	res, err := s.db.DB().Query(
		context.Background(),
		query.String(),
		args...,
	)
	if err != nil {
		return []storage.Event{}, fmt.Errorf("remind event-storage exec query: %w", err)
	}

	defer res.Close()

	events := make([]storage.Event, 0, 10)
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
			return events, fmt.Errorf("remind event-storage scan query: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *EventStorage) DropOldEvents(year int) (dropCount int64, err error) {
	query := `delete from c_event where $1 > date_part('year', start_datetime)`

	res, err := s.db.DB().Query(
		context.Background(),
		query,
		year,
	)

	defer func() {
		res.Close()

		if err == nil {
			dropCount = res.CommandTag().RowsAffected()
		}
	}()

	if err != nil {
		return 0, fmt.Errorf("drop old events event-storage scan query: %w", err)
	}

	return
}
