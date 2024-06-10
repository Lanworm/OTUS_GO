package memorystorage

import (
	"sync"

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

func (s *Storage) Add(item storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()

	s.data[id] = item

	return id, nil
}

func (s *Storage) Update(id string, item storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = item

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

func (s *Storage) FindItem(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var item storage.Event

	item, ok := s.data[id]
	if ok {
		return item, storage.ErrEventNotFound
	}

	return item, nil
}

func (s *Storage) List() ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0, len(s.data))

	for _, item := range s.data {
		result = append(result, item)
	}

	return result, nil
}
