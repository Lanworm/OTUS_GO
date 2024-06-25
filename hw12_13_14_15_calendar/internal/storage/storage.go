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
}
