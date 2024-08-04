package storageconf

import (
	"context"
	"errors"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/config"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/database"
	memorystorage "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/pkg/shortcuts"
)

var ErrUnknownDatabaseType = errors.New("unexpected database type")

func NewStorage(ctx context.Context, config config.Config, logg *logger.Logger) (storage.IStorage, error) {
	switch config.Storage.InDatabase() {
	case true:
		db := database.New(&config.Database, logg)
		err := db.Connect(ctx)
		shortcuts.FatalIfErr(err)

		defer func() {
			err := db.Connect(ctx)
			if err != nil {
				logg.Error(err)
			}
		}()
		return sqlstorage.NewEventStorage(db, logg), err
	case false:
		return memorystorage.New(), nil
	}
	return nil, ErrUnknownDatabaseType
}
