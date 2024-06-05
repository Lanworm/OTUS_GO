package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage"
	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/storage/database"
)

type UserStorage struct {
	db     database.IDatabase
	logger *logger.Logger
}

func NewUserStorage(
	db database.IDatabase,
	logger *logger.Logger,
) *UserStorage {
	return &UserStorage{
		db:     db,
		logger: logger,
	}
}

func (s *UserStorage) Add(item storage.User) (string, error) {
	query := `insert into c_user(first_name, last_name) values($1, $2) returning id`

	var id string
	err := s.db.DB().QueryRow(
		context.Background(),
		query,
		item.FirstName,
		item.LastName,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("add user-storage exec query: %w", err)
	}

	return id, nil
}

func (s *UserStorage) Update(id string, item storage.User) error {
	query := `update c_user set first_name=$1, last_name=$2, updated_at=$3 where id=$5`

	res, err := s.db.DB().Query(
		context.Background(),
		query,
		item.FirstName,
		item.LastName,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("update user-storage exec query: %w", err)
	}

	err = res.Err()
	if err != nil {
		return fmt.Errorf("update user-storage query: %w", err)
	}

	return nil
}

func (s *UserStorage) Delete(id string) error {
	query := `delete from c_user where id=$1`

	res, err := s.db.DB().Query(
		context.Background(),
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete user-storage exec query: %w", err)
	}

	err = res.Err()
	if err != nil {
		return fmt.Errorf("delete user-storage query: %w", err)
	}

	return nil
}

func (s *UserStorage) FindItem(id string) (storage.User, error) {
	query := `select id, first_name, last_name, created_at, updated_at from c_user where id=$1`

	var usr storage.User
	err := s.db.DB().QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&usr.ID,
		&usr.FirstName,
		&usr.LastName,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)
	if err != nil {
		return storage.User{}, fmt.Errorf("find user: %w", err)
	}

	return usr, nil
}

func (s *UserStorage) List() ([]storage.User, error) {
	query := `select id, first_name, last_name, created_at, updated_at from c_user`

	res, err := s.db.DB().Query(
		context.Background(),
		query,
	)
	if err != nil {
		return []storage.User{}, fmt.Errorf("list user-storage exec query: %w", err)
	}

	users := make([]storage.User, 0, 50)
	for res.Next() {
		var usr storage.User
		err = res.Scan(
			&usr.ID,
			&usr.FirstName,
			&usr.LastName,
			&usr.CreatedAt,
			&usr.UpdatedAt,
		)
		if err != nil {
			return users, fmt.Errorf("list user-storage scan query: %w", err)
		}

		users = append(users, usr)
	}

	return users, nil
}
