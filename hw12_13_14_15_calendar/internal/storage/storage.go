package storage

import (
	"time"
)

type IStorage interface {
	Add(item *Event) (string, error)
	Update(item *Event) error
	Delete(id string) error
	FindItem(id string) (*Event, error)
	ListRange(start, end *time.Time) ([]Event, error)
	GetEventRemind(now time.Time) ([]Event, error)
	DropOldEvents(year int) (int64, error)
}
