package sqlstorage

import (
	"context"
	"fmt"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/database"
)

type SendStorage struct {
	db     database.IDatabase
	logger *logger.Logger
}

func NewSendStorage(
	db database.IDatabase,
	logger *logger.Logger,
) *SendStorage {
	return &SendStorage{
		db:     db,
		logger: logger,
	}
}

func (s *SendStorage) Add(message string) error {
	query := `insert into c_send(message) values($1)`

	cmd, err := s.db.DB().Exec(
		context.Background(),
		query,
		message,
	)
	if err != nil {
		return fmt.Errorf("add send-storage exec query: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("add send-storage data has not inserting")
	}

	return nil
}
