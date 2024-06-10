package app

import (
	"context"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	logger  Logger
	storage storage.IStorage
}

type Logger interface { // TODO
}

func New(logger Logger, storage storage.IStorage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
