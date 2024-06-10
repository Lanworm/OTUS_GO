package storage

import "time"

type Event struct {
	ID            string    `db:"id"`
	Title         string    `db:"title"`
	StartDatetime time.Time `db:"start_datetime"`
	EndDatetime   time.Time `db:"end_datetime"`
	Description   string    `db:"description"`
	UserID        string    `db:"user_id"`
	RemindBefore  int       `db:"remind_before"`
}
