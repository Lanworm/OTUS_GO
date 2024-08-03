package memorystorage

import (
	"sync"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:   sync.RWMutex{},
		data: make(map[string]storage.Event, 100),
	}
}

func (s *Storage) Add(item *storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()

	item.StartDatetime = item.StartDatetime.UTC()
	item.EndDatetime = item.EndDatetime.UTC()

	item.ID = id
	s.data[id] = *item

	return id, nil
}

func (s *Storage) Update(item *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	innerItem, ok := s.data[item.ID]
	if !ok {
		return storage.ErrEventNotFound
	}

	innerItem.Title = item.Title
	innerItem.Description = item.Description
	innerItem.StartDatetime = item.StartDatetime
	innerItem.EndDatetime = item.EndDatetime
	innerItem.RemindBefore = item.RemindBefore

	s.data[item.ID] = innerItem

	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[id]
	if !ok {
		return storage.ErrEventNotFound
	}

	delete(s.data, id)

	return nil
}

func (s *Storage) FindItem(id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.data[id]
	if !ok {
		return nil, storage.ErrEventNotFound
	}

	return &item, nil
}

func (s *Storage) ListRange(start, end *time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0, len(s.data))

	for _, item := range s.data {
		if (item.StartDatetime.Equal(*start) || item.StartDatetime.After(*start)) &&
			(item.StartDatetime.Equal(*end) || item.StartDatetime.Before(*end)) {
			result = append(result, item)
		}
	}

	return result, nil
}

func (s *Storage) GetEventRemind(now time.Time) ([]storage.Event, error) {
	start := now
	end := now.Add(time.Minute)

	return s.ListRange(&start, &end)
}

func (s *Storage) DropOldEvents(year int) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dropCount := int64(0)

	for k, item := range s.data {
		if year > item.StartDatetime.Year() {
			delete(s.data, k)
			dropCount++
		}
	}

	return dropCount, nil
}
